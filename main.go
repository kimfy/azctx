package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := NewModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	accounts []account
	cursor   int
}

func NewModel() model {
	return model{
		accounts: make([]account, 0),
		cursor:   0,
	}
}

func (m model) Init() tea.Cmd {
  if len(os.Args) == 1 {
    log.Fatal("No input provided")
  }
	return m.parseStdin
}

func (m model) parseStdin() tea.Msg {
	var raw []byte
	r := bufio.NewReader(os.Stdin)
	raw, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	if !json.Valid(raw) {
		s := "Invalid JSON\n"
		s += fmt.Sprintf("Input: %s\n", raw)
		log.Fatal(s)
	}

	var accounts []account
	var users []user
	err = json.Unmarshal(raw, &accounts)
	if err != nil {
		log.Fatal(err)
	}

	for _, a := range accounts {
		users = append(users, a.User)
	}
	m.accounts = append(m.accounts, accounts...)
	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case model:
		m = msg
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "j", "down":
			if m.cursor < len(m.accounts)-1 {
				m.cursor++
			}

		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "enter":
			return m.enter()
		}
	}
	return m, nil
}

func (m model) View() string {
	var s string
	for i, a := range m.accounts {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s %s\n", cursor, a.Name, a.ID)
	}
	return s
}

func (m model) enter() (tea.Model, tea.Cmd) {
	a := m.accounts[m.cursor]
	var s string
	if _, err := exec.LookPath("az"); err != nil {
		s += "\"az\" executable file not found in $PATH\n"
		s += "Run the following command if you have az installed elsewhere:\n"
		s += fmt.Sprintf("\taz account set --subscription %s\n", a.ID)
		fmt.Print(s)
		return m, tea.Quit
	}

	if err := exec.Command("az", "account", "set", "--subscription", a.ID).Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Set subscription to %s (%s)\n", a.Name, a.ID)
	return m, tea.Quit
}

type account struct {
	EnvironmentName string `json:"environmentName"`
	HomeTenantID    string `json:"homeTenantId"`
	ID              string `json:"id"`
	IsDefault       bool   `json:"isDefault"`
	Name            string `json:"name"`
	State           string `json:"state"`
	TenantID        string `json:"tenantId"`
	User            user   `json:"user"`
}

type user struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
