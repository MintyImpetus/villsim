package main

import (
        "fmt"
        "os"

        tea "github.com/charmbracelet/bubbletea"
)

type model struct {
        buffer string
}

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
        s := "Input: " + m.buffer
        return s
}

func main() {
        p := tea.NewProgram(initialModel())
        if _, err := p.Run(); err != nil {
                fmt.Printf("There has been an error: %v", err)
                os.Exit(1)
        }
}
