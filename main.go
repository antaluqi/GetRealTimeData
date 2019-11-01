package main

import (
	"bytes"
	"database/sql"
	"fmt"
	//"io"
	"math"
	"net/http"
	//"os"
	"strings"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

//下载股票代码列表
func stocklist() []string {
	t1 := time.Now()

	URL := "http://file.tushare.org/tsdata/h/hq.csv" // 代码下载网址
	var codeList []string
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer resp.Body.Close()
	buf := bytes.NewBuffer(make([]byte, 0, 512))
	buf.ReadFrom(resp.Body)
	t := string(buf.Bytes())
	Arr := strings.Split(t, ",,")
	for i, s := range Arr {
		arr := strings.Split(s, ",")
		if i == 0 || len(arr) < 10 {
			continue
		}
		//fmt.Println(string(arr[0][2]))
		if arr[0][2] == '6' {
			codeList = append(codeList, "sh"+arr[0][2:])
		} else if arr[0][2] == '0' || arr[0][2] == '3' {
			codeList = append(codeList, "sz"+arr[0][2:])
		} else {
			codeList = append(codeList, arr[0][2:])
		}

	}
	// 加入预设的指数代码列表
	codeList = append(codeList, index_codelist...)
	fmt.Println("获取股票列表计时: ", time.Since(t1))
	return codeList
}

func get_realtime_data(codelist []string) string {
	//已下日线复权数据怎么办？？
	// http://web.ifzq.gtimg.cn/appstock/app/fqkline/get?_var=kline_dayqfq&param=sz300755,day,2019-10-01,2019-10-12,20,qfq&r=0.123458

	st := 0
	ed := 60
	max_ed := len(codelist)
	url_head := "http://qt.gtimg.cn/q="
	str := ""
	t1 := time.Now()
	for {
		if st >= max_ed {
			break
		}
		url := url_head + strings.Join(codelist[st:min(ed, max_ed)], ",")
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		defer resp.Body.Close()
		buf := bytes.NewBuffer(make([]byte, 0, 512))
		buf.ReadFrom(resp.Body)
		utf8, _ := GbkToUtf8(buf.Bytes())
		str = str + string(utf8)
		st = st + 60
		ed = ed + 60
	}
	//计时显示
	fmt.Println("下载计时: ", time.Since(t1))
	return str
}

func get_realtime_data2(codelist []string, dtype string) []string {
	var iNo int
	var url_head string
	if dtype == "sl" {
		iNo = iNo_sl
		pCount = pCount_sl //下载线程数，需要放入conf.ini
		url_head = "http://qt.gtimg.cn/q="
	} else if dtype == "ss" {
		iNo = iNo_ss
		pCount = pCount_ss //下载线程数，需要放入conf.ini
		url_head = "http://qt.gtimg.cn/q="
		codelist = strArrUnitCombin(codelist, "s_", 0)
	} else if dtype == "sj" {
		iNo = iNo_sj
		pCount = pCount_sj //下载线程数，需要放入conf.ini
		url_head = "http://hq.sinajs.cn/rn=xppzh&list="
	} else {

	}

	pFinish = make(chan int, pCount) // 下载阻塞信道，缓冲大小为下载线程个数(全局)

	c := make(chan string, int(math.Ceil(float64(len(codelist))/float64(iNo)))) // 代码信道
	v := make(chan string, int(math.Ceil(float64(len(codelist))/float64(iNo)))) // 值信道
	st := 0
	ed := iNo
	max_ed := len(codelist)

	//t1 := time.Now()
	for {
		if st >= max_ed {
			break
		}
		url := url_head + strings.Join(codelist[st:min(ed, max_ed)], ",")
		c <- url
		//fmt.Printf("代码%d-%d放入信道\n", st, ed)
		st = st + iNo
		ed = ed + iNo
	}
	fmt.Println("代码存放完毕=======================================")
	// 下载线程开始
	for i := 0; i < pCount; i++ {
		go producer(c, v, i)
	}
	// ---------------------下载阻塞
	for i := 0; i < pCount; i++ {
		<-pFinish
	}

	//从值信道中取出数值合并
	str := ""
	for {
		if len(v) == 0 {

			break
		}
		str = str + <-v
	}
	// 转换成数组
	var arr []string
	if dtype == "sl" || dtype == "ss" {
		nrep := strings.NewReplacer("v_", "'", "s_", "", `="`, "','", "~", "','", `";`, "';")
		strR := strings.Replace(strings.Replace(nrep.Replace(str), "''", "0", -1), "'", "", -1)
		arr = strings.Split(strR, ";")
	} else if dtype == "sj" {
		nrep := strings.NewReplacer("var hq_str_", "", `="`, ",", `";`, ";")
		strR := nrep.Replace(str)
		arr = strings.Split(strR, ";")
	} else {

	}

	//计时显示
	//fmt.Println("下载计时: ", time.Since(t1))

	return arr
}

func to_psql(arr []string, tx *sql.Tx) {
	if len(arr) <= 1 {
		fmt.Println("获取数据为空")
	}
	arr = arr[0 : len(arr)-1]
	stmt, _ := tx.Prepare(pq.CopyIn(realtime_tablename, columeName[realtime_tablename]...))
	for _, istr := range arr {
		iarr := strings.Split(strings.Replace(istr, "\n", "", -1), ",")[0:len(columeName[realtime_tablename])]
		inarr := make([]interface{}, len(iarr))
		for i, v := range iarr {
			inarr[i] = v
		}
		_, err := stmt.Exec(inarr...)
		check(err)

	}
	_, err := stmt.Exec()
	check(err)
	stmt.Close()
}

func download() {
	//初始化 要放在外面
	initialize()
	//下载股票列表
	cl := stocklist()

	if dtype == "sl" {
		realtime_tablename = sina_Lrealtime_tableName
	} else if dtype == "ss" {
		realtime_tablename = sina_Srealtime_tableName
	} else if dtype == "sj" {
		realtime_tablename = sina_Jrealtime_tableName
	} else {

	}

	//连接数据库
	db := getDB("postgresql") //连接数据库
	defer db.Close()

	//检测是否存在表，不存在则建立
	if !checkTable(db, realtime_tablename) {
		createTable(db, realtime_tablename)
		fmt.Println(fmt.Sprintf("重建%s表\n", realtime_tablename))
	} else {
		fmt.Println(fmt.Sprintf("已经存在%s表\n", realtime_tablename))
	}

	for i := 0; i < loopNo; i++ {
		//计时开始
		t1 := time.Now()
		//下载数据
		arr := get_realtime_data2(cl, dtype)
		//arr := get_realtime_data([]string{"sh600118", "sh000001", "sh000002"})
		fmt.Println("下载计时: ", time.Since(t1))
		// 存入数据库
		tx, err := db.Begin()
		check(err)
		_, err = tx.Exec(fmt.Sprintf("truncate table %s", realtime_tablename))
		check(err)
		to_psql(arr, tx)
		tx.Commit()
		//计时显示
		fmt.Println("下载+处理计时: ", time.Since(t1))
		//是否延时判断
		if loopTM > 0 {
			tf := time.Duration(loopTM)*time.Millisecond - time.Since(t1)
			if tf > 0 {
				fmt.Println("延时: ", tf)
				time.Sleep(tf)
			}
		}

	}

}

// 生产者函数（下载）
func producer(c chan string, v chan string, pname int) {
	var str string
	for {
		str = ""
		if len(c) == 0 {
			break
		}
		url := <-c
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		defer resp.Body.Close()
		buf := bytes.NewBuffer(make([]byte, 0, 512))
		buf.ReadFrom(resp.Body)
		utf8, _ := GbkToUtf8(buf.Bytes())
		str = str + string(utf8)
		v <- str

	}
	pFinish <- 0
	fmt.Printf(">>下载线程%d结束<< \n", pname)
}

// main=============================================================================
func main() {

	download()

}
