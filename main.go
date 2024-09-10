package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"
)

var version string

func main() {
	var v bool
	flag.BoolVar(&v, "version", false, "バージョン")
	flag.Parse()

	// バージョン表示
	if v {
		fmt.Printf("Todoner.exe version %s\n", version)
		return
	}

	// 起動時引数のチェック
	if len(flag.Args()) != 1 {
		fmt.Println("Usage: Todoner.exe <input_file>")
		return
	}

	// 入力ファイルのオープン
	fileName := flag.Arg(0)
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// CSVリーダーの作成
	reader := csv.NewReader(file)

	// 1行目はヘッダなので読み飛ばす
	reader.Read()

	// 構造体にデータを格納
	var records []Record

	// 2行目以降の読み込みと処理
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading line:", err)
			return
		}
		if r, err := parseRecord(line); err != nil {
			fmt.Println("Error parse Record:", err)
			return
		} else {
			records = append(records, r)
		}
	}

	// 処理区でグルーピング
	grpKu := lop.GroupBy(records, func(r Record) string {
		return r.GetKu()
	})

	// 1) 処理区毎の使用件数
	// fmt.Println("1) 処理区毎の集計")
	for ku, rs := range grpKu {
		fmt.Printf("区=%s, 件数=%d\n", ku, len(rs))

		// 処理区＋異動事由でグルーピング
		grpCido := lop.GroupBy(rs, func(r Record) string {
			return r.LifeEvent
		})

		for ido, rs := range grpCido {
			totalUsageTime := lo.SumBy(rs, func(r Record) time.Duration {
				return r.UsageTime
			})
			fmt.Printf("区=%s, 異動事由=%s, 件数=%d, 平均使用時間=%v\n", ku, ido, len(rs), totalUsageTime/time.Duration(len(rs)))
		}
	}
}
