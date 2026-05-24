package custom_cmd

import (
	"fmt"
	"strings"

	"github.com/NazarStinyanskiy/ns-data/internal/catalog"
)

var commands = map[string]interface{}{}

func PrepareCommands() {
	create := map[string]interface{}{
		"table": createTable,
	}
	del := map[string]interface{}{
		"table": deleteTable,
	}
	describe := map[string]interface{}{
		"table": describeTable,
	}
	list := map[string]interface{}{
		"tables": listTables,
	}
	commands["create"] = create
	commands["delete"] = del
	commands["describe"] = describe
	commands["list"] = list
}

func Execute(fullCommand string, db *catalog.DB) error {
	split := strings.Split(fullCommand, " ")
	args := split
	curMap := commands
	found := false
	for _, cmd := range split {
		args = append(args[1:])
		cmd = strings.ToLower(cmd)
		if curMap[cmd] == nil {
			return fmt.Errorf("unknown command: %s", cmd)
		}
		switch subMap := curMap[cmd].(type) {
		case map[string]interface{}:
			curMap = subMap
			continue
		case interface{}:
			err := subMap.(func(args []string, db *catalog.DB) error)(args, db)
			if err != nil {
				return err
			}
			found = true
		}
		if found {
			break
		}
	}
	if !found {
		return fmt.Errorf("unknown command: %s", fullCommand)
	}
	return nil
}
