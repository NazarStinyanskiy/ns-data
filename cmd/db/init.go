package db

import (
	"log/slog"
	"os"

	"github.com/NazarStinyanskiy/ns-data/internal/storage"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new db",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			slog.Error("usage: nsdata init <path>")
			os.Exit(1)
		}
		err := storage.InitDb(args[0])
		if err != nil {
			os.Exit(1)
		}
	},
}
