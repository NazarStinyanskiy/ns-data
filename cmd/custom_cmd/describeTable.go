package custom_cmd

import (
	"fmt"
	"strings"

	"github.com/NazarStinyanskiy/ns-data/internal/catalog"
)

func describeTable(args []string, db *catalog.DB) error {
	if len(args) == 0 {
		return fmt.Errorf("expected table name")
	}
	if len(args) > 1 {
		return fmt.Errorf("expected only table name. '%s' - not resolvable", args[1])
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

func printTable(table catalog.Table) {
	type row struct{ name, typ, modifiers string }

	rows := make([]row, len(table.Columns))
	var builder strings.Builder
	for i, col := range table.Columns {
		for mId, mod := range col.Modifiers {
			builder.WriteString(string(mod))
			if mId < len(col.Modifiers)-1 {
				builder.WriteString(", ")
			}
		}
		rows[i] = row{col.Name, string(col.Type), builder.String()}
		builder.Reset()
	}

	w0, w1, w2 := len("COLUMN"), len("TYPE"), len("MODIFIERS")
	for _, r := range rows {
		if len(r.name) > w0 {
			w0 = len(r.name)
		}
		if len(r.typ) > w1 {
			w1 = len(r.typ)
		}
		if len(r.modifiers) > w2 {
			w2 = len(r.modifiers)
		}
	}

	sep := fmt.Sprintf("+-%s-+-%s-+-%s-+", strings.Repeat("-", w0), strings.Repeat("-", w1), strings.Repeat("-", w2))
	rowFmt := fmt.Sprintf("| %%-%ds | %%-%ds | %%-%ds |", w0, w1, w2)

	fmt.Println(sep)
	fmt.Printf(rowFmt+"\n", "COLUMN", "TYPE", "MODIFIERS")
	fmt.Println(sep)
	for _, r := range rows {
		fmt.Printf(rowFmt+"\n", r.name, r.typ, r.modifiers)
	}
	fmt.Println(sep)
}
