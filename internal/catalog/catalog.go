package catalog

import (
	"encoding/json"
	"log/slog"
	"os"
	"slices"

	"github.com/spf13/viper"
)

var NsDataPath = "./.nsdata"
var NsDataMetadataFile = "nsdata_metadata.json"

type DB struct {
	NsDataMetadata NsDataMetadata
}

type NsDataMetadata struct {
	Version string  `json:"version"`
	DbName  string  `json:"dbName"`
	Tables  []Table `json:"tables"`
}

func InitDb(name string) error {
	slog.Debug("Initializing database...")
	path := NsDataPath + "/" + name

	err := os.MkdirAll(path, 0755)
	if err != nil {
		slog.Error("Error creating db catalog", "description", err)
		return err
	}

	data, err := json.Marshal(NsDataMetadata{Version: viper.GetString("version"), DbName: name})
	if err != nil {
		slog.Error("Error creating db catalog", "description", err)
		return err
	}

	err = os.WriteFile(path+"/"+NsDataMetadataFile, data, 0755)
	if err != nil {
		slog.Error("Error creating nsdata_metadata.json file", "description", err)
		err := os.Remove(path)
		if err != nil {
			slog.Error("Error removing catalog", "description", err)
		}
		return err
	}

	slog.Info("Successfully initialized database", "path", path)
	return nil
}

func Connect(name string) (*DB, error) {
	bytes, err := os.ReadFile(NsDataPath + "/" + name + "/" + NsDataMetadataFile)
	if err != nil {
		slog.Error("Error reading catalog", "description", err)
		return &DB{}, err
	}

	var db DB
	err = json.Unmarshal(bytes, &db.NsDataMetadata)
	if err != nil {
		slog.Error("Error parsing catalog", "description", err)
		return &DB{}, err
	}

	return &db, nil
}

func (db *DB) CreateTable(tableName string, columns []Column) error {
	db.NsDataMetadata.Tables = append(db.NsDataMetadata.Tables, Table{Name: tableName, Columns: columns})
	return db.writeMetadata("Error creating table")
}

func (db *DB) DeleteTable(tableName string) error {
	for id := 0; id < len(db.NsDataMetadata.Tables); id++ {
		if db.NsDataMetadata.Tables[id].Name == tableName {
			db.NsDataMetadata.Tables = slices.Delete(db.NsDataMetadata.Tables, id, id+1)
			break
		}
	}
	return db.writeMetadata("Error deleting table")
}

func (db *DB) writeMetadata(errorMsg string) error {
	marshal, err := json.Marshal(db.NsDataMetadata)
	if err != nil {
		slog.Error(errorMsg, "description", err)
		return err
	}

	err = os.WriteFile(db.metadataPath(), marshal, 0644)
	if err != nil {
		slog.Error(errorMsg, "description", err)
		//Rollback NsDataMetadata here!! TODO
		return err
	}
	return nil
}

func (db *DB) metadataPath() string {
	return NsDataPath + "/" + db.NsDataMetadata.DbName + "/" + NsDataMetadataFile
}
