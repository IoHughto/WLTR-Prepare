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

type cfbePerson struct {
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	DCI               int64  `json:"dci"`
	Country           string `json:"country"`
	Status            string `json:"status"`
	Role              string `json:"role"`
	CreatedAt         string `json:"created_at"`
	LastModified      string `json:"last_modified"`
	RegistrationState string `json:"registration_state"`
}

func ReadCFBEFile(file string) []data.Person {
	people := importCFBEPlayers(file)
	newPeople := cleanCFBE(people)

	return newPeople
}

func CheckCSV(file string, rows int) bool {
	handle, _ := os.Open(file)
	reader := csv.NewReader(bufio.NewReader(handle))
	line, e := reader.Read()
	if e != nil {
		return false
	}
	if len(line) != rows {
		return false
	}
	return true
}

func cleanCFBE(people []cfbePerson) []data.Person {
	var newPeople []data.Person
	for _, person := range people {
		if person.RegistrationState == "registered" {
			newPeople = append(newPeople, data.Person{
				FirstName: person.FirstName,
				LastName:  person.LastName,
				DCI:       person.DCI,
				Country:   person.Country,
				Status:    person.Status,
				Role:      person.Role,
				Byes:      0,
			})
		}
	}
	cleanedPeople, _ := data.CleanPeople(newPeople)
	return cleanedPeople
}

func importCFBEPlayers(file string) []cfbePerson {
	csvFile, _ := os.Open(file)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	_, e := reader.Read()
	if e != nil {
		log.Fatal(e)
	}
	var people []cfbePerson
	for {
		line, e := reader.Read()
		if e == io.EOF {
			break
		} else if e != nil {
			log.Fatal(e)
		}
		dci, e := strconv.ParseInt(line[2], 10, 64)
		if e != nil {
			log.Fatal(e)
		}
		people = append(people, cfbePerson{
			FirstName:         strings.Title(line[0]),
			LastName:          strings.Title(line[1]),
			DCI:               dci,
			Country:           line[3],
			Status:            line[4],
			Role:              line[5],
			CreatedAt:         line[8],
			LastModified:      line[9],
			RegistrationState: line[10],
		})
	}
	e = csvFile.Close()
	if e != nil {
		log.Fatal(e)
	}
	return people
}
