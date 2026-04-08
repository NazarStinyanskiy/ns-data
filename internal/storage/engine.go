package storage

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

type DbMetadata struct {
	Version string `json:"version"`
}

func InitDb(path string) error {
	slog.Debug("Initializing database...")
	err := os.Mkdir(path, 0755)
	if err != nil {
		slog.Error("Error creating db catalog", "description", err)
		return err
	}
	data, err := json.Marshal(DbMetadata{viper.GetString("version")})
	err = os.WriteFile(path+"/nsdata_metadata.json", data, 0755)
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
