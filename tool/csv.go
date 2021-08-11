package main

import (
	"encoding/csv"
	"os"
)

// 追加一行
func WriteCsvFile(headline []string) {
	filename := "./data.csv"
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	if err != nil {
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	writer.Comma = ','
	writer.Write(headline)
	writer.Flush()
}
