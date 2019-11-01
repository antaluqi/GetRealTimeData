package main

import (
	"database/sql"
	"fmt"

	"github.com/Unknwon/goconfig"
	_ "github.com/lib/pq"
)

//
//	"time"

//	"github.com/Unknwon/goconfig"
// 读取.ini 文件
func readIni() {
	//读取ini文件
	cfg, err := goconfig.LoadConfigFile("conf.ini")
	if err != nil {
		check(err)

	}
	//-------------------------------[Postgresql]

	host_psql, err = cfg.GetValue("Postgresql", "host_psql")
	if err != nil {
		check(err)
	}
	port_psql, err = cfg.Int("Postgresql", "port_psql")
	if err != nil {
		check(err)
	}
	user_psql, err = cfg.GetValue("Postgresql", "user_psql")
	if err != nil {
		check(err)
	}
	password_psql, err = cfg.GetValue("Postgresql", "password_psql")
	if err != nil {
		check(err)
	}
	dbname_psql, err = cfg.GetValue("Postgresql", "dbname_psql")
	if err != nil {
		check(err)
	}
	//-------------------------------[Download_DateRange]

	dtype, err = cfg.GetValue("Download_DateRange", "dtype")
	if err != nil {
		check(err)
	}
	iNo_sl, err = cfg.Int("Download_DateRange", "iNo_sl")
	if err != nil {
		check(err)
	}
	iNo_ss, err = cfg.Int("Download_DateRange", "iNo_ss")
	if err != nil {
		check(err)
	}
	iNo_sj, err = cfg.Int("Download_DateRange", "iNo_sj")
	if err != nil {
		check(err)
	}

	//-------------------------------[Thread]

	pCount_sl, err = cfg.Int("Thread", "pCount_sl")
	if err != nil {
		fmt.Println(err)
	}
	pCount_ss, err = cfg.Int("Thread", "pCount_ss")
	if err != nil {
		fmt.Println(err)
	}
	pCount_sj, err = cfg.Int("Thread", "pCount_sj")
	if err != nil {
		fmt.Println(err)
	}
	//-------------------------------[Custom]

	loopNo, err = cfg.Int("Custom", "loopNo")
	if err != nil {
		fmt.Println(err)

	}
	loopTM, err = cfg.Float64("Custom", "loopTM")
	if err != nil {
		fmt.Println(err)

	}
}

// 初始化函数
func initialize() {
	readIni()
	// ------------------------表字段信息

	columeName[sina_Lrealtime_tableName] = []string{"code",
		"market", "name", "scode", "price", "preprice", "open", "volume", "buying", "selling", "b1p",
		"b1v", "b2p", "b2v", "b3p", "b3v", "b4p", "b4v", "b5p", "b5v", "s1p",
		"s1v", "s2p", "s2v", "s3p", "s3v", "s4p", "s4v", "s5p", "s5v", "rec_tran_num",
		"time", "change", "p_change", "high", "low", "p_v_a", "volume2", "amount", "turnover", "pe_ttm",
		"status", "high2", "low2", "range", "fvalues", "abvalues", "pb", "limit_up_price", "limit_down_price", "volratio",
		"deviation", "avgprice", "pe_d", "pe_s", "unknow2", "unknow3", "unknow4", "amount2", "unknow5", "unknow6",
		"unknow7", "unknow8", "unknow9", "unknow10", "unknow11"} //日线表 列名
	columeType[sina_Lrealtime_tableName] = []string{"text",
		"text", "text", "text", "real", "real", "real", "real", "real", "real", "real",
		"real", "real", "real", "real", "real", "real", "real", "real", "real", "real",
		"real", "real", "real", "real", "real", "real", "real", "real", "real", "text",
		"text", "real", "real", "real", "real", "text", "real", "real", "real", "real",
		"text", "real", "real", "real", "real", "real", "real", "real", "real", "real",
		"real", "real", "real", "real", "text", "text", "text", "text", "text", "text",
		"text", "text", "text", "text", "text"} //日线表 列属性
	columeName[sina_Srealtime_tableName] = []string{"code", "market", "name", "scode", "price", "change", "p_change", "volume", "amount", "unknow1", "abvalues", "unknow2"}
	columeType[sina_Srealtime_tableName] = []string{"text", "text", "text", "text", "real", "real", "real", "real", "real", "text", "real", "text"}
	columeName[sina_Jrealtime_tableName] = []string{"code", "name", "open", "pre_close", "price", "high", "low", "bid", "ask", "volumn",
		"amount", "b1v", "b1p", "b2v", "b2p", "b3v", "b3p", "b4v", "b4p", "b5v",
		"b5p", "s1v", "s1p", "s2v", "s2p", "s3v", "s3p", "s4v", "s4p", "s5v",
		"s5p", "date", "time"}
	columeType[sina_Jrealtime_tableName] = []string{"text", "text", "real", "real", "real", "real", "real", "real", "real", "real",
		"real", "real", "real", "real", "real", "real", "real", "real", "real", "real",
		"real", "real", "real", "real", "real", "real", "real", "real", "real", "real",
		"real", "date", "time without time zone"}
}

