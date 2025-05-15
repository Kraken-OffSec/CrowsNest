package cmd

import (
	"dehasher/internal/pretty"
	"dehasher/internal/sqlite"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"strings"
)

func init() {
	// Add whois command to root command
	rootCmd.AddCommand(queryCmd)

	// Add flags specific to whois command
	queryCmd.Flags().StringVarP(&dbQueryTableName, "table", "t", "", "Table to query (results, creds, whois, subdomains, history, query_options)")
	queryCmd.Flags().IntVarP(&dbQueryLimitRows, "limit", "l", 100, "Limit number of results")
	queryCmd.Flags().StringVarP(&dbQueryNotNull, "not-null", "n", "", "Filter for non-null values (comma-separated list, e.g., 'password,email')")
	queryCmd.Flags().StringVarP(&dbQueryColumns, "columns", "c", "", "Columns to display in output (comma-separated list, e.g., 'username,email,password')")
	queryCmd.Flags().StringVarP(&dbQueryUserQuery, "user-query", "q", "", "User query to execute")
	queryCmd.Flags().StringVarP(&dbQueryRawQuery, "raw-query", "r", "", "Raw SQL query to execute")
	queryCmd.Flags().BoolVarP(&dbQueryListAll, "list-all", "a", false, "List all columns")

	// Add mutually exclusive flags to query and raw-query
	// Cannot use query and raw-query at the same time
	queryCmd.MarkFlagsMutuallyExclusive("user-query", "raw-query")
	// Raw query does not require a table
	queryCmd.MarkFlagsMutuallyExclusive("user-query", "table")
	// List all columns does not require a query or raw-query
	queryCmd.MarkFlagsMutuallyExclusive("raw-query", "list-all")
}

var (
	dbQueryTableName string
	dbQueryLimitRows int
	dbQueryNotNull   string
	dbQueryColumns   string
	dbQueryUserQuery string
	dbQueryRawQuery  string
	dbQueryListAll   bool

	queryCmd = &cobra.Command{
		Use:   "query",
		Short: "Query the database",
		Long:  `Query the database for various information.`,
		Run: func(cmd *cobra.Command, args []string) {
			// If Raw Query is set, execute it and return
			if dbQueryRawQuery != "" {
				fmt.Println("[*] Executing Raw Query...")
				rawDBQuery()
				os.Exit(1)
			}

			// Determine which table to query based on the tableTypeDBQuery parameter
			table := sqlite.GetTable(dbQueryTableName)
			if table == sqlite.UnknownTable {
				fmt.Printf("Error: Unknown table type '%s'.\n", dbQueryTableName)
				cmd.Help()
				return
			}
			fmt.Println("[*] Querying Database...")
			tableQuery(table)
		},
	}
)

