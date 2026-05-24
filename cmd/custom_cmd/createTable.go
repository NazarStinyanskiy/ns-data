package custom_cmd

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/NazarStinyanskiy/ns-data/internal/catalog"
)

// create table employees ( id int , name varchar )
func createTable(args []string, db *catalog.DB) error {
	args = normalizeArgs(args)
	if len(args) < 5 {
		return fmt.Errorf("expected at least 5 arguments")
	}
	tableName := args[0]
	if isTableExists(tableName, db) {
		return fmt.Errorf("table already exists: %s", tableName)
	}
	if args[1] != "(" || args[len(args)-1] != ")" {
		return fmt.Errorf("expected '(' or ')'")
	}
	isName := true
	isType := false
	isModifier := false
	pkExists := false
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
			modifier := catalog.Modifier(args[i])
			if !slices.Contains(catalog.ValidModifiers, modifier) {
				return fmt.Errorf("invalid modifier: %s", modifier)
			}
			if modifier == catalog.PK {
				if pkExists {
					return fmt.Errorf("duplicate primary key")
				}
				pkExists = true
			}
			columns[columnId].Modifiers = append(columns[columnId].Modifiers, modifier)
			continue
		}
	}
	if !pkExists {
		return fmt.Errorf("missing primary key: %s", tableName)
	}
	err := db.CreateTable(tableName, columns)
	if err != nil {
		return err
	}
	slog.Info("Created table successfully")
	return nil
}

func normalizeArgs(args []string) []string {
	newArgs := make([]string, 0)
	for _, arg := range args {
		if len(arg) > 1 && (strings.HasPrefix(arg, "(") || strings.HasPrefix(arg, ",")) {
			newArgs = append(newArgs, arg[0:1])
			newArgs = append(newArgs, arg[1:])
			continue
		}
		if len(arg) > 1 && (strings.HasSuffix(arg, ")") || strings.HasSuffix(arg, ",")) {
			newArgs = append(newArgs, arg[:len(arg)-1])
			newArgs = append(newArgs, arg[len(arg)-1:])
			continue
		}
		//PK,name - sandwich case
		if len(arg) > 2 && strings.Contains(arg, ",") && strings.LastIndex(arg, ",") != len(arg)-1 {
			commaIndex := strings.Index(arg, ",")
			newArgs = append(newArgs, arg[:commaIndex])
			newArgs = append(newArgs, arg[commaIndex:commaIndex+1])
			newArgs = append(newArgs, arg[commaIndex+1:])
			continue
		}
		newArgs = append(newArgs, arg)
	}
	return newArgs
}

func isTableExists(tableName string, db *catalog.DB) bool {
	for _, table := range db.NsDataMetadata.Tables {
		if table.Name == tableName {
			return true
		}
	}
	return false
}

func getColumnType(t string) (catalog.ColumnType, error) {
	cT := catalog.ColumnType(strings.ToUpper(t))
	if !slices.Contains(catalog.ValidColumnTypes, cT) {
		return cT, fmt.Errorf("invalid column type: %s", t)
	}
	return cT, nil
}
