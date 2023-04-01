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

type item struct {
	name string
	itemType string
	damage int
	ddq int
	description string
}

type player struct {
	health int
	inventory []item
	visibleActions []action
	x float64
	y float64
	renderDistance float64
}

type block struct {
	blockType string
	flickerUp bool
	sinceLastFlicker int
	x float64
	y float64
	endX float64
	endY float64
}

type entity struct {
	name string
	id string
	hp int
	x float64
	y float64
}

type action struct {
	name string
	duration int
}

var blockList = make(map[string]block)
var entityList = make(map[string]entity)
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

func checkIfBlock(x float64, y float64) bool {
	exists := false
	for _, currentObject := range blockList {
		if ( inRangeOfNumbers(x, currentObject.x, currentObject.endX) && inRangeOfNumbers(y, currentObject.y, currentObject.endY)) {
			exists = true
		}
	}
	return exists

}

func HandleMove(PlayerX float64, PlayerY float64, dArray []string) (float64, float64) {
	tmpX := PlayerX
	tmpY := PlayerY
	fmt.Println(PlayerX, PlayerY)
	if strings.TrimSpace(string(dArray[1])) == "up" {
		fmt.Println("up")
		tmpY = tmpY + 1
	}
	if strings.TrimSpace(string(dArray[1])) == "down" {
		fmt.Println("down")
		tmpY = tmpY - 1
	}
	if strings.TrimSpace(string(dArray[1])) == "left" {
		fmt.Println("left")
		tmpX = tmpX - 1
	}
	if strings.TrimSpace(string(dArray[1])) == "right" {
		fmt.Println("right")
		tmpX = tmpX + 1
	}
	blocked := checkIfBlock(tmpX, tmpY)
	if blocked == true {
		fmt.Println("Blocked")
		tmpX = PlayerX
		tmpY = PlayerY
	}
	fmt.Println(PlayerX, PlayerY)
	return tmpX, tmpY
}

func getObjectDistance(startingX float64, startingY float64, x float64, y float64) float64 {
	distanceX := getDifferenceFloat64(startingX, x)
	distanceY := getDifferenceFloat64(startingY, y)
	distance := math.Sqrt(math.Pow(float64(distanceX), 2) + math.Pow(float64(distanceY), 2))
	return distance

}

func displayAnimation(key string, name string, x float64, y float64, duration int) {
	for _, currentPlayer := range playerList {
		if getObjectDistance(x, y, currentPlayer.x, currentPlayer.y) <= playerList[key].renderDistance {
			currentPlayer.visibleActions = append(playerList[key].visibleActions, action{ name: name, duration: duration })
			playerList[key] = currentPlayer
		}
	}
}

func handleActions(connId string, dArray []string) string {

	response := "[ { "
	currentPlayer := playerList[connId]

	if strings.TrimSpace(string(dArray[0])) == "exit" {
		fmt.Println("Exiting game server!")
        }
	if strings.TrimSpace(string(dArray[0])) == "echo" {
		response = response + " { " + strings.TrimSpace(string(dArray[1])) 
        }
	if strings.TrimSpace(string(dArray[0])) == "health" {
		response = (strconv.Itoa(currentPlayer.health))
	}
	if strings.TrimSpace(string(dArray[0])) == "sethealth" {
		var err error
		currentPlayer.health, err = strconv.Atoi(strings.TrimSpace(string(dArray[1])))
		if err != nil {
			fmt.Println(err)
		}
	}
	if strings.TrimSpace(string(dArray[0])) == "animate" {
		displayAnimation(connId, "exampleAnimation", 5, 3, 50)
	}
	if strings.TrimSpace(string(dArray[0])) == "exit" {
		fmt.Println("Player " + connId + " has left the game.")
		response += `"exit": "successful"`
		delete(playerList, connId)
	}
	if strings.TrimSpace(string(dArray[0])) == "move" {
		blocked := "false"
		currentPlayer.x, currentPlayer.y = HandleMove(currentPlayer.x, currentPlayer.y, dArray)
		if currentPlayer.x == playerList[connId].x && currentPlayer.y == playerList[connId].y {
			blocked = "true"
		}
		playerList[connId] = currentPlayer
		fmt.Println("PlayerX:", currentPlayer.x, "PlayerY:", currentPlayer.y)
		response = response + `"x": ` + strconv.Itoa(int(currentPlayer.x)) + `, "y": ` + strconv.Itoa(int(currentPlayer.y)) + `, "blocked": ` + blocked
	}

	response = response + " } ] "

	return response
}

func getActions(playerId string) string {
return ""
}

