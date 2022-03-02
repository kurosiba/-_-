package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	var logData [100][2]string //タイムアウトしているサーバーのタイムスタンプとアドレスを格納する配列
	serverCount := 0           //タイムアウトしているサーバーの数
	serverChecker := false     //同じアドレスから連続でタイムアウトが出ていないかの確認するための変数
	fileName := "logfile.txt"  //開きたいファイルの名前

	f, err := os.Open(fileName)
	if err != nil { //ファイルをOpenしたときのエラー処理
		fmt.Println("open error")
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	/*
	*一行ずつログファイルを読み取り、カンマごとに分割し、アドレスと復帰時間をコンソールに出力する。
	 */
	for scanner.Scan() { //ログファイルを一行ずつ読む
		arry1 := strings.Split(scanner.Text(), ",") //カンマ毎に分割して、{タイムスタンプ,アドレス,応答時間}として格納している

		for i := 0; i < serverCount; i++ { //同じアドレスから連続でタイムアウトが出ていないかの確認
			if arry1[1] == logData[i][1] {
				serverChecker = true
			}
		}

		if arry1[2] == "-" && serverChecker == false { //タイムアウトが出たサーバーのタイムスタンプとアドレスをlogData[]に代入している
			serverCount++
			logData[serverCount-1][0] = arry1[0]
			logData[serverCount-1][1] = arry1[1]
		}

		for i := 0; i < serverCount; i++ {
			if arry1[1] == logData[i][1] && arry1[2] != "-" { //サーバーが復帰したタイミングでアドレスと故障時間をコンソールに出力する
				time1, _ := strconv.Atoi(arry1[0])
				time2, _ := strconv.Atoi(logData[i][0])
				fmt.Println(logData[i][1])
				str := fmt.Sprintf("%014d", time1-time2)
				fmt.Println(str)

				logData[i][0] = "" //タイムアウトしていたサーバーのタイムスタンプの初期化
				logData[i][1] = "" //タイムアウトしていたサーバーのアドレスの初期化
			}
		}
	}

	if err = scanner.Err(); err != nil { //読み取りでエラーが出たときのエラー処理
		fmt.Println("scan error")
	}

}
