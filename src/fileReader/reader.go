package fileReader

import (
	"data"
	"fmt"
	"log"
	"os"
	"playerImporter"
)

func ReadFile(file string) []data.Person {
	var people []data.Person
	handle, e := os.Open(file)
	if e != nil {
		fmt.Println("File", file, "doesn't exist")
	} else {
		if handle.Close() != nil {
			log.Fatal(e)
		}
		if playerImporter.CheckCSV(file, 12) {
			people = playerImporter.ReadCFBEFile(file)
		} else if playerImporter.CheckWER(file) {
			people = playerImporter.ReadWERFile(file)
		} else if playerImporter.CheckCSV(file, 9) {
			people = playerImporter.ReadWLTRFile(file)
		}
	}
	return people
}
