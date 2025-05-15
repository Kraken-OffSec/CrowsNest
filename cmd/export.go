package cmd

import (
	"dehasher/internal/export"
	"dehasher/internal/files"
	"dehasher/internal/sqlite"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"strings"
)

func init() {
	// Add Subcommand to db command
	rootCmd.AddCommand(exportCmd)

	// Add flags specific to export command
	exportCmd.Flags().IntVarP(&exportLimitRows, "limit", "l", 100, "Limit number of results")
	exportCmd.Flags().BoolVarP(&exportListAll, "list-all", "a", false, "List all tables and their columns")
	exportCmd.Flags().StringVarP(&exportTableName, "table", "t", "", "Table to export (results, creds, whois, subdomains, history, runs)")
	exportCmd.Flags().StringVarP(&exportNotNull, "not-null", "n", "", "Filter for non-null values (comma-separated list, e.g., 'password,email')")
	exportCmd.Flags().StringVarP(&exportColumns, "columns", "c", "", "Columns to display in output (comma-separated list, e.g., 'username,email,password')")
	exportCmd.Flags().StringVarP(&exportUserQuery, "user-query", "q", "", "User query to execute")
	exportCmd.Flags().StringVarP(&exportRawQuery, "raw-query", "r", "", "Raw SQL query to execute")
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "json", "Output format (json, yaml, xml, txt)")
	exportCmd.Flags().StringVarP(&exportFile, "file", "o", "export", "File to output results to including extension")

	// Add mutually exclusive flags to query and raw-query
	// Cannot use query and raw-query at the same time
	exportCmd.MarkFlagsMutuallyExclusive("user-query", "raw-query")
	// Raw query does not require a table
	exportCmd.MarkFlagsMutuallyExclusive("user-query", "table")
	// List all columns does not require a query or raw-query
	exportCmd.MarkFlagsMutuallyExclusive("raw-query", "list-all")
}

// DB export command
var (
	exportLimitRows int
	exportListAll   bool
	exportTableName string
	exportNotNull   string
	exportColumns   string
	exportUserQuery string
	exportRawQuery  string
	exportFormat    string
	exportFile      string

	exportCmd = &cobra.Command{
		Use:   "export",
		Short: "Export database to file",
		Run: func(cmd *cobra.Command, args []string) {
			// If list-all flag is set, list all tables and columns
			if exportListAll {
				listAvailableTables()
				return
			}

			fmt.Println("[*] Exporting database...")

			// If Raw Query is set, execute it and export
			if exportRawQuery != "" {
				fmt.Println("[*] Executing Raw Query...")
				exportRawDBQuery()
				return
			}

			// Validate table name
			if exportTableName == "" {
				fmt.Println("[!] Error: Table name is required. Use -t or --table to specify a table.")
				fmt.Println("[*] Available tables: results, creds, whois, subdomains, history, runs")
				fmt.Println("[*] Use --list-all to see all tables and their columns.")
				return
			}

			if !isValidTable(exportTableName) {
				fmt.Printf("[!] Error: Unknown table '%s'.\n", exportTableName)
				fmt.Println("[*] Available tables: results, creds, whois, subdomains, history, runs")
				fmt.Println("[*] Use --list-all to see all tables and their columns.")
				return
			}

			// Validate columns if specified
			if exportColumns != "" {
				columns := strings.Split(exportColumns, ",")
				invalidColumns := validateColumns(exportTableName, columns)
				if len(invalidColumns) > 0 {
					fmt.Printf("[!] Error: Invalid column(s) for table '%s': %s\n",
						exportTableName, strings.Join(invalidColumns, ", "))
					fmt.Println("[*] Available columns for this table:")
					for i := 0; i < len(availableTables[exportTableName]); i += 5 {
						end := i + 5
						if end > len(availableTables[exportTableName]) {
							end = len(availableTables[exportTableName])
						}
						fmt.Printf("    %s\n", strings.Join(availableTables[exportTableName][i:end], ", "))
					}
					return
				}
			}

			// Validate not-null fields if specified
			if exportNotNull != "" {
				notNullFields := strings.Split(exportNotNull, ",")
				invalidFields := validateColumns(exportTableName, notNullFields)
				if len(invalidFields) > 0 {
					fmt.Printf("[!] Error: Invalid not-null field(s) for table '%s': %s\n",
						exportTableName, strings.Join(invalidFields, ", "))
					fmt.Println("[*] Available columns for this table:")
					for i := 0; i < len(availableTables[exportTableName]); i += 5 {
						end := i + 5
						if end > len(availableTables[exportTableName]) {
							end = len(availableTables[exportTableName])
						}
						fmt.Printf("    %s\n", strings.Join(availableTables[exportTableName][i:end], ", "))
					}
					return
				}
			}

			// Determine which table to query based on the tableTypeDBQuery parameter
			table := sqlite.GetTable(exportTableName)
			if table == sqlite.UnknownTable {
				fmt.Printf("[!] Error: Unknown table type '%s'.\n", exportTableName)
				fmt.Println("[*] Available tables: results, creds, whois, subdomains, history, runs")
				fmt.Println("[*] Use --list-all to see all tables and their columns.")
				return
			}

			fmt.Println("[*] Querying Database...")
			exportTableQuery(table)
		},
	}
)

