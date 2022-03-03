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
	m := 3                    //直近のpingの応答時間を何個とるか
	t := 2                    //pingの平均応答時間

	var timeOutLogData [100][3]string //タイムアウトしているサーバーのタイムスタンプとアドレス,連続でタイムアウトした回数を格納する配列
	serverCount := 0                  //タイムアウトしているサーバーの数
	serverChecker := false            //同じアドレスから連続でタイムアウトが出ていないかの確認するための変数

	logData := make([][]string, 100) //タイムアウトしていないときのサーバーのアドレス、タイムスタンプ、pingの応答時間を格納する配列
	for i := range logData {
		logData[i] = make([]string, 2*m+1)
	}

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

		/*
		 *タイムアウトや故障に関する部分
		 */
		for i := 0; i < serverCount; i++ { //同じアドレスから連続でタイムアウトが出ていないかの確認
			serverChecker = false

			if arry1[1] == timeOutLogData[i][1] {
				serverChecker = true //timeOutLogData[][]に初期値を入れるかどうかの判定

				if arry1[2] == "-" { //タイムアウトした回数の計測と保存
					timeOutCount, _ := strconv.Atoi(timeOutLogData[i][2])
					timeOutCount++
					timeOutLogData[i][2] = strconv.Itoa(timeOutCount)
				}

				break
			}
		}

		if arry1[2] == "-" && serverChecker == false { //タイムアウトが出たサーバーのタイムスタンプとアドレス、タイムアウトの回数をlogData[]に代入している(初期化)
			serverCount++
			timeOutCount := 1
			str := strconv.Itoa(timeOutCount)
			timeOutLogData[serverCount-1][0] = arry1[0]
			timeOutLogData[serverCount-1][1] = arry1[1]
			timeOutLogData[serverCount-1][2] = str
		}

		for i := 0; i < serverCount; i++ {
			timeOutCount, _ := strconv.Atoi(timeOutLogData[i][2])

			if arry1[1] == timeOutLogData[i][1] && timeOutCount >= N { //サーバーが復帰したタイミングでアドレスと故障時間をコンソールに出力する

				if arry1[2] != "-" {
					time1, _ := strconv.Atoi(arry1[0])
					time2, _ := strconv.Atoi(timeOutLogData[i][0])
					str := fmt.Sprintf("%014d", time1-time2)

					fmt.Println(timeOutLogData[i][1] + " Server failure")
					fmt.Println(timeOutLogData[i][1] + " Server failure period is " + str)

					timeOutLogData[i][0] = "" //タイムアウトしていたサーバーのタイムスタンプの初期化
					timeOutLogData[i][1] = "" //タイムアウトしていたサーバーのアドレスの初期化
					timeOutLogData[i][2] = "0"
					serverChecker = false
				} else {
					fmt.Println(timeOutLogData[i][1] + " Server failure")
				}
			}

		}

		/*
		 *サーバーの可負荷に関する部分
		 */
		if arry1[2] != "-" {
			for i := range logData {
				if logData[i][0] != arry1[1] { //始めて出てきたタイムアウトしていないサーバーのアドレス、タイムスタンプ、応答時間を格納する
					if logData[i][0] == "" {
						logData[i][0] = arry1[1]   //アドレスの格納
						logData[i][m] = arry1[0]   //タイムスタンプの格納
						logData[i][2*m] = arry1[2] //応答時間の格納
						break
					} else {
						continue
					}
				} else { //すでに出てきたタイムアウトしていないサーバーのアドレス、タイムスタンプ、応答時間を並び替えて格納する
					for j := 0; j < m; j++ { //タイムスタンプの並べ替え
						if m-(j+1) == 0 {
							for k := 1; k < j+1; k++ {
								logData[i][m-(j+1)+k] = logData[i][m-(j+1)+k+1]
							}
							logData[i][m] = arry1[0]
							break
						} else if logData[i][m-(j+1)] == "" {
							for k := 0; k < j+1; k++ {
								logData[i][m-(j+1)+k] = logData[i][m-(j+1)+k+1]
							}
							logData[i][m] = arry1[0]
							break
						} else {
							continue
						}
					}

					for j := 0; j < m; j++ { //応答時間の並べ替え
						if 2*m-(j+1) == m {
							for k := 1; k < j+1; k++ {
								logData[i][2*m-(j+1)+k] = logData[i][2*m-(j+1)+k+1]
							}
							logData[i][2*m] = arry1[2]
							break
						} else if logData[i][2*m-(j+1)] == "" {
							for k := 0; k < j+1; k++ {
								logData[i][2*m-(j+1)+k] = logData[i][2*m-(j+1)+k+1]
							}
							logData[i][2*m] = arry1[2]
							break
						} else {
							continue
						}
					}

					if logData[i][1+m] != "" { //m回分のデータが溜まった時の処理
						sec := 0
						for j := m + 1; j < 2*m+1; j++ {
							s, _ := strconv.Atoi(logData[i][j])
							sec += s
						}
						rTime := sec / m //平均応答時間の算出
						if rTime >= t {
							time1, _ := strconv.Atoi(logData[i][1])
							time2, _ := strconv.Atoi(logData[i][m])
							fmt.Printf(logData[i][0]+" Server's overload period is "+"%014d\n", time2-time1)
						}

						break
					}
				}

			}
		}
	}

	if err = scanner.Err(); err != nil { //読み取りでエラーが出たときのエラー処理
		fmt.Println("scan error")
	}

}
