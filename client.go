package main

import (
        "fmt"
        "os"
	"strconv"
        tea "github.com/charmbracelet/bubbletea"
)

type player struct {
	knownLocations []string
	base           int
	newsFeed       []article
	money          int
}

type article struct {
	title   string
	content string
	date    string
}

type model struct {
        buffer string
}

var user player

func initialModel() model {
        return model{}
}

func (m model) Init() tea.Cmd {
        return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        switch msg := msg.(type) {
        case tea.KeyMsg:
                switch msg.String() {
                case "ctrl+c", "q":
                        return m, tea.Quit
                case "esc":
                        break
                case "backspace":
                        if len(m.buffer) > 0 {
                                m.buffer = m.buffer[:len(m.buffer)-1]
                        }
                        return m, nil
                default:
                        m.buffer += msg.String()
                        return m, nil
                }
        }
        return m, nil
}

func (m model) View() string {
        s := ""
	s += "Coins: " + strconv.Itoa(user.money) + "\n"
        s += "Input: " + m.buffer
        return s
}

func main() {
	fmt.Println(`
 _____  _  _  _  _____  _        
|  |  ||_|| || ||   __||_| _____ 
|  |  || || || ||__   || ||     |
 \___/ |_||_||_||_____||_||_|_|_|

VERSION: somewhere-in-alpha
`)
        p := tea.NewProgram(initialModel())
        if _, err := p.Run(); err != nil {
                fmt.Printf("There has been an error: %v", err)
                os.Exit(1)
        }
}
