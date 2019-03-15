package playerImporter

import (
	"bufio"
	"data"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type wltrPerson struct {
	FirstName string
	LastName  string
	DCI       int64
	Country   string
	Status    string
	Byes      int
}

func ReadWLTRFile(file string) []data.Person {
	people := importWLTRPlayers(file)
	newPeople := cleanWLTR(people)

	return newPeople
}

func importWLTRPlayers(file string) []wltrPerson {
	csvFile, _ := os.Open(file)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for i := 0; i < 2; i++ {
		_, e := reader.Read()
		if e != nil {
			log.Fatal(e)
		}
	}
	var people []wltrPerson
	for {
		line, e := reader.Read()
		if e == io.EOF {
			break
		} else if e != nil {
			log.Fatal(e)
		}
		names := strings.Split(line[0], ", ")
		dci, e := strconv.ParseInt(line[1], 10, 64)
		if e != nil {
			log.Fatal(e)
		}
		byes := 0
		if line[6] != "-" {
			byes, e = strconv.Atoi(line[6])
		}
		if e != nil {
			log.Fatal(e)
		}
		people = append(people, wltrPerson{
			FirstName: strings.Title(strings.ToLower(names[1])),
			LastName:  strings.Title(strings.ToLower(names[0])),
			DCI:       dci,
			Country:   line[2],
			Status:    line[3],
			Byes:      byes,
		})
	}
	e := csvFile.Close()
	if e != nil {
		log.Fatal(e)
	}
	return people
}

func cleanWLTR(people []wltrPerson) []data.Person {
	var newPeople []data.Person
	for _, person := range people {
		newPeople = append(newPeople, data.Person{
			FirstName:    person.FirstName,
			LastName:     person.LastName,
			DCI:          person.DCI,
			Country:      person.Country,
			Status:       person.Status,
			Role:         "Player",
			Byes:         person.Byes,
			RawFirstName: person.FirstName,
			RawLastName:  person.LastName,
			RawDCI:       person.DCI,
			Dropped:      false,
		})
	}
	cleanedPeople, _ := data.CleanPeople(newPeople)
	return cleanedPeople
}
