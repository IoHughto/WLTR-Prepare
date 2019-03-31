package main

import (
	"byes"
	"data"
	"encoding/json"
	"fileReader"
	"flag"
	"log"
	"os"
)

type stringArray []string

func (i *stringArray) String() string {
	return "my string representation"
}

func (i *stringArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var addFiles stringArray
	flag.Var(&addFiles, "add", "List of files")
	dataPtr := flag.String("data", "./local.json", "File containing stored data")
	judgePtr := flag.String("judge", "", "File containing staff list")
	byesPtr := flag.String("byes", "", "CSV file with byes")
	savePtr := flag.Bool("save", false, "Save file")
	outPtr := flag.String("out", "./local.json", "Output file")
	printPtr := flag.String("print", "", "Print players to <file> in CSV format (stdout prints to terminal instead)")
	obfuscatePtr := flag.Bool("obfuscate", false, "Print obfuscated info instead of full info")
	dropPtr := flag.String("drop", "", "Comma-separated list (no spaces!) of DCI numbers to drop")
	flag.Parse()

	var tournament data.Tournament
	tournament = data.LoadData(*dataPtr)
	if len(addFiles) > 0 {
		for _, file := range addFiles {
			tournament.Players = data.AddPlayers(fileReader.ReadFile(file), tournament.Players)
		}
	}
	if *judgePtr != "" {
		tournament.Judges = data.AddPlayers(fileReader.ReadFile(*judgePtr), tournament.Judges)
	}

	if *byesPtr != "" {
		tournament.Byes = byes.LoadByes(*byesPtr)
		tournament = data.UpdateByes(tournament)
	}

	if *dropPtr != "" {
		tournament = data.DropPlayers(tournament, *dropPtr)
	}

	if *printPtr != "" {
		tournament.WriteCSV(*printPtr, *obfuscatePtr)
	}
	if *savePtr {
		tournamentJson, _ := json.Marshal(tournament)
		jsonFile, err := os.Create(*outPtr)
		if err != nil {
			log.Fatal(err)
		}
		_, err = jsonFile.Write(tournamentJson)
		if err != nil {
			log.Fatal(err)
		}
	}
	if len(tournament.Players) > 0 {
		data.PrintDuplicates(tournament.Players)
		tournament.PrintSummary()
	}
}
