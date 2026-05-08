package db

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/NazarStinyanskiy/ns-data/cmd/custom_cmd"
	"github.com/NazarStinyanskiy/ns-data/internal/catalog"
	"github.com/spf13/cobra"
)

var connect = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a db to perform operations",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			slog.Error("usage: nsdata db connect <name>")
			os.Exit(1)
		}
		db, err := catalog.Connect(args[0])
		if err != nil {
			log.Fatal(err)
			return err
		}

		fmt.Printf("Connected to %s\n", args[0])
		fmt.Printf("Version: %s\n", db.NsDataMetadata.Version)

		reader := bufio.NewReader(os.Stdin)
		custom_cmd.PrepareCommands()
		for {
			fmt.Print("nsdata > ")

			line, err := reader.ReadString('\n')
			if err != nil {
				return err
			}

			line = strings.TrimSpace(line)

			switch line {
			case "exit", "quit":
				return nil
			}

			err = custom_cmd.Execute(line, db)
			if err != nil {
				slog.Error("Failed to execute command: "+line, "Description", err)
			}
		}
	},
}
