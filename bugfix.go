package main

import (
	"fmt"
	"github.com/google/uuid"
)

type player struct {
	knownLocations []string
	base           string
	newsFeed       []article
	money          int
	movingSoldiersList []movingSoldiers
}

type location struct {
	name          string
	class         string
	information   int
	members       []string
	population    int
	averageIncome int
	tax           int
	frequency     float64
	start         string
	end           string
	events        []event
	distance      float64
	baracks       map[string]int
	soldiers      map[string]int
}

type movingSoldiers struct {
	destination string
	origin string
	population int
	distance float64
}

type article struct {
	title   string
	content string
	date    string
	id      string
}

type event struct {
	newsworthiness int
	title          string
	content        string
	id             string
	time           int
}

var locationList = make(map[string]location)

var turn int = 1

var testPlayer player

func deleteElement(list []string, i int) []string {
	return append(list[:i], list[i+1:]...)
}

func genUUID() string {
	id := uuid.New()
	return id.String()
}

func getLocationId(name string) string {
	for key, currentLocation := range locationList {
		if currentLocation.name == name {
			return key
		}
	}
	return ""
}

func getLocationName(id string) string {
	for key, currentLocation := range locationList {
		if key == id {
			return currentLocation.name
		}
	}
	return ""
}

func makeSoldiersTravel() {
	index := 0
	currentLocation := locationList[getLocationId(testPlayer.movingSoldiersList[index].destination)]
	currentLocation.soldiers["testPlayer"] += testPlayer.movingSoldiersList[index].population
	locationList[getLocationId(testPlayer.movingSoldiersList[index].destination)] = currentLocation
}

func indexLocation(name string, class string, population int, frequency int, averageIncome int, tax int, start string, end string, distance int) {
	locationId := genUUID()
	if name == "" {
		name = "Unamed-Village"
	}
	if class == "" {
		class = "hub"
	}
	if class != "path" {
		if population == 0 {
			fmt.Println("Village", name, "created without population")
		}
		if averageIncome == 0 {
			fmt.Println("Village", name, "created without average income")
		}
		if tax == 0 {
			fmt.Println("Village", name, "created without tax")
		}
	}
	locationList[locationId] = location{name: name, class: class, population: population, averageIncome: averageIncome, tax: tax, start: start, end: end}
	
	currentLocation := locationList[locationId]
	currentLocation.baracks = make(map[string]int)
	currentLocation.soldiers = make(map[string]int)
	locationList[locationId] = currentLocation
}

func generateWorld() {
	indexLocation("Random-Village", "hub", 200, 0, 1, 30, "", "", 0)
	indexLocation("Small-Town", "hub", 700, 0, 1, 30, "", "", 0)
	indexLocation("Somewhat-popular-road", "path", 0, 4, 1, 30, getLocationId("Random-Village"), getLocationId("Small-Town"), 20)
	indexLocation("Far-Away-Town", "hub", 1000, 0, 3, 20, "", "", 0)
	indexLocation("More-popular-road", "path", 0, 8, 0, 0, getLocationId("Far-Away-Town"), getLocationId("Random-Village"), 30)
	indexLocation("A-Fork-Village", "hub", 300, 0, 1, 30, "", "", 0)
}

func main() {
	fmt.Println("Generating world...")
	generateWorld()
	fmt.Println("Finished generating world")
}
