package data

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Person struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	DCI          int64  `json:"dci"`
	Country      string `json:"country"`
	Status       string `json:"status"`
	Role         string `json:"role"`
	Byes         int    `json:"byes"`
	RawFirstName string `json:"rawFirstName"`
	RawLastName  string `json:"rawLastName"`
	RawDCI       int64  `json:"rawDCI"`
	Dropped      bool   `json:"dropped"`
}

func firstNameChange(newPerson Person, oldPerson Person) bool {
	return newPerson.LastName == oldPerson.LastName &&
		newPerson.DCI == oldPerson.DCI &&
		newPerson.Country == oldPerson.Country &&
		newPerson.Status == oldPerson.Status &&
		newPerson.Role == oldPerson.Role &&
		newPerson.Byes <= oldPerson.Byes
}

func sameRawPeople(newPerson Person, oldPerson Person) bool {
	return newPerson.RawFirstName == oldPerson.RawFirstName &&
		newPerson.RawLastName == oldPerson.RawLastName &&
		newPerson.RawDCI == oldPerson.RawDCI
}

func Contains(oldPeople []Person, newPerson Person) bool {
	for _, oldPerson := range oldPeople {
		if oldPerson == newPerson {
			return true
		}
		if sameRawPeople(newPerson, oldPerson) {
			if newPerson.Byes > oldPerson.Byes {
				fmt.Println(newPerson, "is a bye dup")
				validator := DCIValidator{DCI: newPerson.DCI}
				validator.Init()
			}
			return true
		}
		if firstNameChange(newPerson, oldPerson) {
			return true
		}
	}
	return false
}

func isDCIMatch(person Person, check Person) bool {
	if person.DCI == check.DCI {
		return true
	}
	personValidator := DCIValidator{DCI: person.DCI}
	personValidator.Init()
	checkValidator := DCIValidator{DCI: check.DCI}
	checkValidator.Init()
	if (DCIValidator{DCI: person.DCI}).Contains(check.DCI) ||
		(DCIValidator{DCI: check.DCI}).Contains(person.DCI) {
		return true
	}
	return false
}

func PrintDuplicates(people []Person) {
	for i, person := range people {
		//for i := 0; i<len(people); i++ {
		for j, otherPerson := range people {
			if i < j {
				//for j := i+1; j<len(people); j++ {
				if otherPerson == person {
					fmt.Println("Complete duplicate", otherPerson, person, i, j)
				} else {
					if isDCIMatch(person, otherPerson) {
						fmt.Println("DCI duplicate", otherPerson, person)
					}
					if otherPerson.FirstName == person.FirstName && otherPerson.LastName == person.LastName {
						fmt.Println("Name duplicate", otherPerson, person)

					}
				}
			}
		}
	}
}

func CleanPeople(people []Person) ([]Person, []Person) {
	var cleanedPeople []Person
	var problemPeople []Person
	for _, person := range people {
		if Contains(cleanedPeople, person) {
			problemPeople = append(problemPeople, person)
		} else {
			cleanedPeople = append(cleanedPeople, person)
		}
	}
	if len(problemPeople) > 0 {
		fmt.Println(len(problemPeople), " problem people")
		for _, person := range problemPeople {
			fmt.Println(person)
		}
	}
	return cleanedPeople, problemPeople
}

func getObfuscatedRow(person Person) string {
	DCIString := fmt.Sprintf("%v", person.DCI)
	var DCI string
	if len(DCIString) < 4 {
		DCI = DCIString
	} else {
		DCI = DCIString[len(DCIString)-4:]
	}
	Byes := fmt.Sprintf("%v", person.Byes)

	return person.LastName + "," + person.FirstName + "," + DCI + "," + Byes
}

func WriteCSV(people []Person, fileName string, obfuscate bool) {
	if len(people) == 0 {
		fmt.Println("There are no records to print")
		return
	}
	var file *os.File
	var e error
	if fileName != "stdout" {
		file, e = os.Create(fileName)
		check(e)
	}
	if obfuscate {
		headerString := "Last name,First Name,Last 4 of DCI,Byes"
		if fileName != "stdout" {
			_, e = fmt.Fprintln(file, headerString)
			check(e)
		} else {
			fmt.Println(headerString)
		}
		for _, person := range people {
			if !person.Dropped {
				if fileName != "stdout" {
					_, e = fmt.Fprintln(file, getObfuscatedRow(person))
					check(e)
				} else {
					fmt.Println(getObfuscatedRow(person))
				}
			}
		}
	} else {
		var w *csv.Writer
		if fileName != "stdout" {
			w = csv.NewWriter(file)
		} else {
			w = csv.NewWriter(os.Stdout)
		}
		t := reflect.TypeOf(people[0])
		names := make([]string, t.NumField())
		for i := range names {
			names[i] = t.Field(i).Name
		}
		err := w.Write(names)
		check(err)

		for _, record := range people {
			if !record.Dropped {
				err := w.Write(record.ValueStrings())
				check(err)
			}
		}

		w.Flush()

		err = w.Error()
		check(err)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (person Person) ValueStrings() []string {
	v := reflect.ValueOf(person)
	ss := make([]string, v.NumField())
	for i := range ss {
		ss[i] = fmt.Sprintf("%v", v.Field(i))
	}
	return ss
}

func CombinePeople(new []Person, old []Person) ([]Person, []Person) {
	allPeople := old
	var newPeople []Person
	for _, person := range new {
		if !Contains(old, person) {
			allPeople = append(allPeople, person)
			newPeople = append(newPeople, person)
		}
	}
	return allPeople, newPeople
}

func LoadData(file string) []Person {
	var people []Person

	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return people
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return people
	}

	err = json.Unmarshal(byteValue, &people)
	if err != nil {
		fmt.Println(err)
	}

	return people
}

func DropPlayers(dropList string, all []Person, new []Person) ([]Person, []Person) {
	dciNumbersToDrop := parseDropList(dropList)
	for _, drop := range dciNumbersToDrop {
		allDropIndex := playerIndexExists(all, drop)
		newDropIndex := playerIndexExists(new, drop)
		if allDropIndex != -1 {
			all[allDropIndex].Dropped = true
		}
		if newDropIndex != -1 {
			new[newDropIndex].Dropped = true
		}
		if allDropIndex == -1 && newDropIndex == -1 {
			fmt.Println("DCI", drop, "not in event")
		}
	}
	return all, new
}

func parseDropList(dropList string) []int64 {
	dropStrings := strings.Split(dropList, ",")
	var drops []int64
	for _, drop := range dropStrings {
		newDCI, e := strconv.ParseInt(drop, 10, 64)
		check(e)
		drops = append(drops, newDCI)
	}
	return drops
}

func playerIndexExists(people []Person, dci int64) int {
	for index, person := range people {
		if dci == person.DCI {
			return index
		}
	}
	return -1
}
