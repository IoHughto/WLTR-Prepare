package byes

import (
	"data"
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

func UpdateByes(people []data.Person, byeMap map[int64]ByePlayer) []data.Person {
	for i := 0; i < len(people); i++ {
		validator := data.DCIValidator{DCI: people[i].DCI}.Init()
		people[i].Byes = getMaxByes(validator, byeMap)
	}

	return people
}

func getMaxByes(validator data.DCIValidator, byeMap map[int64]ByePlayer) int {
	var byeList []int
	for _, dci := range validator.ValidDCIs {
		if byeMap[dci] != (ByePlayer{}) {
			byeList = append(byeList, byeMap[dci].Byes)
		}
	}
	return max(byeList)
}

func max(ints []int) int {
	if len(ints) == 0 {
		return 0
	}
	returnInt := ints[0]
	for _, tmpInt := range ints {
		if tmpInt > returnInt {
			returnInt = tmpInt
		}
	}
	return returnInt
}
