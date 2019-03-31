package data

import (
	"byes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Tournament struct {
	Players []Person                 `json:"players"`
	Judges  []Person                 `json:"judges"`
	Byes    map[int64]byes.ByePlayer `json:"byes"`
}

type Bye struct {
	DCI  int64 `json:"dci"`
	Byes int   `json:"byes"`
}

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

func UpdateByes(tournament Tournament) Tournament {
	for player := range tournament.Players {
		validator := DCIValidator{DCI: tournament.Players[player].DCI}.Init()
		tournament.Players[player].Byes = GetMaxByes(validator, tournament.Byes)
	}
	return tournament
}

func GetMaxByes(validator DCIValidator, byeMap map[int64]byes.ByePlayer) int {
	var byeList []int
	for _, dci := range validator.ValidDCIs {
		if byeMap[dci] != (byes.ByePlayer{}) {
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
		for j, otherPerson := range people {
			if i < j {
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

func (tournament Tournament) PrintSummary() {

	activePlayerCount := 0
	inactivePlayerCount := 0
	judgeCount := len(tournament.Judges)

	for _, player := range tournament.Players {
		if player.Dropped {
			inactivePlayerCount++
		} else {
			activePlayerCount++
		}
	}

	fmt.Println("Summary")
	fmt.Printf("%-18v%v\n", "Active Players:", activePlayerCount)
	fmt.Printf("%-18v%v\n", "Dropped Players:", inactivePlayerCount)
	fmt.Printf("%-18v%v\n", "Judges:", judgeCount)

}

func (tournament Tournament) WriteCSV(fileName string, obfuscate bool) {
	if len(tournament.Players) == 0 && len(tournament.Judges) == 0 {
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
		for _, person := range tournament.Players {
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
		var t reflect.Type
		if len(tournament.Players) > 0 {
			t = reflect.TypeOf(tournament.Players[0])
		} else {
			t = reflect.TypeOf(tournament.Judges[0])
		}
		names := make([]string, t.NumField())
		for i := range names {
			names[i] = t.Field(i).Name
		}
		err := w.Write(names)
		check(err)

		for _, record := range tournament.Judges {
			if !record.Dropped {
				err := w.Write(record.ValueStrings())
				check(err)
			}
		}

		for _, record := range tournament.Players {
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

func LoadData(file string) Tournament {
	var tournament Tournament

	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return tournament
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return tournament
	}

	err = json.Unmarshal(byteValue, &tournament)
	if err != nil {
		fmt.Println(err)
	}

	return tournament
}

func DropPlayers(tournament Tournament, dropList string) Tournament {
	dciNumbersToDrop := parseDropList(dropList)
	for _, drop := range dciNumbersToDrop {
		dropIndex := playerIndexExists(tournament.Players, drop)
		if dropIndex != -1 {
			tournament.Players[dropIndex].Dropped = true
		} else {
			fmt.Println("DCI", drop, "not in event")
		}
	}
	return tournament
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

func AddPlayers(newPeople []Person, previousPlayers []Person) []Person {

	for _, newPerson := range newPeople {
		if !Contains(previousPlayers, newPerson) {
			previousPlayers = append(previousPlayers, newPerson)
		}
	}

	return previousPlayers
}
