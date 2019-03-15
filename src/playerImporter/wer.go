package playerImporter

import (
	"data"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type localPlayers struct {
	XMLName xml.Name    `xml:"LocalPlayers"`
	Player  []werPlayer `xml:"Player"`
}

type werPlayer struct {
	XMLName       xml.Name `xml:"Player"`
	FirstName     string   `xml:"FirstName,attr"`
	LastName      string   `xml:"LastName,attr"`
	MiddleInitial string   `xml:"MiddleInitial,attr"`
	DCI           string   `xml:"DciNumber,attr"`
	Country       string   `xml:"CountryCode,attr"`
	IsJudge       string   `xml:"IsJudge,attr"`
}

func CheckWER(file string) bool {
	handle, _ := os.Open(file)

	byteValue, e := ioutil.ReadAll(handle)
	if e != nil {
		return false
	}

	var players localPlayers

	e = xml.Unmarshal(byteValue, &players)
	if e != nil {
		return false
	}

	if len(players.Player) == 0 {
		return false
	}
	return true
}

func ReadWERFile(file string) []data.Person {
	people := importWERPlayers(file)
	newPeople := cleanWER(people)

	return newPeople
}

func cleanWER(people localPlayers) []data.Person {
	var newPeople []data.Person
	for _, person := range people.Player {
		dci, e := strconv.ParseInt(person.DCI, 10, 64)
		if e != nil {
			log.Fatal(e)
		}
		newPeople = append(newPeople, data.Person{
			FirstName:    strings.Title(strings.ToLower(person.FirstName)),
			LastName:     strings.Title(strings.ToLower(person.LastName)),
			DCI:          dci,
			Country:      person.Country,
			Status:       "Enrolled",
			Role:         "Player",
			Byes:         0,
			RawFirstName: strings.Title(strings.ToLower(person.FirstName)),
			RawLastName:  strings.Title(strings.ToLower(person.LastName)),
			RawDCI:       dci,
			Dropped:      false,
		})
	}
	cleanedPeople, _ := data.CleanPeople(newPeople)
	return cleanedPeople
}

func importWERPlayers(file string) localPlayers {
	xmlFile, _ := os.Open(file)

	byteValue, e := ioutil.ReadAll(xmlFile)
	if e != nil {
		fmt.Println(e)
	}

	var players localPlayers

	e = xml.Unmarshal(byteValue, &players)
	if e != nil {
		fmt.Println(e)
	}

	return players
}
