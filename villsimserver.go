/*
Doing right now:

Todo:
Make soldiers appear, and do stuff.

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
	baracks        map[string]int
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

type character struct {
	//Name and health will not be implemented for a bit, while this game is still an information transfer simulator, and does not have any real characters.
	name   string
	health int
	class  string
	//Currently, characters are just located wherever their struct is held, but maybe in the future I can add the ability for them to be in transit between areas.
	//This probably doesn't need to be a coardinate, as a information game, you should be unaware of their location until they interact with someone. The location can be an estimate done on the client side.

	//I do not think these characters need an "id." Having a uuid for their location in a map should be enough.
}

type action struct {
	name     string
	duration int
}

var playerList = make(map[string]player)
var locationList = make(map[string]location)
var characterList = make(map[string]character)
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

func getDifferenceFloat64(a float64, b float64) float64 {
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

func getObjectDistance(startingX float64, startingY float64, x float64, y float64) float64 {
	distanceX := getDifferenceFloat64(startingX, x)
	distanceY := getDifferenceFloat64(startingY, y)
	distance := math.Sqrt(math.Pow(float64(distanceX), 2) + math.Pow(float64(distanceY), 2))
	return distance

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
		if getLocationId(strings.TrimSpace(string(dArray[1]))) == "" {
			successful = false
			response += "no such location"
		} else {
			if currentPlayer.money > 99 {
				currentPlayer.baracks[getLocationId(strings.TrimSpace(string(dArray[1])))] += 1
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
			fmt.Println("Info is get")
			eventlocation := strings.TrimSpace(string(dArray[1]))
			eventcontent := strings.TrimSpace(string(dArray[2]))
			eventnewsworthiness, err := strconv.Atoi(strings.TrimSpace(string(dArray[3])))
			if err != nil {
				fmt.Println(err)
			} else {
				for key, location := range locationList {
					if location.name == eventlocation {
						fmt.Println("Location got", location.name)
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
	response := "{ "
	response += `"player": { `
	response += `"base": "` + locationList[currentPlayer.base].name + `", `
	response += `"money": ` + strconv.Itoa(currentPlayer.money) + ", "
	response += " } "
	response += "}"
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

	currentPlayer.baracks = make(map[string]int)

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

		/*
			actions := getActions(connId)

		*/
		updates = updateClient(connId)

		//message = response + actions + renderUpdates + "\n"
		message = `{ "updates": ` + updates + ", " + `"command": ` + response + "  }" + "\n"
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
			currentPlayer.money += locationList[currentPlayer.base].population * locationList[currentPlayer.base].averageIncome * locationList[currentPlayer.base].tax / 100
			playerList[key] = currentPlayer
		}
		turn = turn + 1
	}
}

func GenerateWorld() {
	village1 := genUUID()
	locationList[village1] = location{name: "Random-Village", class: "hub", population: 200, averageIncome: 1, tax: 30}

	village2 := genUUID()
	locationList[village2] = location{name: "Small-Town", class: "hub", population: 700, averageIncome: 1, tax: 30}

	pathid := genUUID()
	locationList[pathid] = location{name: "Somewhat-popular-road", class: "path", frequency: 4, start: village1, end: village2, distance: 20}

	village3 := genUUID()
	locationList[village3] = location{name: "Far-Away-Town", class: "hub", population: 1000, averageIncome: 1, tax: 30}

	pathid = genUUID()
	locationList[pathid] = location{name: "More-popular-road", class: "path", frequency: 8, start: village3, end: village1, distance: 30}

	village4 := genUUID()
	locationList[village4] = location{name: "A-Fork-Village", class: "hub", population: 300, averageIncome: 1, tax: 30}

	pathid = genUUID()
	locationList[pathid] = location{name: "rainbow-road", class: "path", frequency: 2, start: village4, end: village1, distance: 10}
}

func main() {
	fmt.Println("Generating world...")
	GenerateWorld()
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
