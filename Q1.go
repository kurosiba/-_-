package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	FileName := "logfile.txt" //開きたいファイルの名前

	f, err := os.Open(FileName)
	if err != nil { //ファイルをOpenしたときのエラー処理
		fmt.Println("open error")
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	/*一行ずつ読み取り、カンマごとに分割し、-が出るとその一つ前の要素(アドレス)をコンソールに出力する。*/
	for scanner.Scan() {
		arry1 := strings.Split(scanner.Text(), ",")
		if arry1[2] == "-" {
			fmt.Println(arry1[1])
		}
	}

	if err = scanner.Err(); err != nil { //読み取りでエラーが出たときのエラー処理
		fmt.Println("scan error")
	}

}
