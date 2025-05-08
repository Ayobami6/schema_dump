package internal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/Ayobami6/schema_dump/utils"
	"github.com/zalando/go-keyring"
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

// transformToORMModel takes a language and transforms the SQL schema to the ORM model
// It uses the AzureAIClient to send a request to the Azure OpenAI API
func TransformToORMModel(lang, tableName string, db *sql.DB) error {
	hasDumped := make(chan bool, 1)
	// Generate the schema
	go func() {
		dumpTableSchema(db, tableName)
		hasDumped <- true

	}()

	osName := runtime.GOOS
	serviceName := "token"
	var token string
	var apiKey string
	// get the apikey
	apiKey, err := keyring.Get("apiKey", osName)
	if err != nil {
		// get token from keychain
		token, err = keyring.Get(serviceName, osName)
		if err != nil {
			// log.Printf("Error getting token from keychain: %v\n", err)
			if err == keyring.ErrNotFound {
				// if token is not found, fetch it from the token service
				token, err = FetchToken()
				if err != nil {
					// log.Printf("Error fetching token: %v\n", err)
					return fmt.Errorf("failed to fetch token: %w", err)
				}
				// save the token keychain
				err = keyring.Set(serviceName, osName, token)
				if err != nil {
					// log.Printf("Error saving token to keychain: %v\n", err)
					return fmt.Errorf("failed to save token to keychain: %w", err)
				}
				// get the api key
				apiKey, err = FetchAPIKey(token)
				if err != nil {
					// log.Printf("Error fetching api key: %v\n", err)
					return fmt.Errorf("failed to fetch api key: %w", err)
				}
				err = keyring.Set("apiKey", osName, apiKey)
				if err != nil {
					// log.Printf("Error saving api key to keychain: %v\n", err)
					return fmt.Errorf("failed to save api key to keychain: %w", err)
				}
				// set add the api key to the keyring
			} else {
				// log.Printf("Error getting token from keychain: %v\n", err)
				return fmt.Errorf("failed to get token from keychain: %w", err)
			}

		}
		apiKey, err = FetchAPIKey(token)
		if err != nil {
			// log.Printf("Error fetching api key: %v\n", err)
			return fmt.Errorf("failed to fetch api key: %w", err)
		}
		err = keyring.Set("apiKey", osName, apiKey)
		if err != nil {
			// log.Printf("Error saving api key to keychain: %v\n", err)
			return fmt.Errorf("failed to save api key to keychain: %w", err)
		}
	}
	// makes sure the dumpschema goroutine completes before sending prompt
	<-hasDumped
	// Read the schema from the file
	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}
	// Create the prompt
	prompt := fmt.Sprintf("Transform the following SQL schema to %s ORM model:\n%s", lang, string(schema))

	// Create the AzureAIClient
	client := utils.NewAzureAIClient("https://models.github.ai/inference/chat/completions", apiKey)
	// Send the request
	response, err := client.CreateCompletions(prompt)
	if err != nil {
		log.Printf("Error creating completions: %v", err)
		return fmt.Errorf("failed to create completions: %w", err)
	}
	// Write the response to a file
	ormModelFile, err := os.Create("orm_model.md")
	if err != nil {
		log.Printf("Error creating ORM model file: %v", err)
		return fmt.Errorf("failed to create ORM model file: %w", err)
	}
	_, err = ormModelFile.WriteString(response.(map[string]interface{})["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string))
	if err != nil {
		log.Printf("Error writing to ORM model file: %v", err)
		return fmt.Errorf("failed to write to ORM model file: %w", err)
	}
	return nil

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

type TokenResponse struct {
	Data       Data   `json:"data"`
	Message    string `json:"message"`
	Status     string `json:"status"`
	StatusCode int64  `json:"status_code"`
}

type Data struct {
	CreatedAt string `json:"created_at"`
	ID        int64  `json:"id"`
	IPAddress string `json:"ip_address"`
	Token     string `json:"token"`
}

type APIKeyResponse struct {
	Data       APIKeyData `json:"data"`
	Message    string     `json:"message"`
	Status     string     `json:"status"`
	StatusCode int64      `json:"status_code"`
}

type APIKeyData struct {
	AIApikey  string `json:"ai_apikey"`
	CreatedAt string `json:"created_at"`
	ID        int64  `json:"id"`
	Name      string `json:"name"`
}

func FetchAPIKey(token string) (string, error) {
	headers := map[string]interface{}{
		"Content-Type": "application/json",
		"Accept":       "application/json",
		"Token":        token,
	}
	// Create the request
	req := utils.NewRequest("GET", "https://utils.logizon.com/utils/first", nil, headers)
	resp, err := req.Send()
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	// Read the response
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get api key: %s", resp.Status)
	}
	// Decode the response
	var apiKeyResponse APIKeyResponse
	decoderErr := json.NewDecoder(resp.Body).Decode(&apiKeyResponse)
	if decoderErr != nil {
		return "", fmt.Errorf("failed to decode response: %w", decoderErr)
	}
	// Return the api key
	return apiKeyResponse.Data.AIApikey, nil

}

func FetchToken() (string, error) {
	tokenChan := make(chan string, 1)
	tokenError := make(chan error, 1)

	go func() {
		headers := map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		}
		req := utils.NewRequest("GET", "https://utils.logizon.com/tokens", nil, headers)

		resp, err := req.Send()
		if err != nil {
			tokenError <- fmt.Errorf("failed to send request: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			tokenError <- fmt.Errorf("failed to get token: %s", resp.Status)
			return
		}

		var tokenResponse TokenResponse
		if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
			tokenError <- fmt.Errorf("failed to decode response: %w", err)
			return
		}

		tokenChan <- tokenResponse.Data.Token
	}()

	// Wait for result or error (or timeout if needed)
	select {
	case token := <-tokenChan:
		return token, nil
	case err := <-tokenError:
		return "", err
		// case <-time.After(5 * time.Second):
		// 	return "", fmt.Errorf("timeout fetching token")
	}
}
