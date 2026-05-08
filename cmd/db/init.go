package db

import (
	"log/slog"
	"os"

	"github.com/NazarStinyanskiy/ns-data/internal/catalog"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new db",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			slog.Error("usage: nsdata db init <name>")
			os.Exit(1)
		}
		err := catalog.InitDb(args[0])
		if err != nil {
			os.Exit(1)
		}
	},
}
