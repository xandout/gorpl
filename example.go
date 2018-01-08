package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/xandout/repl/repl"
)

var mode = "table"

func main() {

	f := repl.New("> ", ";")

	f.AddAction("exit", func(args ...interface{}) (interface{}, error) {
		fmt.Println("Bye!")
		os.Exit(0)
		return nil, nil
	})
	f.AddAction("describe", func(args ...interface{}) (interface{}, error) {
		fmtString := "SELECT COLUMN_NAME,DATA_TYPE_NAME,LENGTH,IS_NULLABLE, SCHEMA_NAME FROM TABLE_COLUMNS WHERE TABLE_NAME = '%s';"
		if len(args) != 1 {
			return nil, errors.New("describe function requires a table name to be supplied")
		}
		fmt.Println(fmt.Sprintf(fmtString, args[0]))
		return "", nil
	})
	f.AddAction("mode", func(args ...interface{}) (interface{}, error) {
		modes := map[string]bool{
			"csv":   true,
			"table": true,
		}
		if len(args) != 1 {
			return nil, errors.New("mode function requires a valid output mode")
		}
		attemptedMode := args[0]
		if modes[attemptedMode.(string)] {
			mode = attemptedMode.(string)
		} else {
			fmt.Println("Valid modes are:")
			for k := range modes {
				fmt.Printf("\t%s\n", k)
			}
		}
		return "", nil

	})
	f.AddAction("show-mode", func(args ...interface{}) (interface{}, error) {
		fmt.Println(mode)
		return "", nil
	})
	f.Default = repl.Action{
		Action: func(args ...interface{}) (interface{}, error) {
			fmt.Println(args)
			return "", nil
		},
	}
	f.Start()

}
