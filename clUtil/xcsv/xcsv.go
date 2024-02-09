package xcsv

import (
	"bytes"
	"encoding/csv"
)

func ParseXCSV(header []string, rows [][]string) []byte {
	//		wr.WriteString("\xEF\xBB\xBF")
	b := &bytes.Buffer{}
	wr := csv.NewWriter(b)
	wr.Write(header)
	for _, row := range rows {
		wr.Write(row)
	}
	wr.Flush()
	return b.Bytes()
}
