package custom_cmd

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/NazarStinyanskiy/ns-data/internal/catalog"
)

var commands = map[string]interface{}{}

func PrepareCommands() {
	create := map[string]interface{}{
		"table": createTable,
	}
	describe := map[string]interface{}{
		"table": describeTable,
	}
	list := map[string]interface{}{
		"tables": listTables,
	}
	commands["create"] = create
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

// create table employees ( id int , name varchar )
func createTable(args []string, db *catalog.DB) error {
	if len(args) < 5 {
		return fmt.Errorf("expected at least 5 arguments")
	}
	tableName := args[0]
	if args[1] != "(" || args[len(args)-1] != ")" {
		return fmt.Errorf("expected '(' or ')'")
	}
	isName := true
	isType := false
	isModifier := false
	var columns []catalog.Column
	columnId := 0
	for i := 2; i < len(args)-1; i++ {
		if isName {
			columns = append(columns, catalog.Column{Name: args[i]})
			isName = false
			isType = true
			continue
		}
		if isType {
			columnType, err := getColumnType(args[i])
			if err != nil {
				return err
			}
			columns[columnId].Type = columnType
			isType = false
			isModifier = true
			continue
		}
		if args[i] == "," {
			columnId++
			isName = true
			isType = false
			isModifier = false
			continue
		}
		if isModifier {
			columns[columnId].Modifiers = append(columns[columnId].Modifiers, catalog.Modifier(args[i]))
			continue
		}
	}
	err := db.CreateTable(tableName, columns)
	if err != nil {
		return err
	}
	slog.Info("Created table successfully")
	return nil
}

// describe table a
func describeTable(args []string, db *catalog.DB) error {
	if len(args) == 0 {
		return fmt.Errorf("expected table name")
	}
	if len(args) > 1 {
		return fmt.Errorf("expected only table name. '" + fmt.Sprint(args[1]) + "' - not resolvable")
	}

	tableName := args[0]
	found := false
	for i := range db.NsDataMetadata.Tables {
		if db.NsDataMetadata.Tables[i].Name == tableName {
			found = true
			printTable(db.NsDataMetadata.Tables[i])
			break
		}
	}
	if !found {
		return fmt.Errorf("unknown table: %s", tableName)
	}
	return nil
}

func listTables(args []string, db *catalog.DB) error {
	if len(args) != 0 {
		return fmt.Errorf("too many arguments")
	}
	fmt.Println("TABLES:")
	for _, table := range db.NsDataMetadata.Tables {
		fmt.Println(table.Name)
	}
	return nil
}

func getColumnType(t string) (catalog.ColumnType, error) {
	cT := catalog.ColumnType(strings.ToUpper(t))
	switch cT {
	case catalog.INTEGER:
		return catalog.INTEGER, nil
	case catalog.VARCHAR:
		return catalog.VARCHAR, nil
	default:
		return "", fmt.Errorf("Wrong type: " + t)
	}
}

// rewrite using %-20s
func printTable(table catalog.Table) {
	fmt.Println("|COLUMN\t|TYPE\t|MODIFIERS\t|")
	var builder strings.Builder
	for i := range table.Columns {
		column := table.Columns[i]
		for mId, modifier := range column.Modifiers {
			builder.WriteString(string(modifier))
			if mId < len(column.Modifiers)-1 {
				builder.WriteString(", ")
			}
		}
		fmt.Println("|" + column.Name + "\t|" + string(column.Type) + "\t|" + builder.String() + "\t|")
		builder.Reset()
	}
}
