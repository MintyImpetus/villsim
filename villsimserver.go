/*
Doing right now:

Todo:
Make soldiers appear, and do stuff.

Add time to the soldiers and baracks. Eg, make it take time to build the baracks and hear back if they were successful.

Possibly make it so errors returned in json can be any string depicting the error, so clients can just output it.

Multithread the math done by attemptInfoTransfer to speed up gane loop (not important, it is only basic maths.)
*/

package main

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"math"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
	// "encoding/json"
)

type player struct {
	knownLocations []string
	base           string
	newsFeed       []article
	money          int
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

var playerList = make(map[string]player)
var locationList = make(map[string]location)
var connList = make(map[string]net.Conn)

var increment float64 = 0.01

var turn int = 1

func deleteElement(list []string, i int) []string {
	return append(list[:i], list[i+1:]...)
}

func getLocationId(name string) string {
	for key, currentLocation := range locationList {
		if currentLocation.name == name {
			return key
		}
	}
	return ""
}

func inRangeOfNumbers(query float64, low float64, high float64) bool {
	if query >= low && query <= high {
		return true
	} else {
		return false
	}
}

func getDifference(a int, b int) int {
	if a > b {
		return a - b
	} else {
		return b - a
	}
}

func genUUID() string {
	id := uuid.New()
	return id.String()
}

func getDaySinceGenesis() string {
	day := int(math.Floor(float64(turn) / 120))
	hour := int(math.Floor(math.Mod(float64(turn), 120)) / 5)
	return strconv.Itoa(day) + "-" + strconv.Itoa(hour)
}

func handleActions(connId string, dArray []string) (string, bool) {
	toClose := false
	response := "{ "

	if strings.TrimSpace(string(dArray[0])) == "" {
		response += `"result": "empty"`
	} else if strings.TrimSpace(string(dArray[0])) == "echo" {
		response += `"output": "`
		response += strings.TrimSpace(string(dArray[1]))
		response += `", `
		response += `"result": "success"`
	} else if strings.TrimSpace(string(dArray[0])) == "base" {
		response += `"output": "`
		response += locationList[playerList[connId].base].name
		response += `", `
		response += `"result": "success"`
	} else if strings.TrimSpace(string(dArray[0])) == "barack" {
		successful := true
		response += `"output": "`
		currentPlayer := playerList[connId]
		currentLocationId := getLocationId(strings.TrimSpace(string(dArray[1])))
		if getLocationId(strings.TrimSpace(string(dArray[1]))) == "" {
			successful = false
			response += "no such location"
		} else {
			if currentPlayer.money > 99 {
				//currentPlayer.baracks[getLocationId(strings.TrimSpace(string(dArray[1])))] += 1
				currentLocation := locationList[currentLocationId]
				currentLocation.baracks[connId] += 1
				locationList[currentLocationId] = currentLocation
				currentPlayer.money -= 99
			} else {
				successful = false
			}
		}
		playerList[connId] = currentPlayer
		response += ": " + strings.TrimSpace(string(dArray[1]))
		response += `", `
		if successful == true {
			response += `"result": "success"`
		} else {
			response += `"result": "fail"`
		}
	} else if strings.TrimSpace(string(dArray[0])) == "income" {
		response += `"output": "`
		response += strconv.Itoa(playerList[connId].money)
		response += `", `
		response += `"result": "success"`
	} else if strings.TrimSpace(string(dArray[0])) == "news" {
		response += `"output": [`
		for index, newsItem := range playerList[connId].newsFeed {
			response += ` "`
			response += newsItem.content + `"`
			if index < len(playerList[connId].newsFeed)-1 {
				response += ","
			}
			response += " "
		}
		response += "], "
		response += `"result": "success"`
	} else if strings.TrimSpace(string(dArray[0])) == "info" {
		if len(dArray) > 3 {
			eventlocation := strings.TrimSpace(string(dArray[1]))
			eventcontent := strings.TrimSpace(string(dArray[2]))
			eventnewsworthiness, err := strconv.Atoi(strings.TrimSpace(string(dArray[3])))
			if err != nil {
				fmt.Println(err)
			} else {
				for key, location := range locationList {
					if location.name == eventlocation {
						infoId := genUUID()
						currentEvent := event{newsworthiness: eventnewsworthiness, content: eventcontent, id: infoId, time: turn}
						location.events = append(location.events, currentEvent)
						locationList[key] = location
						fmt.Println("Added news item", currentEvent, "to", eventlocation)
						break
					}
				}
			}
		}
	} else if strings.TrimSpace(string(dArray[0])) == "exit" {
		fmt.Println("Player " + connId + " has left the game.")
		response += `"result": "success"`
		toClose = true
	} else if strings.TrimSpace(string(dArray[0])) == "list" {
		response += `"output": [`
		for index, currentLocationId := range playerList[connId].knownLocations {
			currentLocation := locationList[currentLocationId]
			fmt.Println(currentLocation)
			response += ` "`
			response += currentLocation.name + `"`
			if index < len(playerList[connId].knownLocations)-1 {
				response += ","
			}
			response += " "
		}
		response += "], "
		response += `"result": "success"`
	} else {
		fmt.Println("Command " + strings.TrimSpace(string(dArray[0])) + " not recognised from player " + connId)
		response += `"result": "invalid"`
	}
	response = response + " }"
	return response, toClose
}

func updateClient(connId string) string {
	currentPlayer := playerList[connId]
	response := `"player": { `
	response += `"base": "` + locationList[currentPlayer.base].name + `", `
	response += `"money": ` + strconv.Itoa(currentPlayer.money) + " "
	response += "}, "
	response += `"time": "` + getDaySinceGenesis() + `" `
	return response
}

func handleConnections(connId string) {

	var message string

	var response string
	var updates string
	//var renderUpdates string

	toClose := false

	numberOfHubs := 0
	for _, location := range locationList {
		if location.class == "hub" {
			numberOfHubs += 1
		}
	}

	randomBase := int(rand.Intn(numberOfHubs))

	i := 0
	currentPlayer := playerList[connId]
	for key, _ := range locationList {
		currentPlayer.knownLocations = append(currentPlayer.knownLocations, key)
		if locationList[key].class == "hub" {
			if i == randomBase {
				currentPlayer.base = key
			}
			i += 1
		}
	}

	currentPlayer.newsFeed = append(currentPlayer.newsFeed, article{content: "Some crazy news"})

	playerList[connId] = currentPlayer

	fmt.Println("Player " + connId + " created")
	connList[connId].Write([]byte("Connection accepted, account " + connId + " created.\n"))

	for {
		//Get what the player wants to do and then send a response.
		message = ""

		response = ""
		updates = ""

		data, err := bufio.NewReader(connList[connId]).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		dArray := strings.Split(data, " ")
		response, toClose = handleActions(connId, dArray)

		updates = updateClient(connId)

		message = `{ "updates": { ` + updates + "}, " + `"command": ` + response + " }" + "\n"
		connList[connId].Write([]byte(message))

		if toClose {
			connList[connId].Close()
			delete(connList, connId)
			delete(playerList, connId)
			return
		}
	}
}

func attemptInfoTransfer(theevent event, place location, origin string, destination string) {
	timesince := turn - theevent.time
	rotations := float64(timesince) / 4
	chance := float64(theevent.newsworthiness) * math.Pow(rotations, float64(theevent.newsworthiness)) * place.frequency * float64(locationList[destination].population)
	chance /= place.distance
	chance /= 100
	chance += float64(rotations / 2)
	randomNumber := rand.Intn(int(chance) + 70)
	fmt.Println(chance)
	fmt.Println(randomNumber)
	if float64(randomNumber) < chance {
		fmt.Println("Event information transfered from", locationList[origin].name, "to", locationList[destination].name)
		currentLocation := locationList[destination]
		theevent.time = turn
		currentLocation.events = append(currentLocation.events, theevent)
		locationList[destination] = currentLocation
		fmt.Println(locationList[destination].events)
	}
}

func gameLoop() {
	for {
		time.Sleep(time.Second)
		for _, place := range locationList {
			if place.class == "path" {
				for _, startevent := range locationList[place.start].events {
					transfered := false
					for _, endevent := range locationList[place.end].events {
						if endevent.id == startevent.id {
							transfered = true
							break
						}
					}
					if transfered == false {
						attemptInfoTransfer(startevent, place, place.start, place.end)
					}
				}
			}
		}
		for _, place := range locationList {
			if place.class == "path" {
				for _, endevent := range locationList[place.end].events {
					transfered := false
					for _, startevent := range locationList[place.start].events {
						if startevent.id == endevent.id {
							transfered = true
							break
						}
					}
					if transfered == false {
						attemptInfoTransfer(endevent, place, place.end, place.start)
					}
				}
			}
		}
		for key, currentPlayer := range playerList {
			for _, locationevent := range locationList[currentPlayer.base].events {
				alreadyseen := false
				for _, currentArticle := range currentPlayer.newsFeed {
					if currentArticle.id == locationevent.id {
						alreadyseen = true
					}
				}
				if alreadyseen == false {
					currentPlayer.newsFeed = append(currentPlayer.newsFeed, article{title: locationevent.title, content: locationevent.content, date: "Add a function to find the time, or something.", id: locationevent.id})
					playerList[key] = currentPlayer
				}
			}
		}
		for key, currentPlayer := range playerList {
			currentPlayer.money += locationList[currentPlayer.base].population * locationList[currentPlayer.base].averageIncome * locationList[currentPlayer.base].tax / 100 / 10
			playerList[key] = currentPlayer
		}
		for key, currentLocation := range locationList {
			for key, barackNumber := range currentLocation.baracks {
				if barackNumber > 0 {
					currentLocation.soldiers[key] = barackNumber * currentLocation.population / 100
				}
			}
			locationList[key] = currentLocation
		}
		turn = turn + 1
	}
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
	PORT := ":9876"
	dstream, err := net.Listen("tcp", PORT)
	fmt.Println("Started listener on port", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dstream.Close()
	go gameLoop()
	for {
		conn, err := dstream.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			fmt.Println("Connection received")
			connId := genUUID()
			connList[connId] = conn
			go handleConnections(connId)
		}
	}
}
