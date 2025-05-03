package internal

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type Column struct {
	ColumnName string
	DataType   string
	IsNullable string
}

type ForeignKey struct {
	SourceTable    string
	SourceColumn   string
	TargetTable    string
	TargetColumn   string
	ConstraintName string
}

func PostgresSchemaDump(db *sql.DB, tableName string) {
	if tableName != "" {
		dumpTableSchema(db, tableName)
	} else {
		dumpAllSchemas(db)
	}
}

func dumpTableSchema(db *sql.DB, tableName string) {
	rows, err := db.Query(`
        SELECT table_name, column_name, data_type, is_nullable
        FROM information_schema.columns
        WHERE table_schema = 'public'
        ORDER BY table_name, ordinal_position;
    `)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	schema := make(map[string][]Column)

	for rows.Next() {
		var table, colName, dataType, nullable string
		err := rows.Scan(&table, &colName, &dataType, &nullable)
		if err != nil {
			log.Fatal(err)
		}

		if table == tableName {
			schema[table] = append(schema[table], Column{
				ColumnName: colName,
				DataType:   dataType,
				IsNullable: nullable,
			})
		}

	}

	outFile, err := os.Create("schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	pkRows, err := db.Query(`
        SELECT tc.table_name, kcu.column_name
			FROM information_schema.table_constraints AS tc
			JOIN information_schema.key_column_usage AS kcu
			ON tc.constraint_name = kcu.constraint_name
			AND tc.table_schema = kcu.table_schema
			WHERE tc.constraint_type = 'PRIMARY KEY'
			AND tc.table_schema = 'public';

    `)
	if err != nil {
		log.Fatal(err)
	}
	defer pkRows.Close()

	fkRows, err := db.Query(`
    SELECT
        tc.table_name AS source_table,
        kcu.column_name AS source_column,
        ccu.table_name AS target_table,
        ccu.column_name AS target_column,
        tc.constraint_name
    FROM information_schema.table_constraints AS tc
    JOIN information_schema.key_column_usage AS kcu
        ON tc.constraint_name = kcu.constraint_name
        AND tc.table_schema = kcu.table_schema
    JOIN information_schema.constraint_column_usage AS ccu
        ON ccu.constraint_name = tc.constraint_name
        AND ccu.table_schema = tc.table_schema
    WHERE tc.constraint_type = 'FOREIGN KEY'
      AND tc.table_schema = 'public';
`)
	if err != nil {
		log.Fatal(err)
	}
	defer fkRows.Close()

	foreignKeys := make(map[string][]ForeignKey)
	for fkRows.Next() {
		var fk ForeignKey
		if err := fkRows.Scan(&fk.SourceTable, &fk.SourceColumn, &fk.TargetTable, &fk.TargetColumn, &fk.ConstraintName); err != nil {
			log.Fatal(err)
		}
		fmt.Println("lets see the source table: ", fk.SourceTable)

		foreignKeys[fk.SourceTable] = append(foreignKeys[fk.SourceTable], fk)
	}

	primaryKeys := make(map[string][]string)
	for pkRows.Next() {
		var table, col string
		err := pkRows.Scan(&table, &col)
		if err != nil {
			log.Fatal(err)
		}
		primaryKeys[table] = append(primaryKeys[table], col)
	}

	for table, columns := range schema {
		fmt.Fprintf(outFile, "CREATE TABLE %s (\n", table)
		for _, col := range columns {
			null := "NOT NULL"
			if col.IsNullable == "YES" {
				null = "NULL"
			}
			comma := ","
			// if i == len(columns)-1 {
			// 	comma = ""
			// }
			fmt.Fprintf(outFile, "    %s %s %s%s\n", col.ColumnName, col.DataType, null, comma)
		}
		if fks, ok := foreignKeys[table]; ok {
			for _, fk := range fks {
				fmt.Fprintf(outFile, "    CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s),\n",
					fk.ConstraintName, fk.SourceColumn, fk.TargetTable, fk.TargetColumn)
			}
		}
		// Add primary key constraint
		if pkCols, ok := primaryKeys[table]; ok {
			fmt.Fprintf(outFile, "    PRIMARY KEY (%s)\n", join(pkCols, ", "))
		}

		fmt.Fprintln(outFile, ");\n")
	}

	fmt.Println("Schema written to schema.sql")
}

func dumpAllSchemas(db *sql.DB) {
	rows, err := db.Query(`
        SELECT table_name, column_name, data_type, is_nullable
        FROM information_schema.columns
        WHERE table_schema = 'public'
        ORDER BY table_name, ordinal_position;
    `)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	schema := make(map[string][]Column)

	for rows.Next() {
		var table, colName, dataType, nullable string
		err := rows.Scan(&table, &colName, &dataType, &nullable)
		if err != nil {
			log.Fatal(err)
		}

		schema[table] = append(schema[table], Column{
			ColumnName: colName,
			DataType:   dataType,
			IsNullable: nullable,
		})
	}

	outFile, err := os.Create("schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	pkRows, err := db.Query(`
        SELECT tc.table_name, kcu.column_name
			FROM information_schema.table_constraints AS tc
			JOIN information_schema.key_column_usage AS kcu
			ON tc.constraint_name = kcu.constraint_name
			AND tc.table_schema = kcu.table_schema
			WHERE tc.constraint_type = 'PRIMARY KEY'
			AND tc.table_schema = 'public';

    `)
	if err != nil {
		log.Fatal(err)
	}
	defer pkRows.Close()

	fkRows, err := db.Query(`
    SELECT
        tc.table_name AS source_table,
        kcu.column_name AS source_column,
        ccu.table_name AS target_table,
        ccu.column_name AS target_column,
        tc.constraint_name
    FROM information_schema.table_constraints AS tc
    JOIN information_schema.key_column_usage AS kcu
        ON tc.constraint_name = kcu.constraint_name
        AND tc.table_schema = kcu.table_schema
    JOIN information_schema.constraint_column_usage AS ccu
        ON ccu.constraint_name = tc.constraint_name
        AND ccu.table_schema = tc.table_schema
    WHERE tc.constraint_type = 'FOREIGN KEY'
      AND tc.table_schema = 'public';
`)
	if err != nil {
		log.Fatal(err)
	}
	defer fkRows.Close()

	foreignKeys := make(map[string][]ForeignKey)
	for fkRows.Next() {
		var fk ForeignKey
		if err := fkRows.Scan(&fk.SourceTable, &fk.SourceColumn, &fk.TargetTable, &fk.TargetColumn, &fk.ConstraintName); err != nil {
			log.Fatal(err)
		}
		foreignKeys[fk.SourceTable] = append(foreignKeys[fk.SourceTable], fk)
	}

	primaryKeys := make(map[string][]string)
	for pkRows.Next() {
		var table, col string
		err := pkRows.Scan(&table, &col)
		if err != nil {
			log.Fatal(err)
		}
		primaryKeys[table] = append(primaryKeys[table], col)
	}

	for table, columns := range schema {
		fmt.Fprintf(outFile, "CREATE TABLE %s (\n", table)
		for _, col := range columns {
			null := "NOT NULL"
			if col.IsNullable == "YES" {
				null = "NULL"
			}
			comma := ","
			// if i == len(columns)-1 {
			// 	comma = ""
			// }
			fmt.Fprintf(outFile, "    %s %s %s%s\n", col.ColumnName, col.DataType, null, comma)
		}
		if fks, ok := foreignKeys[table]; ok {
			for _, fk := range fks {
				fmt.Fprintf(outFile, "    CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s),\n",
					fk.ConstraintName, fk.SourceColumn, fk.TargetTable, fk.TargetColumn)
			}
		}
		// Add primary key constraint
		if pkCols, ok := primaryKeys[table]; ok {
			fmt.Fprintf(outFile, "    PRIMARY KEY (%s)\n", join(pkCols, ", "))
		}

		fmt.Fprintln(outFile, ");\n")
	}

	fmt.Println("Schema written to schema.sql")
}

func join(slice []string, sep string) string {
	out := ""
	for i, s := range slice {
		out += s
		if i < len(slice)-1 {
			out += sep
		}
	}
	return out
}

func Tables(sb *sql.DB) []string {
	rows, err := sb.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			log.Fatal(err)
		}
		tables = append(tables, table)
	}

	return tables
}
