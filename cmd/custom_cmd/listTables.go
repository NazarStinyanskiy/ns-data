package custom_cmd

import (
	"fmt"

	"github.com/NazarStinyanskiy/ns-data/internal/catalog"
)

func listTables(args []string, db *catalog.DB) error {
	if len(args) != 0 {
		return fmt.Errorf("too many arguments")
	}
	if len(db.NsDataMetadata.Tables) == 0 {
		fmt.Printf("No tables found\n")
		return nil
	}
	fmt.Println("TABLES:")
	for _, table := range db.NsDataMetadata.Tables {
		fmt.Println(table.Name)
	}
	return nil
}
