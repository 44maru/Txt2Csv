package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const itemName = "商品名"
const color = "カラー"
const itemSize = "商品サイズ"
const buyerName = "購入者名"
const buyerMail = "購入者メールアドレス"
const orderNo = "注文オーダー番号"

var enc = simplifiedchinese.GBK

func main() {
	if len(os.Args) != 2 {
		failOnError("main.exeにテキストファイルをドラッグ&ドロップしてください", nil)
	}
	convertTxt2Csv(os.Args[1])
	waitEnter()
}

func failOnError(errMsg string, err error) {
	errs := errors.WithStack(err)
	fmt.Println(errMsg)
	fmt.Printf("%+v\n", errs)
	waitEnter()
	//panic(err)
	os.Exit(1)
	//log.Fatal(err)
}

func waitEnter() {
	fmt.Println("エンターを押すと処理を終了します。")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
}

func convertTxt2Csv(fileName string) {

	fp, err := os.Open(fileName)
	if err != nil {
		failOnError("ファイル読込に失敗しました", err)
	}
	defer fp.Close()

	file, err := os.OpenFile("./result.csv", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		failOnError("result.csvのオープンに失敗しました", err)
	}
	defer file.Close()

	err = file.Truncate(0) // ファイルを空っぽにする(実行2回目以降用)
	if err != nil {
		failOnError("CSVファイルの初期化に失敗しました", err)
	}

	isCsvItemLine := false
	isOrderNoLine := false
	var itemList []string
	orderNoMap := map[string]bool{}

	writer := csv.NewWriter(transform.NewWriter(file, japanese.ShiftJIS.NewEncoder()))

	r := transform.NewReader(fp, japanese.ShiftJIS.NewDecoder())
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {

		if strings.Contains(scanner.Text(), itemName) ||
			strings.Contains(scanner.Text(), color) ||
			strings.Contains(scanner.Text(), itemSize) ||
			strings.Contains(scanner.Text(), buyerName) ||
			strings.Contains(scanner.Text(), buyerMail) {

			isCsvItemLine = true

		} else if strings.Contains(scanner.Text(), orderNo) {
			isOrderNoLine = true

		} else if isCsvItemLine {
			itemList = append(itemList, scanner.Text())
			isCsvItemLine = false

		} else if isOrderNoLine {
			orderNo := scanner.Text()
			var isAlreadyExists bool = orderNoMap[orderNo]

			if isAlreadyExists {
				itemList = []string{}
				isOrderNoLine = false
				continue
			}

			itemList = append(itemList, orderNo)
			isOrderNoLine = false
			writer.Write(itemList)
			orderNoMap[orderNo] = true
			itemList = []string{}
		}
	}

	writer.Flush()
	fmt.Println("result.csvを出力しました")
}
