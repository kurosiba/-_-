package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	fileName := "logfile.txt" //開きたいファイルの名前
	N := 2                    //故障とみなすためのタイムアウトの回数

	var logData [100][3]string //タイムアウトしているサーバーのタイムスタンプとアドレス,連続でタイムアウトした回数を格納する配列
	serverCount := 0           //タイムアウトしているサーバーの数
	serverChecker := false     //同じアドレスから連続でタイムアウトが出ていないかの確認するための変数

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
			serverChecker = false

			if arry1[1] == logData[i][1] {
				serverChecker = true //logData[][]に初期値を入れるかどうかの判定

				if arry1[2] == "-" { //タイムアウトした回数の計測と保存
					timeOutCount, _ := strconv.Atoi(logData[i][2])
					timeOutCount++
					logData[i][2] = strconv.Itoa(timeOutCount)
				}

				break
			}
		}

		if arry1[2] == "-" && serverChecker == false { //タイムアウトが出たサーバーのタイムスタンプとアドレス、タイムアウトの回数をlogData[]に代入している(初期化)
			serverCount++
			timeOutCount := 1
			str := strconv.Itoa(timeOutCount)
			logData[serverCount-1][0] = arry1[0]
			logData[serverCount-1][1] = arry1[1]
			logData[serverCount-1][2] = str
		}

		for i := 0; i < serverCount; i++ {
			timeOutCount, _ := strconv.Atoi(logData[i][2])

			if arry1[1] == logData[i][1] && timeOutCount >= N { //サーバーが復帰したタイミングでアドレスと故障時間をコンソールに出力する

				if arry1[2] != "-" {
					time1, _ := strconv.Atoi(arry1[0])
					time2, _ := strconv.Atoi(logData[i][0])
					str := fmt.Sprintf("%014d", time1-time2)

					fmt.Println(logData[i][1] + " Server failure")
					fmt.Println(logData[i][1] + " Server failure period is " + str)

					logData[i][0] = "" //タイムアウトしていたサーバーのタイムスタンプの初期化
					logData[i][1] = "" //タイムアウトしていたサーバーのアドレスの初期化
					logData[i][2] = "0"
					serverChecker = false
				} else {
					fmt.Println(logData[i][1] + " Server failure")
				}
			}

		}
	}

	if err = scanner.Err(); err != nil { //読み取りでエラーが出たときのエラー処理
		fmt.Println("scan error")
	}

}
