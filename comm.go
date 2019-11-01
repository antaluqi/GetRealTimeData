package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// 整形最小值（两个数）
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// 整形最大值（两个数）
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// 打印错误
func check(err error) {

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

//字符串编码Gbk -> Utf8     go get golang.org/x/text
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

//字符串编码Utf8 -> Gbk     go get golang.org/x/text
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func strArrUnitCombin(strArr []string, cStr string, pos int) []string {

	if len(strArr) == 0 {
		return strArr
	}
	var strResult []string
	for _, str := range strArr {
		if pos <= 0 {
			strResult = append(strResult, cStr+str)
		} else {
			strResult = append(strResult, str+cStr)
		}
	}
	return strResult
}
