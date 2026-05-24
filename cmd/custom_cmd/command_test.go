package custom_cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/NazarStinyanskiy/ns-data/internal/catalog"
)

func TestMain(m *testing.M) {
	PrepareCommands()
	os.Exit(m.Run())
}

// emptyDB returns a DB with no tables and no disk backing (sufficient for error-path tests).
func emptyDB() *catalog.DB {
	return &catalog.DB{
		NsDataMetadata: catalog.NsDataMetadata{
			Version: "0.1.0",
			DbName:  "testdb",
		},
	}
}

// diskDB creates a temp directory, writes an empty metadata file into it, overrides
// catalog.NsDataPath for the duration of the test, and returns the DB + cleanup func.
func diskDB(t *testing.T) (*catalog.DB, func()) {
	t.Helper()
	dir := t.TempDir()

	origPath := catalog.NsDataPath
	catalog.NsDataPath = dir

	dbDir := filepath.Join(dir, "testdb")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		t.Fatal(err)
	}

	db := &catalog.DB{
		NsDataMetadata: catalog.NsDataMetadata{
			Version: "0.1.0",
			DbName:  "testdb",
		},
	}
	meta, _ := json.Marshal(db.NsDataMetadata)
	if err := os.WriteFile(filepath.Join(dbDir, catalog.NsDataMetadataFile), meta, 0644); err != nil {
		t.Fatal(err)
	}

	return db, func() { catalog.NsDataPath = origPath }
}

// --- Execute routing ---

func TestExecute_UnknownTopLevelCommand(t *testing.T) {
	err := Execute("foo bar", emptyDB())
	if err == nil {
		t.Fatal("expected error for unknown command")
	}
}

func TestExecute_UnknownSubCommand(t *testing.T) {
	err := Execute("create index foo", emptyDB())
	if err == nil {
		t.Fatal("expected error for unknown sub-command")
	}
}

func TestExecute_IncompleteCommandNoLeaf(t *testing.T) {
	// "create" alone resolves to a map, never a function → not found
	err := Execute("create", emptyDB())
	if err == nil {
		t.Fatal("expected error for incomplete command with no leaf function")
	}
}

func TestExecute_CaseInsensitive(t *testing.T) {
	db, cleanup := diskDB(t)
	defer cleanup()

	// Command words (CREATE, TABLE) are lowercased during routing; modifier must stay uppercase.
	err := Execute("CREATE TABLE employees ( id INTEGER PK )", db)
	if err != nil {
		t.Fatalf("uppercase command routing failed: %v", err)
	}
}

func TestExecute_CreateTable_ModifierCaseSensitive(t *testing.T) {
	// Modifiers are NOT normalized — lowercase "pk" is rejected while column types are uppercased.
	err := Execute("create table employees ( id INTEGER pk )", emptyDB())
	if err == nil {
		t.Fatal("expected error: modifier 'pk' should not match 'PK' (modifiers are case-sensitive)")
	}
}

// --- create table ---

func TestExecute_CreateTable_TooFewArgs(t *testing.T) {
	err := Execute("create table employees", emptyDB())
	if err == nil {
		t.Fatal("expected error for too few arguments")
	}
}

func TestExecute_CreateTable_MissingParens(t *testing.T) {
	err := Execute("create table employees id INTEGER pk", emptyDB())
	if err == nil {
		t.Fatal("expected error for missing parentheses")
	}
}

func TestExecute_CreateTable_InvalidColumnType(t *testing.T) {
	err := Execute("create table employees ( id BLOB pk )", emptyDB())
	if err == nil {
		t.Fatal("expected error for invalid column type")
	}
}

func TestExecute_CreateTable_MissingPrimaryKey(t *testing.T) {
	err := Execute("create table employees ( id INTEGER )", emptyDB())
	if err == nil {
		t.Fatal("expected error when primary key is absent")
	}
}

func TestExecute_CreateTable_DuplicatePrimaryKey(t *testing.T) {
	err := Execute("create table employees ( id INTEGER pk , name VARCHAR pk )", emptyDB())
	if err == nil {
		t.Fatal("expected error for duplicate primary key")
	}
}

func TestExecute_CreateTable_InvalidModifier(t *testing.T) {
	err := Execute("create table employees ( id INTEGER notnull )", emptyDB())
	if err == nil {
		t.Fatal("expected error for invalid modifier")
	}
}

