package db

import "github.com/spf13/cobra"

var Db = &cobra.Command{
	Use:   "db",
	Short: "Manage databases",
}

func init() {
	Db.AddCommand(initCmd)
}