// exportTableQuery queries a table and exports the results
func exportTableQuery(table sqlite.Table) {
	// Get the columns to query
	columns := []string{"*"}
	if exportColumns != "" {
		columns = strings.Split(exportColumns, ",")
	}

	// Get the not null fields
	notNullFields := []string{}
	if exportNotNull != "" {
		notNullFields = strings.Split(exportNotNull, ",")
	}

	// Get the user query
	userQuery := ""
	if exportUserQuery != "" {
		userQuery = exportUserQuery
	}

	// Get the limit
	limit := exportLimitRows

	// Get the object for the table
	object := table.Object()

	// Check if object is nil (invalid table)
	if object == nil {
		fmt.Printf("[!] Error: Table '%s' is not valid or does not exist.\n", exportTableName)
		return
	}

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
		zap.L().Error("export_query",
			zap.String("message", "failed to execute query"),
			zap.Error(err),
		)
		fmt.Printf("[!] Error executing query: %v\n", err)
		return
	}
	defer rows.Close()

	// Get the columns
	cols, err := rows.Columns()
	if err != nil {
		zap.L().Error("export_query",
			zap.String("message", "failed to get columns from query"),
			zap.Error(err),
		)
		fmt.Printf("[!] Error getting columns from query: %v\n", err)
		return
	}

	// Prepare data for export
	var results []map[string]interface{}

	// Process the rows
	for rows.Next() {
		values := make([]interface{}, len(cols))
		pointers := make([]interface{}, len(cols))
		for i := range values {
			pointers[i] = &values[i]
		}
		if err := rows.Scan(pointers...); err != nil {
			zap.L().Error("export_query",
				zap.String("message", "failed to scan row from query"),
				zap.Error(err),
			)
			fmt.Printf("[!] Error scanning row from query: %v\n", err)
			return
		}

		// Create a map for this row
		rowMap := make(map[string]interface{})
		for i, col := range cols {
			val := values[i]
			rowMap[col] = val
		}

		results = append(results, rowMap)
	}

	// Export the results
	exportResults(results)
}

// exportRawDBQuery executes a raw query and exports the results
func exportRawDBQuery() {
	db := sqlite.GetDB()
	rows, err := db.Raw(exportRawQuery).Rows()
	if err != nil {
		zap.L().Error("export_raw_query",
			zap.String("message", "failed to execute raw query"),
			zap.Error(err),
		)
		fmt.Printf("[!] Error executing raw query: %v\n", err)
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		zap.L().Error("export_raw_query",
			zap.String("message", "failed to get columns from raw query"),
			zap.Error(err),
		)
		fmt.Printf("[!] Error getting columns from raw query: %v\n", err)
		return
	}

	// Prepare data for export
	var results []map[string]interface{}

	// Process the rows
	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		if err := rows.Scan(pointers...); err != nil {
			zap.L().Error("export_raw_query",
				zap.String("message", "failed to scan row from raw query"),
				zap.Error(err),
			)
			fmt.Printf("[!] Error scanning row from raw query: %v\n", err)
			return
		}

		// Create a map for this row
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			rowMap[col] = val
		}

		results = append(results, rowMap)
	}

	// Export the results
	exportResults(results)
}

// exportResults exports the results to a file
func exportResults(results []map[string]interface{}) {
	// Get file type
	fileType := files.GetFileType(exportFormat)

	// Export results
	err := export.WriteQueryResultsToFile(results, exportFile, fileType)
	if err != nil {
		zap.L().Error("export_results",
			zap.String("message", "failed to write to file"),
			zap.Error(err),
		)
		fmt.Printf("[!] Error writing to file: %v\n", err)
		return
	}

	fmt.Printf("[+] Exported %d records to file: %s%s\n", len(results), exportFile, fileType.Extension())
}
