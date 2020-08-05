package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type ActionInterface interface {
	Info() string
}

type Action struct {
	Type string `json:"type"`
}

type ActionHTTP struct {
	Action
	URL    string `json:"url"`
	Method string `json:"method"`
}

func (a ActionHTTP) Info() string {
	return fmt.Sprintf("%s - %s - %s", a.Type, a.URL, a.Method)
}

type ActionMail struct {
	Action
	Receiver string `json:"receiver"`
	Subject  string `json:"subject"`
}

func (a ActionMail) Info() string {
	return fmt.Sprintf("%s - %s - %s", a.Type, a.Receiver, a.Subject)
}

type Model struct {
	Actions    map[string]ActionInterface `json:"ractions"`
	ActionsRaw map[string]json.RawMessage `json:"actions"`
}

func (m *Model) Process() error {
	if m.ActionsRaw != nil {
		m.Actions = make(map[string]ActionInterface, len(m.ActionsRaw))
		for k, v := range m.ActionsRaw {
			a := Action{}
			if err := json.Unmarshal(v, &a); err != nil {
				return err
			}
			switch strings.ToLower(a.Type) {
			case "http":
				ah := ActionHTTP{}
				if err := json.Unmarshal(v, &ah); err != nil {
					return err
				}
				m.Actions[k] = ah
			case "mail":
				am := ActionMail{}
				if err := json.Unmarshal(v, &am); err != nil {
					return err
				}
				m.Actions[k] = am
			default:
				return errors.New("Unknown type")
			}
		}
	}
	return nil
}

func main() {
	str := `{"actions": {"action1": {"type": "http", "url": "http://localhost", "method": "GET"}, "action2": {"type": "mail", "receiver": "tui", "subject": "mail subject"}}}`
	m := Model{}
	if err := json.Unmarshal([]byte(str), &m); err != nil {
		panic(err)
	}
	if err := m.Process(); err != nil {
		panic(err)
	}
	fmt.Printf("%+v", m)
}
