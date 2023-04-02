/*
Doing right now:

TODO:

Complete the communication protocol and begin working on the client in pygame.
	Make it so the message variable is repeatedly updated to include what a player within render distance does, and every time a player updates, send them the response to their previous action, all the actions of visible players and then all visible block data.



Whenever action is done, run a goroutine function that adds the action to all nearby players list of actions. Have a goroutine running that constantly lowers the tick-till-end number for each of the actions and destroys them when done. 

Whenever the client asks for an update, send them all the actions and how long 'till they end.

At some point in the future, I could attempt to make all major tasks run through goroutines, as it may speed up the program.
Increased performance could also be gotten by adding more exit statements to loops.
*/

package main

import (
        "bufio"
        "fmt"
        "net"
        "strings"
	"strconv"
	"math"
	"time"
	"github.com/google/uuid"
//	"encoding/json"
)

type player struct {
	knownLocations []string
	//When I work on updates, they will be a list of updates. This will be a string, but also with metadata about who told them that, and what group it is about.
	updates string
	base string
	newsFeed []article
}

type location struct {
	name string
	class string
	information int
	members []string
	population int
	frequency float64
	start string
	end string
	events []event
	distance float64
}

type article struct {
	title string
	content string
	date string
	id string
}

type event struct {
	newsworthiness int
	title string
	content string
	id string
	time int
}

type character struct {
	//Name and health will not be implemented for a bit, while this game is still an information transfer simulator, and does not have any real characters.
	name string
	health int
	class string
	//Currently, characters are just located wherever their struct is held, but maybe in the future I can add the ability for them to be in transit between areas.
	//This probably doesn't need to be a coardinate, as a information game, you should be unaware of their location until they interact with someone. The location can be an estimate done on the client side.

	//I do not think these characters need an "id." Having a uuid for their location in a map should be enough.
}

type action struct {
	name string
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

func inRangeOfNumbers( query float64, low float64, high float64) bool {
	if (query >= low && query <= high) {
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
	response := "[ { "
	
	if strings.TrimSpace(string(dArray[0])) == "echo" {
		response += `"output": "`
		response += strings.TrimSpace(string(dArray[1]))
		response += `", `
		response += `"result": "success"`
	} else if strings.TrimSpace(string(dArray[0])) == "base" {
		response += `"output": "`
		response += locationList[playerList[connId].base].name
		response += `", `
		response += `"result": "success"`
	} else if strings.TrimSpace(string(dArray[0])) == "news" {
		response += `"output": [`
		for index, newsItem := range playerList[connId].newsFeed {
			response += ` "`
			response += newsItem.content + `"`
			if index < len(playerList[connId].newsFeed) - 1{
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
		for index, currentLocationId:= range playerList[connId].knownLocations {
			currentLocation := locationList[currentLocationId]
			fmt.Println(currentLocation)
			response += ` "`
			response += currentLocation.name + `"`
			if index < len(playerList[connId].knownLocations) - 1 {
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
	response = response + " } ] "
	return response, toClose
}

func handleConnections(connId string) {

	var message string

	var response string
	//var renderUpdates string

	toClose := false

	currentPlayer := playerList[connId]
	for key,_ := range locationList {
		currentPlayer.knownLocations = append(currentPlayer.knownLocations, key)
		if locationList[key].class == "hub" {
			currentPlayer.base = key
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
		//renderUpdates = ""

                data, err := bufio.NewReader(connList[connId]).ReadString('\n')
                if err != nil {
                        fmt.Println(err)
                        return
                }
		dArray := strings.Split(data, " ")
		response, toClose = handleActions(connId, dArray)
		
		/*
		actions := getActions(connId)

		renderUpdates = updateClient(connId)
		*/
		
		//message = response + actions + renderUpdates + "\n"
		message = response + "\n"
                connList[connId].Write([]byte(message))

		if toClose {
			connList[connId].Close()
			delete(connList, connId)
			delete(playerList, connId)
			return
		}
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
						fmt.Println("Not transfered", startevent)
						rotations := turn - startevent.time 
						rotations = rotations / 4
						chance := float64(startevent.newsworthiness) * math.Pow(float64(rotations), float64(startevent.newsworthiness)) * place.frequency * float64(locationList[place.end].population) / place.distance / 10000
						chance += float64(rotations / 2)
						fmt.Println(chance)
						if turn - startevent.time > 10 {
							currentLocation := locationList[place.end]
							currentLocation.events = append(currentLocation.events, startevent)
							locationList[place.end] = currentLocation
						}
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
						fmt.Println("Not transfered", endevent)
						if turn - endevent.time > 10 {
							currentLocation := locationList[place.start]
							currentLocation.events = append(currentLocation.events, endevent)
							locationList[place.start] = currentLocation
						}

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
		turn = turn + 1
	}
}

func GenerateWorld() {
	village1 := genUUID()
	locationList[village1] = location{name: "Random-Village", class: "hub", population: 200}

	village2 := genUUID()
	locationList[village2] = location{name: "Small-Town", class: "hub", population: 700}

	pathid := genUUID()
	locationList[pathid] = location{name: "Somewhat-popular-road", class: "path", frequency: 4, start: village1, end: village2}
}

func main() {
	GenerateWorld()

        PORT := ":9876"
        dstream, err := net.Listen("tcp", PORT)
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
