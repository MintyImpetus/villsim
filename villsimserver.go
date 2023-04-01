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
	knownLocations []location
	//When I work on updates, they will be a list of updates. This will be a string, but also with metadata about who told them that, and what group it is about.
	updates string
}

type location struct {
	information int
	members []character
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
var connList = make(map[string]net.Conn)

var increment float64 = 0.01

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

func handleActions(connId string, dArray []string) string {

	response := "[ { "
	currentPlayer := playerList[connId]
	
	if strings.TrimSpace(string(dArray[0])) == "echo" {
		response = response + " { " + strings.TrimSpace(string(dArray[1])) 
	}
	else if strings.TrimSpace(string(dArray[0])) == "exit" {
		fmt.Println("Player " + connId + " has left the game.")
		response += `"exit": "successful"`
		delete(playerList, connId)
	} else {
		fmt.Println("Command " + strings.TrimSpace(string(dArray[0])) + " not recognised from player " + connId)
	}
	response = response + " } ] "
	return response
}

func handleConnections(connId string) {

	var message string

	var response string
	//var renderUpdates string

	playerList[connId] = player{health: 20, x: 0, y: 1, renderDistance: 3}
	
//Send opening information to the player.

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
		response = handleActions(connId, dArray)
		
		/*
		actions := getActions(connId)

		renderUpdates = updateClient(connId)
		*/
		
		//message = response + actions + renderUpdates + "\n"
		message = response + "\n"
                connList[connId].Write([]byte(message))
        }
}

func gameLoop() {
	for {
		time.Sleep(time.Second)
	}
}

func main() {
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
			connId := genUUID()
			connList[connId] = conn
			go handleConnections(connId)
		}
	}
}