func tableQuery(table sqlite.Table) {

	// Get the columns to query
	columns := []string{"*"}
	if dbQueryColumns != "" {
		columns = strings.Split(dbQueryColumns, ",")
	}

	// Get the not null fields
	notNullFields := []string{}
	if dbQueryNotNull != "" {
		notNullFields = strings.Split(dbQueryNotNull, ",")
	}

	// Get the user query
	userQuery := ""
	if dbQueryUserQuery != "" {
		userQuery = dbQueryUserQuery
	}

	// Get the limit
	limit := dbQueryLimitRows

	// Get the object for the table
	object := table.Object()

	// Query the database
	db := sqlite.GetDB()
	query := db.Model(object).Select(columns)
	if len(notNullFields) > 0 {
		for _, field := range notNullFields {
			query = query.Where(fmt.Sprintf("%s IS NOT NULL", field))
		}
	}
	if userQuery != "" {
		query = query.Where(userQuery)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	rows, err := query.Rows()
	if err != nil {
		zap.L().Error("db_query",
			zap.String("message", "failed to execute query"),
			zap.Error(err),
		)
		fmt.Printf("[!] Error executing query: %v\n", err)
	}
	defer rows.Close()

	// Get the columns
	cols, err := rows.Columns()
	if err != nil {
		zap.L().Error("db_query",
			zap.String("message", "failed to get columns from query"),
			zap.Error(err),
		)
		fmt.Printf("[!] Error getting columns from query: %v\n", err)
	}

	// Prepare data for pretty.Table
	headers := cols
	var tableRows [][]string

	// Process the rows
	for rows.Next() {
		values := make([]interface{}, len(cols))
		pointers := make([]interface{}, len(cols))
		for i := range values {
			pointers[i] = &values[i]
		}
		if err := rows.Scan(pointers...); err != nil {
			zap.L().Error("db_query",
				zap.String("message", "failed to scan row from query"),
				zap.Error(err),
			)
			fmt.Printf("[!] Error scanning row from query: %v\n", err)
		}

		// Convert row values to strings
		rowStrings := make([]string, len(values))
		for i, value := range values {
			if value == nil {
				rowStrings[i] = " "
			} else {
				// Check if the value is a slice or array
				switch v := value.(type) {
				case []string:
					// Join string slices with commas, no brackets
					rowStrings[i] = strings.Join(v, ", ")
				case []interface{}:
					// Convert interface slice to strings and join
					strSlice := make([]string, len(v))
					for j, item := range v {
						if item == nil {
							strSlice[j] = ""
						} else {
							strSlice[j] = fmt.Sprintf("%v", item)
						}
					}
					rowStrings[i] = strings.Join(strSlice, ", ")
				case string:
					// Handle JSON strings that might be arrays
					if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
						// Try to unmarshal JSON array
						var strArray []string
						if err := json.Unmarshal([]byte(v), &strArray); err == nil {
							rowStrings[i] = strings.Join(strArray, ", ")
						} else {
							rowStrings[i] = v
						}
					} else {
						rowStrings[i] = v
					}
				default:
					rowStrings[i] = fmt.Sprintf("%v", v)
				}
			}
		}
		tableRows = append(tableRows, rowStrings)
	}

	// Display the table
	pretty.Table(headers, tableRows)
}

func rawDBQuery() {
	db := sqlite.GetDB()
	rows, err := db.Raw(dbQueryRawQuery).Rows()
	if err != nil {
		zap.L().Error("raw_query",
			zap.String("message", "failed to execute raw query"),
			zap.Error(err),
		)
		fmt.Printf("[!] Error executing raw query: %v\n", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		zap.L().Error("raw_query",
			zap.String("message", "failed to get columns from raw query"),
			zap.Error(err),
		)
		fmt.Printf("[!] Error getting columns from raw query: %v\n", err)
	}

	// Prepare data for pretty.Table
	headers := columns
	var tableRows [][]string

	// Process the rows
	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		if err := rows.Scan(pointers...); err != nil {
			zap.L().Error("raw_query",
				zap.String("message", "failed to scan row from raw query"),
				zap.Error(err),
			)
			fmt.Printf("[!] Error scanning row from raw query: %v\n", err)
		}

		// Convert row values to strings
		rowStrings := make([]string, len(values))
		for i, value := range values {
			if value == nil {
				rowStrings[i] = " "
			} else {
				// Check if the value is a slice or array
				switch v := value.(type) {
				case []string:
					// Join string slices with commas, no brackets
					rowStrings[i] = strings.Join(v, ", ")
				case []interface{}:
					// Convert interface slice to strings and join
					strSlice := make([]string, len(v))
					for j, item := range v {
						if item == nil {
							strSlice[j] = ""
						} else {
							strSlice[j] = fmt.Sprintf("%v", item)
						}
					}
					rowStrings[i] = strings.Join(strSlice, ", ")
				case string:
					// Handle JSON strings that might be arrays
					if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
						// Try to unmarshal JSON array
						var strArray []string
						if err := json.Unmarshal([]byte(v), &strArray); err == nil {
							rowStrings[i] = strings.Join(strArray, ", ")
						} else {
							rowStrings[i] = v
						}
					} else {
						rowStrings[i] = v
					}
				default:
					rowStrings[i] = fmt.Sprintf("%v", v)
				}
			}
		}
		tableRows = append(tableRows, rowStrings)
	}

	// Display the table
	pretty.Table(headers, tableRows)
}