func updateClient(playerId string) string {
	//var tempMessage string
	message := "["
	i := 0
	for key, currentBlock := range blockList {
		visible := false
		for x := currentBlock.x; x <= currentBlock.endX; x += increment {
			for y := currentBlock.y; y <= currentBlock.endY; y += increment  {
				/* This is the old code that is the same as the getObjectDistance function. It should be removed if proven unessecary.
				distanceX := getDifferenceFloat64(playerList[playerId].x, x)
				distanceY := getDifferenceFloat64(playerList[playerId].y, y)
				lineDistance := math.Sqrt(math.Pow(float64(distanceX), 2) + math.Pow(float64(distanceY), 2))
				*/
				if getObjectDistance(playerList[playerId].x, playerList[playerId].y, x, y) <= playerList[playerId].renderDistance {
					visible = true
					//Should add an exit here to optimise.
					/*
					if tempMessage != "" {
						tempMessage = tempMessage + `{"x": ` + strconv.Itoa(int(x)) + `, "y": ` + strconv.Itoa(int(y)) + `}, `
					}
					*/
					/*
					if tempMessage != "" {
						tempMessage = tempMessage + ", "
					}
					tempMessage = tempMessage + `{"x": ` + strconv.Itoa(int(x)) + `, "y": ` + strconv.Itoa(int(y)) + `}`
					*/
					// It seems like the only way to make this work is to either create a queue for blocks that are visible that need to be sent, !!!or to add the , to the previous part, and not add it if it is the first one! 
					//After reading this comment from long ago, I think I will simply make it so the program checks every coardinate in a block, and if it is visible, it checks if the block with the same uuid has already been shown. If so, it doesn't add anything.
					//I will also change how blocks work, making it so each block has a starting coardinate set, and an ending one, otherwise it will cause rendering issues. < Done
					//An alternative is to change the width and height as you send it to the client, depending on what is visible, but I think that introduces unnessacary complexity, and restricts the client unessircatly.
					//if y != currentBlock.y + currentBlock.height {
					//}
				}
			}
		}
		if visible == true {
			// Each block needs to be surrounded by curly brackets, within square, eg [{x 5 y 6}, {x 6 y 7}]
			message = message + ` { "id": "` + key + ` ", "blockType": "` + currentBlock.blockType  + `", "x": ` + strconv.Itoa(int(currentBlock.x)) + `, "y": ` + strconv.Itoa(int(currentBlock.y)) + `, "endX": ` + strconv.Itoa(int(currentBlock.endX)) + `, "endY": ` + strconv.Itoa(int(currentBlock.endY)) + ` ] }`
			//This line is from when all the visible blocks had their coardinates individual coardinates sent to the client, not just their details. Not sending them is likely a better way to do things, but I have left this here to 
			//message = message + ` { "id": "` + key + ` ", "blockType": "` + currentBlock.blockType  + `", "x": ` + strconv.Itoa(int(currentBlock.x)) + `, "y": ` + strconv.Itoa(int(currentBlock.y)) + `, "endX": ` + strconv.Itoa(int(currentBlock.endX)) + `, "endY": ` + strconv.Itoa(int(currentBlock.endY)) + `, ` + `"blocks": [ ` + tempMessage + ` ] }`
			if i != len(blockList) - 1 {
				message = message + ", "
			}
		}
		i = i + 1
	}
	message = message + "]"
	return message
}


func handleConnections(connId string) {

	var message string

	var response string
	var renderUpdates string

	playerList[connId] = player{health: 20, x: 0, y: 1, renderDistance: 3}
	
//Send opening information to the player.

        for {
		//Get what the player wants to do and then send a response.
		message = ""

		response = ""
		renderUpdates = ""

                data, err := bufio.NewReader(connList[connId]).ReadString('\n')
                if err != nil {
                        fmt.Println(err)
                        return
                }
		dArray := strings.Split(data, " ")
		response = handleActions(connId, dArray)

		actions := getActions(connId)

		renderUpdates = updateClient(connId)
		
		message = response + actions + renderUpdates + "\n"
                connList[connId].Write([]byte(message))
        }
}

func gameLoop() {
	for {
		for key, currentBlock := range blockList {
			if currentBlock.blockType == "flicker" {
				if currentBlock.sinceLastFlicker >= 5 {
				if currentBlock.flickerUp == true {
						currentBlock.y = currentBlock.y - 1
						currentBlock.endY -= 1
						currentBlock.flickerUp = false
					} else {
						currentBlock.y = currentBlock.y + 1
						currentBlock.endY += 1
						currentBlock.flickerUp = true
					}
				} else {
					currentBlock.sinceLastFlicker = currentBlock.sinceLastFlicker + 1
				}
			}
			blockList[key] = currentBlock
		}
		time.Sleep(time.Second)
	}
}

func main() {
	blockList[genUUID()] = block{ blockType: "basic", x: 4, y: 4, endX: 7, endY: 7}
//	blockList[genUUID()] = block{ blockType: "flicker", x: 6, y: 4, endX: 7, endY: 6}
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
