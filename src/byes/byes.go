package byes

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"log"
)

type ByePlayer struct {
	DCI       int64
	FirstName string
	LastName  string
	Byes      int
}

func LoadByes(file string) map[int64]ByePlayer {
	byeMap := make(map[int64]ByePlayer)
	excelFile, err := xlsx.OpenFile(file)
	if err != nil {
		fmt.Println(err)
		return byeMap
	}
	firstRecord := true
	for _, row := range excelFile.Sheets[0].Rows {
		if firstRecord {
			firstRecord = false
			continue
		}
		dci, err := row.Cells[0].Int64()
		if err != nil {
			log.Fatal(err)
		}
		byes, err := row.Cells[11].Int()
		if err != nil {
			log.Fatal(err)
		}
		byeMap[dci] = ByePlayer{
			DCI:       dci,
			FirstName: row.Cells[1].String(),
			LastName:  row.Cells[2].String(),
			Byes:      byes,
		}
	}
	return byeMap
}
