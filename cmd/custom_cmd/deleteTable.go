package custom_cmd

import (
	"fmt"

	"github.com/NazarStinyanskiy/ns-data/internal/catalog"
)

func deleteTable(args []string, db *catalog.DB) error {
	if len(args) == 0 {
		return fmt.Errorf("expected table name")
	}
	if len(args) > 1 {
		return fmt.Errorf("expected only table name. '%s' - not resolvable", args[1])
	}

	tableName := args[0]
	for i := range db.NsDataMetadata.Tables {
		if db.NsDataMetadata.Tables[i].Name == tableName {
			fmt.Printf("Deleting table '%s'\n", tableName)
			return db.DeleteTable(tableName)
		}
	}
	return fmt.Errorf("unknown table: %s", tableName)
}