// 获取数据库连接
func getDB(DBname string) *sql.DB {
	var db *sql.DB
	switch DBname {
	case "postgresql":
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host_psql, port_psql, user_psql, password_psql, dbname_psql)
		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}

		err = db.Ping()
		if err != nil {
			panic(err)
		}
		fmt.Println("%s Successfully connected!", DBname)

		return db
	default:
		fmt.Println("没有这个数据库模块")
		return db

	}
}

// 检测有没有与预设相同的表
func checkTable(db *sql.DB, tableName string) bool {
	in_tablelist := false
	for _, tb := range table_name_list {
		if tableName == tb {
			in_tablelist = true
			break
		}
	}
	if !in_tablelist {
		fmt.Printf("表格 %s 的预设不存在！\n", tableName)
		return in_tablelist
	}

	rows, err := db.Query("select tablename from pg_tables where schemaname='public'")
	check(err)
	result := false
	var row string
	for rows.Next() {
		err := rows.Scan(&row)
		check(err)
		if row == tableName {
			result = true
			break
		}
	}
	if !result {
		fmt.Printf("%s 中没有表格 %s 存在！\n", dbname_psql, tableName)
		return result
	}

	rows, err = db.Query(fmt.Sprintf(`SELECT
											column_name
											,data_type
									  FROM information_schema.columns
									  WHERE table_schema = 'public' and table_name='%s'`, tableName))

	var colume_name, data_type string
	i := -1
	for rows.Next() {
		i++
		if i+1 > len(columeName[tableName]) {
			fmt.Printf("表格%s列数与预设不符,预设为%d列，实际已经超过\n", tableName, len(columeName[tableName]))
			return false
		}
		err := rows.Scan(&colume_name, &data_type)
		check(err)

		if columeName[tableName][i] != colume_name || columeType[tableName][i] != data_type {
			fmt.Printf("表格%s中第%d列与预设不符合,预设为%s(%s),实际为%s(%s)，\n", tableName, i+1, columeName[tableName][i], columeType[tableName][i], colume_name, data_type)
			return false
		}

	}
	if i+1 != len(columeName[tableName]) {
		fmt.Printf("表格%s列数与预设不符,预设为%d列,实际为%d列\n", tableName, len(columeName[tableName]), i+1)
		return false
	}
	return true
}

// 创建或替换表
func createTable(db *sql.DB, tableName string) bool {
	in_tablelist := false
	for _, tb := range table_name_list {
		if tableName == tb {
			in_tablelist = true
			break
		}
	}
	if !in_tablelist {
		fmt.Printf("表格 %s 的预设不存在！\n", tableName)
		return in_tablelist
	}
	strBody := ""
	strHeader := fmt.Sprintf("create table %s(\n", tableName)
	strTail := ")"
	for i, _ := range columeName[tableName] {
		strBody = strBody + columeName[tableName][i] + " " + columeType[tableName][i] + ",\n"
	}
	sqlStr := strHeader + strBody[0:len(strBody)-2] + strTail

	_, err := db.Exec(fmt.Sprintf("drop table if exists %s CASCADE", tableName))
	check(err)
	_, err = db.Exec(sqlStr)
	check(err)
	return true
}