func TestExecute_CreateTable_AlreadyExists(t *testing.T) {
	db := emptyDB()
	db.NsDataMetadata.Tables = []catalog.Table{{Name: "employees"}}

	err := Execute("create table employees ( id INTEGER pk )", db)
	if err == nil {
		t.Fatal("expected error when table already exists")
	}
}

func TestExecute_CreateTable_Success(t *testing.T) {
	db, cleanup := diskDB(t)
	defer cleanup()

	err := Execute("create table employees ( id INTEGER PK , name VARCHAR )", db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(db.NsDataMetadata.Tables) != 1 {
		t.Fatalf("expected 1 table, got %d", len(db.NsDataMetadata.Tables))
	}
	if db.NsDataMetadata.Tables[0].Name != "employees" {
		t.Fatalf("expected table name 'employees', got %q", db.NsDataMetadata.Tables[0].Name)
	}
}

func TestExecute_CreateTable_SuccessMultipleColumns(t *testing.T) {
	db, cleanup := diskDB(t)
	defer cleanup()

	err := Execute("create table users ( id INTEGER PK , email VARCHAR , age INTEGER )", db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cols := db.NsDataMetadata.Tables[0].Columns
	if len(cols) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(cols))
	}
}

// --- create table: spacing variants (normalizeArgs) ---

func TestExecute_CreateTable_SpacingVariants(t *testing.T) {
	tests := []struct {
		name string
		cmd  string
	}{
		{
			"paren attached to first column, comma attached to second column name",
			"create table t (id INTEGER PK ,name VARCHAR)",
		},
		{
			"paren attached to first column, spaces around comma",
			"create table t (id INTEGER PK , name VARCHAR)",
		},
		{
			"standard spacing with paren attached to close",
			"create table t ( id INTEGER PK , name VARCHAR)",
		},
		{
			"comma attached after modifier, space before next column",
			"create table t ( id INTEGER PK, name VARCHAR)",
		},
		{
			"paren attached to first column and comma attached after modifier",
			"create table t (id INTEGER PK, name VARCHAR)",
		},
		{
			"all tokens packed: paren+first, comma+second, close paren attached",
			"create table t (id INTEGER PK ,name VARCHAR)",
		},
		{
			"three columns with mixed attachment",
			"create table t (id INTEGER PK ,name VARCHAR ,age INTEGER)",
		},
		{
			"sandwich: comma between modifier and next column name, no spaces",
			"create table t ( id INTEGER PK,name VARCHAR)",
		},
		{
			"sandwich: comma between type and next column name, no spaces",
			"create table t ( id INTEGER,name VARCHAR PK)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, cleanup := diskDB(t)
			defer cleanup()

			if err := Execute(tc.cmd, db); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(db.NsDataMetadata.Tables) != 1 {
				t.Fatalf("expected 1 table, got %d", len(db.NsDataMetadata.Tables))
			}
		})
	}
}

// --- list tables ---

func TestExecute_ListTables_Empty(t *testing.T) {
	err := Execute("list tables", emptyDB())
	if err != nil {
		t.Fatalf("unexpected error listing empty tables: %v", err)
	}
}

func TestExecute_ListTables_WithTables(t *testing.T) {
	db := emptyDB()
	db.NsDataMetadata.Tables = []catalog.Table{{Name: "employees"}, {Name: "orders"}}

	err := Execute("list tables", db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecute_ListTables_TooManyArgs(t *testing.T) {
	err := Execute("list tables extra", emptyDB())
	if err == nil {
		t.Fatal("expected error for extra arguments")
	}
}

// --- describe table ---

func TestExecute_DescribeTable_NoTableName(t *testing.T) {
	err := Execute("describe table", emptyDB())
	if err == nil {
		t.Fatal("expected error when table name is missing")
	}
}

func TestExecute_DescribeTable_TooManyArgs(t *testing.T) {
	err := Execute("describe table employees extra", emptyDB())
	if err == nil {
		t.Fatal("expected error for extra arguments after table name")
	}
}

func TestExecute_DescribeTable_UnknownTable(t *testing.T) {
	err := Execute("describe table employees", emptyDB())
	if err == nil {
		t.Fatal("expected error for unknown table")
	}
}

func TestExecute_DescribeTable_Success(t *testing.T) {
	db := emptyDB()
	db.NsDataMetadata.Tables = []catalog.Table{
		{
			Name: "employees",
			Columns: []catalog.Column{
				{Name: "id", Type: catalog.INTEGER, Modifiers: []catalog.Modifier{catalog.PK}},
				{Name: "name", Type: catalog.VARCHAR},
			},
		},
	}

	err := Execute("describe table employees", db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
