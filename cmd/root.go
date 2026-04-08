package cmd

import (
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/NazarStinyanskiy/ns-data/cmd/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "nsdata",
	Short: "CLI util for managing nsdata DBMS",
	Long:  "CLI util for managing nsdata DBMS",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}
		level := slog.LevelInfo
		if verbose {
			level = slog.LevelDebug
		}
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
		slog.SetDefault(logger)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func init() {
	loadConfig()
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Version = viper.GetString("version")

	rootCmd.PersistentFlags().BoolP("verbose", "V", false, "verbose output")
	rootCmd.AddCommand(db.Db)
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./cfg/")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}
