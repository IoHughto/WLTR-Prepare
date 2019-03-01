package main

import (
	"byes"
	"data"
	"encoding/json"
	"fileReader"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {

	newPtr := flag.Bool("new", false, "Only consider new players")
	dataPtr := flag.String("data", "./local.json", "File containing stored data")
	byesPtr := flag.String("byes", "", "CSV file with byes")
	addPtr := flag.String("add", "", "New file to add")
	savePtr := flag.Bool("save", false, "Save file")
	outPtr := flag.String("out", "./local.json", "Output file")
	printPtr := flag.Bool("print", false, "Print players in CSV format")
	obfuscatePtr := flag.Bool("obfuscate", false, "Print obfuscated info instead of full info")
	dropPtr := flag.String("drop", "", "Comma-separated list (no spaces!) of DCI numbers to drop")
	flag.Parse()

	var people []data.Person
	var stored []data.Person
	var byeMap map[int64]byes.ByePlayer
	if *addPtr != "" {
		people = fileReader.ReadFile(*addPtr)
	}
	stored = data.LoadData(*dataPtr)
	allPeople, newPeople := data.CombinePeople(people, stored)
	fmt.Println("All players:", len(allPeople))
	fmt.Println("New players:", len(newPeople))

	if *byesPtr != "" {
		byeMap = byes.LoadByes(*byesPtr)
		allPeople = byes.UpdateByes(allPeople, byeMap)
	}

	if *dropPtr != "" {
		allPeople, newPeople = data.DropPlayers(*dropPtr, allPeople, newPeople)
	}

	if *printPtr {
		if *newPtr {
			data.WriteCSV(newPeople, *obfuscatePtr)
		} else {
			data.WriteCSV(allPeople, *obfuscatePtr)
		}
	}
	if *savePtr {
		peopleJson, _ := json.Marshal(allPeople)
		jsonFile, err := os.Create(*outPtr)
		if err != nil {
			log.Fatal(err)
		}
		_, err = jsonFile.Write(peopleJson)
		if err != nil {
			log.Fatal(err)
		}
	}
	if len(allPeople) > 0 {
		data.PrintDuplicates(allPeople)
		fmt.Println(len(allPeople), "players")
		fmt.Println(len(people), "news")
		fmt.Println(len(stored), "old")
	}
}
