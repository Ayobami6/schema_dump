package cmd

import (
	"database/sql"
	"log"

	"github.com/Ayobami6/schema_dump/internal"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var (
	dbType    string
	dbURL     string
	tableName string
)

var RootCmd = &cobra.Command{
	Use:   "Models Dump",
	Short: "Django Models Table Schema Dump",
}

// dumpSchemaCmd represents the dump-schema command
var dumpSchemaCmd = &cobra.Command{
	Use:   "dump-schema",
	Short: "Dump SQL schema from a live database",
	Run: func(cmd *cobra.Command, args []string) {
		if dbType != "postgres" {
			log.Fatalf("Only postgres is supported for now")
		}

		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		defer db.Close()

		internal.PostgresSchemaDump(db, tableName)
	},
}

func init() {
	RootCmd.AddCommand(dumpSchemaCmd)

	dumpSchemaCmd.Flags().StringVar(&dbType, "db", "", "Database type (e.g., postgres)")
	dumpSchemaCmd.Flags().StringVar(&dbURL, "url", "", "Database connection URL")
	dumpSchemaCmd.Flags().StringVar(&tableName, "table", "", "Table name to dump (optional)")

	dumpSchemaCmd.MarkFlagRequired("db")
	dumpSchemaCmd.MarkFlagRequired("url")
}
