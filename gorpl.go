package gorpl

import (
	"errors"
	"log"
	"strings"

	"github.com/chzyer/readline"

	"fmt"
)

// Action is a named func
type Action struct {
	Name   string
	Action func(args ...interface{}) (interface{}, error)
}

// Repl houses all of our config data
type Repl struct {
	RL         *readline.Instance
	Actions    map[string]Action
	Default    Action
	Prefix     string
	Terminator string
}

// New sets up the Repl
func New(prefix string, term string) Repl {
	r, err := readline.NewEx(&readline.Config{
		Prompt:                 "> ",
		HistoryFile:            "/tmp/repl.history",
		DisableAutoSaveHistory: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	return Repl{
		RL:         r,
		Prefix:     "/",
		Terminator: ";",
		Actions:    make(map[string]Action),
	}
}

// AddAction registers a named function, Action
func (repl *Repl) AddAction(cmd string, action func(args ...interface{}) (interface{}, error)) error {
	if strings.Count(cmd, " ") > 0 {
		return errors.New("cmd string can not contain white spaces")
	}

	repl.Actions[cmd] = Action{Name: cmd, Action: action}
	return nil

}

// Start runs the Repl as configured
func (repl *Repl) Start() error {

	defer repl.RL.Close()
	// Setup readline completer
	var pci []readline.PrefixCompleterInterface
	for _, action := range repl.Actions {
		pci = append(pci, readline.PcItem(action.Name))
	}

	repl.RL.Config.AutoComplete = readline.NewPrefixCompleter(pci...)
	// History
	var cmdHist []string
	var lines []string

	// If user does not supply default Action then we just echo
	if repl.Default.Action == nil {
		repl.Default = Action{
			Action: func(args ...interface{}) (interface{}, error) {
				fmt.Println(args)
				return nil, nil
			},
		}
	}

	for {
		line, err := repl.RL.Readline()
		if err != nil {
			log.Fatal(err)
		}
		// Remove extra whitespace
		line = strings.TrimSpace(line)

		cmd := strings.Split(line, " ")

		// Empty input
		if len(line) == 0 {
			continue
		}

		cmdHist = append(cmdHist, line)

		// Is this a built in command?
		if option, ok := repl.Actions[cmd[0]]; ok {
			args := make([]interface{}, len(cmd[1:]))
			for i, v := range cmd[1:] {
				args[i] = v
			}
			repl.RL.SetPrompt("> ")
			repl.RL.SaveHistory(line)
			ret, err := option.Action(args...)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Println(ret)
			}
		} else {
			// is not a registered Action, treat as string with multiline and pass to default.
			if !strings.HasSuffix(line, repl.Terminator) {
				lines = append(lines, line)
				repl.RL.SetPrompt(">>> ")
				repl.RL.SaveHistory(line)
				continue
			} else {
				args := make([]interface{}, 1)
				lines = append(lines, line)
				repl.RL.SetPrompt("> ")
				repl.RL.SaveHistory(line)
				args[0] = strings.Join(lines, " ")
				ret, err := repl.Default.Action(args...)
				if err != nil {
					log.Println(err)
				} else {
					fmt.Println(ret)
				}
				lines = []string{}
				continue
			}

		}

	}
}
