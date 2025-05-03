package cmd

import (
	"database/sql"
	"log"
	"os"

	"github.com/Ayobami6/schema_dump/internal"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var (
	dbType    string
	dbURL     string
	tableName string
	lang      string
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

var listTableCommand = &cobra.Command{
	Use:   "list-tables",
	Short: "List tables in the database",
	Run: func(cmd *cobra.Command, args []string) {
		if dbType != "postgres" {
			log.Fatalf("Only postgres is supported for now")
		}

		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		defer db.Close()

		tables := internal.Tables(db)
		// write the table to json
		outFile, err := os.Create("tables.json")
		if err != nil {
			log.Fatalf("Failed to create file: %v", err)
		}
		defer outFile.Close()
		// write to the outfile
		_, err = outFile.WriteString("[\n")
		if err != nil {
			log.Fatalf("Failed to write to file: %v", err)
		}
		for i, table := range tables {
			_, err = outFile.WriteString("{\"table_name\": \"" + table + "\"}")
			if err != nil {
				log.Fatalf("Failed to write to file: %v", err)
			}
			if i != len(tables)-1 {
				_, err = outFile.WriteString(",\n")
				if err != nil {
					log.Fatalf("Failed to write to file: %v", err)
				}
			}
		}
		_, err = outFile.WriteString("]\n")
		if err != nil {
			log.Fatalf("Failed to write to file: %v", err)
		}
		log.Println("Tables dumped to tables.json")

	},
}

var transformCommand = &cobra.Command{
	Use:   "transform",
	Short: "Transform SQL schema to a Language Model",
	Run: func(cmd *cobra.Command, args []string) {
		if dbType != "postgres" {
			log.Fatalf("Only postgres is supported for now")
		}

		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		defer db.Close()
		supportedLangs := map[string]bool{
			"python":     true,
			"typescript": true,
			"java":       true,
			"rust":       true,
			"go":         true,
		}
		if _, ok := supportedLangs[lang]; !ok {
			log.Fatalf("Language %s is not supported", lang)
		}
		err = internal.TransformToORMModel(lang, tableName, db)
		if err != nil {
			log.Fatalf("Failed to transform schema: %v", err)
		}
		log.Printf("Schema transformed to %s ORM model", lang)

	},
}

func init() {
	RootCmd.AddCommand(dumpSchemaCmd)
	RootCmd.AddCommand(listTableCommand)
	RootCmd.AddCommand(transformCommand)

	listTableCommand.Flags().StringVar(&dbType, "db", "", "Database type (e.g., postgres)")
	listTableCommand.Flags().StringVar(&dbURL, "url", "", "Database connection URL")

	dumpSchemaCmd.Flags().StringVar(&dbType, "db", "", "Database type (e.g., postgres)")
	dumpSchemaCmd.Flags().StringVar(&dbURL, "url", "", "Database connection URL")
	dumpSchemaCmd.Flags().StringVar(&tableName, "table", "", "Table name to dump (optional)")
	transformCommand.Flags().StringVar(&dbType, "db", "", "Database type (e.g., postgres)")
	transformCommand.Flags().StringVar(&dbURL, "url", "", "Database connection URL")
	transformCommand.Flags().StringVar(&tableName, "table", "", "Table name to dump (optional)")
	transformCommand.Flags().StringVar(&lang, "lang", "", "Language to transform to (e.g., python, typescript, java, rust, go)")

	dumpSchemaCmd.MarkFlagRequired("db")
	dumpSchemaCmd.MarkFlagRequired("url")
	listTableCommand.MarkFlagRequired("db")
	listTableCommand.MarkFlagRequired("url")
	transformCommand.MarkFlagRequired("db")
	transformCommand.MarkFlagRequired("url")
	transformCommand.MarkFlagRequired("lang")
	transformCommand.MarkFlagRequired("table")
}
