package cmd

import (
	"dehasher/internal/export"
	"dehasher/internal/files"
	"dehasher/internal/pretty"
	"dehasher/internal/sqlite"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"strings"
	"time"
)

var (
	// DB command flags
	dbPath string

	// DB query command flags
	usernameDBQuery              string
	emailDBQuery                 string
	ipDBQuery                    string
	passwordDBQuery              string
	hashDBQuery                  string
	nameDBQuery                  string
	vinDBQuery                   string
	licensePlateDBQuery          string
	addressDBQuery               string
	phoneDBQuery                 string
	socialDBQuery                string
	cryptoCurrencyAddressDBQuery string
	domainDBQuery                string
	limitResultsDB               int
	exactMatchDBQuery            bool
	outputFormatDB               string
	nonEmptyFieldsDBQuery        string
	displayFieldsDBQuery         string
	tableTypeDBQuery             string

	// DB runs command flags
	startDateDBRuns     string
	endDateDBRuns       string
	containsQueryDBRuns string
	lastXRunsDBRuns     int

	// DB command
	dbCmd = &cobra.Command{
		Use:   "db",
		Short: "Database operations for Dehasher",
		Long:  `Perform database operations like export, import, and query on the local Dehasher database.`,
	}
)

func init() {
	// Add subcommands to db command
	dbCmd.AddCommand(dbExportCmd)
	dbCmd.AddCommand(dbQueryCmd)
	dbCmd.AddCommand(dbRunsCmd)
	dbCmd.AddCommand(dbCredsCmd)

	// Add flags specific to db command
	dbCmd.PersistentFlags().StringVarP(&dbPath, "db-path", "D", "", "Path to database (default: ~/.local/share/Dehasher/dehashed.db)")

	// Add flags specific to db query command
	dbQueryCmd.Flags().StringVarP(&usernameDBQuery, "username", "u", "", "Filter by username")
	dbQueryCmd.Flags().StringVarP(&emailDBQuery, "email", "e", "", "Filter by email")
	dbQueryCmd.Flags().StringVarP(&ipDBQuery, "ip", "i", "", "Filter by IP address")
	dbQueryCmd.Flags().StringVarP(&passwordDBQuery, "password", "p", "", "Filter by password")
	dbQueryCmd.Flags().StringVarP(&hashDBQuery, "hash", "H", "", "Filter by hashed password")
	dbQueryCmd.Flags().StringVarP(&nameDBQuery, "name", "n", "", "Filter by name")
	dbQueryCmd.Flags().StringVarP(&vinDBQuery, "vin", "v", "", "Filter by VIN")
	dbQueryCmd.Flags().StringVarP(&licensePlateDBQuery, "license", "L", "", "Filter by license plate")
	dbQueryCmd.Flags().StringVarP(&addressDBQuery, "address", "a", "", "Filter by address")
	dbQueryCmd.Flags().StringVarP(&phoneDBQuery, "phone", "P", "", "Filter by phone number")
	dbQueryCmd.Flags().StringVarP(&socialDBQuery, "social", "s", "", "Filter by social media handle")
	dbQueryCmd.Flags().StringVarP(&cryptoCurrencyAddressDBQuery, "crypto", "c", "", "Filter by cryptocurrency address")
	dbQueryCmd.Flags().StringVarP(&domainDBQuery, "domain", "d", "", "Filter by domain/URL")
	dbQueryCmd.Flags().IntVarP(&limitResultsDB, "limit", "l", 100, "Limit number of results")
	dbQueryCmd.Flags().BoolVarP(&exactMatchDBQuery, "exact", "x", false, "Use exact matching instead of partial matching")
	dbQueryCmd.Flags().StringVarP(&outputFormatDB, "format", "f", "table", "Output format (json, table, simple)")
	dbQueryCmd.Flags().StringVar(&nonEmptyFieldsDBQuery, "non-empty", "", "Filter for non-empty fields (comma-separated list, e.g., 'password,email')")
	dbQueryCmd.Flags().StringVar(&displayFieldsDBQuery, "display", "", "Fields to display in output (comma-separated list, e.g., 'username,email,password')")
	dbQueryCmd.Flags().StringVarP(&tableTypeDBQuery, "table", "t", "results", "Table to query (results, runs, creds)")

	// Add flags specific to db runs command
	dbRunsCmd.Flags().StringVarP(&startDateDBRuns, "start-date", "s", "", "Start date for filtering runs (format: YYYY-MM-DD)")
	dbRunsCmd.Flags().StringVarP(&endDateDBRuns, "end-date", "e", "", "End date for filtering runs (format: YYYY-MM-DD)")
	dbRunsCmd.Flags().StringVarP(&containsQueryDBRuns, "contains", "c", "", "Filter runs containing this query string")
	dbRunsCmd.Flags().IntVarP(&lastXRunsDBRuns, "last", "x", 0, "Show the last X runs")
	dbRunsCmd.Flags().IntVarP(&limitResultsDB, "limit", "l", 100, "Limit number of results")
	dbRunsCmd.Flags().StringVarP(&outputFormatDB, "format", "f", "table", "Output format (json, table, simple)")

	// Add flags specific to db creds command
	dbCredsCmd.Flags().StringVarP(&usernameDBQuery, "username", "u", "", "Filter by username")
	dbCredsCmd.Flags().StringVarP(&emailDBQuery, "email", "e", "", "Filter by email")
	dbCredsCmd.Flags().StringVarP(&passwordDBQuery, "password", "p", "", "Filter by password")
	dbCredsCmd.Flags().IntVarP(&limitResultsDB, "limit", "l", 100, "Limit number of results")
	dbCredsCmd.Flags().BoolVarP(&exactMatchDBQuery, "exact", "x", false, "Use exact matching instead of partial matching")
	dbCredsCmd.Flags().StringVarP(&outputFormatDB, "format", "f", "table", "Output format (json, table, simple)")
	dbCredsCmd.Flags().StringVar(&nonEmptyFieldsDBQuery, "non-empty", "", "Filter for non-empty fields (comma-separated list, e.g., 'password,email')")
	dbCredsCmd.Flags().StringVar(&displayFieldsDBQuery, "display", "", "Fields to display in output (comma-separated list, e.g., 'username,email,password')")
}

// DB export command
var dbExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export database to file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Exporting database...")
		// Create DBOptions with the provided parameters
		options := &sqlite.DBOptions{
			Username:       usernameDBQuery,
			Email:          emailDBQuery,
			IPAddress:      ipDBQuery,
			Password:       passwordDBQuery,
			HashedPassword: hashDBQuery,
			Name:           nameDBQuery,
			Limit:          limitResultsDB,
			ExactMatch:     exactMatchDBQuery,
		}

		// Parse non-empty fields if provided
		if nonEmptyFieldsDBQuery != "" {
			options.NonEmptyFields = strings.Split(nonEmptyFieldsDBQuery, ",")
		}

		// Parse display fields if provided
		if displayFieldsDBQuery != "" {
			options.DisplayFields = strings.Split(displayFieldsDBQuery, ",")
		}

		// Check if at least one search parameter is provided
		if options.Username == "" && options.Email == "" && options.IPAddress == "" &&
			options.Password == "" && options.HashedPassword == "" && options.Name == "" &&
			len(options.NonEmptyFields) == 0 {
			fmt.Println("Error: At least one search parameter is required.")
			cmd.Help()
			return
		}

		// Get the count of matching results
		count, err := sqlite.GetResultsCount(options)
		if err != nil {
			fmt.Printf("Error counting results: %v\n", err)
			return
		}

		// Query the database
		results, err := sqlite.QueryResults(options)
		if err != nil {
			fmt.Printf("Error querying database: %v\n", err)
			return
		}
		dhResults := sqlite.DehashedResults{Results: results}

		fmt.Printf("Found %d results (showing %d):\n", count, len(results))

		// Output results based on format
		ft := files.GetFileType(outputFormatDB)
		err = export.WriteToFile(dhResults, "dehasher_export", ft)
		if err != nil {
			zap.L().Error("write_to_file",
				zap.String("message", "failed to write to file"),
				zap.Error(err),
			)
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}
		fmt.Printf("Exported successfully to file: dehasher_export%s\n", ft.Extension())
	},
}

// DB query command
var dbQueryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query local database",
	Long:  `Query the local database for previously run dehasher queries based on various parameters.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Determine which table to query based on the tableTypeDBQuery parameter
		switch tableTypeDBQuery {
		case "results":
			queryResultsTable(cmd)
		case "runs":
			queryRunsTable()
		case "creds":
			queryCredsTable(cmd)
		default:
			fmt.Printf("Error: Unknown table type '%s'. Valid options are: results, runs, creds\n", tableTypeDBQuery)
			cmd.Help()
			return
		}
	},
}

func queryRunsTable() {
	// Parse date strings to time.Time
	var startDate, endDate time.Time
	var err error

	if startDateDBRuns != "" {
		startDate, err = time.Parse("2006-01-02", startDateDBRuns)
		if err != nil {
			fmt.Printf("Error parsing start date: %v\n", err)
			return
		}
	}

	if endDateDBRuns != "" {
		endDate, err = time.Parse("2006-01-02", endDateDBRuns)
		if err != nil {
			fmt.Printf("Error parsing end date: %v\n", err)
			return
		}
		// Set end date to end of day
		endDate = endDate.Add(24*time.Hour - time.Second)
	}

	// Get the count of matching runs
	count, err := sqlite.GetRunsCount(lastXRunsDBRuns, startDate, endDate, containsQueryDBRuns)
	if err != nil {
		fmt.Printf("Error counting runs: %v\n", err)
		return
	}

	// Query the database
	runs, err := sqlite.QueryRuns(limitResultsDB, lastXRunsDBRuns, startDate, endDate, containsQueryDBRuns)
	if err != nil {
		fmt.Printf("Error querying runs: %v\n", err)
		return
	}

	displayRunsResults(count, runs)
}

func displayRunsResults(count int64, runs []sqlite.QueryOptions) {
	// Display the results
	fmt.Printf("Found %d runs (showing %d):\n", count, len(runs))

	if len(runs) == 0 {
		fmt.Println("No runs found.")
		return
	}

	// Output results based on format
	switch outputFormatDB {
	case "json":
		data, err := json.MarshalIndent(runs, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting results: %v\n", err)
			return
		}
		fmt.Println(string(data))
	case "table":
		// Define headers and rows for the table
		headers := []string{"ID", "Created At", "Max Records", "Username Query", "Email Query", "IP Query", "Password Query", "Hash Query", "Name Query", "Domain Query"}
		rows := make([][]string, len(runs))

		for i, run := range runs {
			rows[i] = []string{
				fmt.Sprintf("%d", run.ID),
				run.CreatedAt.Format("2006-01-02 15:04:05"),
				fmt.Sprintf("%d", run.MaxRecords),
				truncate(run.UsernameQuery, 20),
				truncate(run.EmailQuery, 20),
				truncate(run.IpQuery, 20),
				truncate(run.PassQuery, 20),
				truncate(run.HashQuery, 20),
				truncate(run.NameQuery, 20),
				truncate(run.DomainQuery, 20),
			}
		}

		pretty.Table(headers, rows)
	default:
		// Simple output
		for _, run := range runs {
			fmt.Printf("Run ID: %d\n", run.ID)
			fmt.Printf("  Created At: %s\n", run.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("  Max Records: %d\n", run.MaxRecords)
			fmt.Printf("  Max Requests: %d\n", run.MaxRequests)
			fmt.Printf("  Starting Page: %d\n", run.StartingPage)
			fmt.Printf("  Output Format: %s\n", run.OutputFormat.String())
			fmt.Printf("  Output File: %s\n", run.OutputFile)
			fmt.Printf("  Regex Match: %t\n", run.RegexMatch)
			fmt.Printf("  Wildcard Match: %t\n", run.WildcardMatch)
			fmt.Printf("  Username Query: %s\n", run.UsernameQuery)
			fmt.Printf("  Email Query: %s\n", run.EmailQuery)
			fmt.Printf("  IP Query: %s\n", run.IpQuery)
			fmt.Printf("  Password Query: %s\n", run.PassQuery)
			fmt.Printf("  Hash Query: %s\n", run.HashQuery)
			fmt.Printf("  Name Query: %s\n", run.NameQuery)
			fmt.Printf("  Domain Query: %s\n", run.DomainQuery)
			fmt.Printf("  VIN Query: %s\n", run.VinQuery)
			fmt.Printf("  License Plate Query: %s\n", run.LicensePlateQuery)
			fmt.Printf("  Address Query: %s\n", run.AddressQuery)
			fmt.Printf("  Phone Query: %s\n", run.PhoneQuery)
			fmt.Printf("  Social Query: %s\n", run.SocialQuery)
			fmt.Printf("  Crypto Address Query: %s\n", run.CryptoAddressQuery)
			fmt.Printf("  Print Balance: %t\n", run.PrintBalance)
			fmt.Printf("  Creds Only: %t\n", run.CredsOnly)
			fmt.Println()
		}
	}
}

// queryResultsTable queries the results table
func queryResultsTable(cmd *cobra.Command) {
	// Create DBOptions with the provided parameters
	options := &sqlite.DBOptions{
		Username:              usernameDBQuery,
		Email:                 emailDBQuery,
		IPAddress:             ipDBQuery,
		Password:              passwordDBQuery,
		HashedPassword:        hashDBQuery,
		Name:                  nameDBQuery,
		Vin:                   vinDBQuery,
		LicensePlate:          licensePlateDBQuery,
		Address:               addressDBQuery,
		Phone:                 phoneDBQuery,
		Social:                socialDBQuery,
		CryptoCurrencyAddress: cryptoCurrencyAddressDBQuery,
		Domain:                domainDBQuery,
		Limit:                 limitResultsDB,
		ExactMatch:            exactMatchDBQuery,
	}

	// Parse non-empty fields if provided
	if nonEmptyFieldsDBQuery != "" {
		options.NonEmptyFields = strings.Split(nonEmptyFieldsDBQuery, ",")
	}

	// Parse display fields if provided
	if displayFieldsDBQuery != "" {
		options.DisplayFields = strings.Split(displayFieldsDBQuery, ",")
	}

	// Check if at least one search parameter is provided
	if options.Empty() {
		fmt.Println("Error: At least one search parameter is required.")
		cmd.Help()
		return
	}

	// Get the count of matching results
	count, err := sqlite.GetResultsCount(options)
	if err != nil {
		fmt.Printf("Error counting results: %v\n", err)
		return
	}

	// Query the database
	results, err := sqlite.QueryResults(options)
	if err != nil {
		fmt.Printf("Error querying database: %v\n", err)
		return
	}

	// Display the results
	displayResultsTable(count, results, options)
}

// displayResultsTable displays the results from the results table
func displayResultsTable(count int64, results []sqlite.Result, options *sqlite.DBOptions) {
	// Display the results
	fmt.Printf("Found %d results (showing %d):\n", count, len(results))

	if len(results) == 0 {
		fmt.Println("No results found.")
		return
	}

	// Output results based on format
	switch outputFormatDB {
	case "json":
		data, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting results: %v\n", err)
			return
		}
		fmt.Println(string(data))
	case "table":
		// Determine which fields to display
		type FieldInfo struct {
			Name   string
			Width  int
			Getter func(result sqlite.Result) string
		}

		// Define all available fields
		allFields := []FieldInfo{
			{"Username", 20, func(r sqlite.Result) string { return truncate(arrayToString(r.Username), 20) }},
			{"Email", 30, func(r sqlite.Result) string { return truncate(arrayToString(r.Email), 30) }},
			{"IP Address", 15, func(r sqlite.Result) string { return truncate(arrayToString(r.IpAddress), 15) }},
			{"Password", 20, func(r sqlite.Result) string { return truncate(arrayToString(r.Password), 20) }},
			{"Hashed Password", 20, func(r sqlite.Result) string { return truncate(arrayToString(r.HashedPassword), 20) }},
			{"Name", 20, func(r sqlite.Result) string { return truncate(arrayToString(r.Name), 20) }},
			{"VIN", 20, func(r sqlite.Result) string { return truncate(arrayToString(r.Vin), 20) }},
			{"License Plate", 15, func(r sqlite.Result) string { return truncate(arrayToString(r.LicensePlate), 15) }},
			{"Address", 30, func(r sqlite.Result) string { return truncate(arrayToString(r.Address), 30) }},
			{"Phone", 15, func(r sqlite.Result) string { return truncate(arrayToString(r.Phone), 15) }},
			{"Social", 20, func(r sqlite.Result) string { return truncate(arrayToString(r.Social), 20) }},
			{"Crypto Address", 20, func(r sqlite.Result) string { return truncate(arrayToString(r.CryptoCurrencyAddress), 20) }},
			{"Domain/URL", 30, func(r sqlite.Result) string { return truncate(arrayToString(r.Url), 30) }},
		}

		// Select fields to display
		var fieldsToDisplay []FieldInfo
		var headers []string
		if len(options.DisplayFields) > 0 {
			// Use specified fields
			for _, fieldName := range options.DisplayFields {
				fieldName = strings.ToLower(strings.TrimSpace(fieldName))
				for _, field := range allFields {
					if strings.ToLower(field.Name) == fieldName ||
						(fieldName == "ip" && strings.ToLower(field.Name) == "ip address") ||
						(fieldName == "hash" && strings.ToLower(field.Name) == "hashed password") ||
						(fieldName == "license" && strings.ToLower(field.Name) == "license plate") ||
						(fieldName == "crypto" && strings.ToLower(field.Name) == "crypto address") ||
						(fieldName == "url" && strings.ToLower(field.Name) == "domain/url") {
						fieldsToDisplay = append(fieldsToDisplay, field)
						headers = append(headers, field.Name)
						break
					}
				}
			}
		} else {
			// Default fields (first 6)
			fieldsToDisplay = allFields[:6]
		}

		var rows [][]string
		for _, result := range results {
			rowValues := []string{}
			for _, field := range fieldsToDisplay {
				rowValues = append(rowValues, field.Getter(result))
			}
			rows = append(rows, rowValues)
		}

		pretty.Table(headers, rows)
	default:
		// Simple output
		for i, result := range results {
			fmt.Printf("Result %d:\n", i+1)

			// Determine which fields to display
			if len(options.DisplayFields) > 0 {
				// Display only specified fields
				for _, field := range options.DisplayFields {
					field = strings.ToLower(strings.TrimSpace(field))
					switch field {
					case "username":
						fmt.Printf("  Username: %s\n", result.Username)
					case "email":
						fmt.Printf("  Email: %s\n", result.Email)
					case "ip", "ipaddress", "ip_address":
						fmt.Printf("  IP Address: %s\n", result.IpAddress)
					case "password":
						fmt.Printf("  Password: %s\n", result.Password)
					case "hash", "hashed_password":
						fmt.Printf("  Hashed Password: %s\n", result.HashedPassword)
					case "name":
						fmt.Printf("  Name: %s\n", result.Name)
					case "vin":
						fmt.Printf("  VIN: %s\n", result.Vin)
					case "license", "license_plate":
						fmt.Printf("  License Plate: %s\n", result.LicensePlate)
					case "address":
						fmt.Printf("  Address: %s\n", result.Address)
					case "phone":
						fmt.Printf("  Phone: %s\n", result.Phone)
					case "social":
						fmt.Printf("  Social: %s\n", result.Social)
					case "crypto", "cryptocurrency_address":
						fmt.Printf("  Crypto Address: %s\n", result.CryptoCurrencyAddress)
					case "domain", "url":
						fmt.Printf("  Domain/URL: %s\n", result.Url)
					}
				}
			} else {
				// Display default fields
				fmt.Printf("  Username: %s\n", result.Username)
				fmt.Printf("  Email: %s\n", result.Email)
				fmt.Printf("  IP Address: %s\n", result.IpAddress)
				fmt.Printf("  Password: %s\n", result.Password)
				fmt.Printf("  Hashed Password: %s\n", result.HashedPassword)
				fmt.Printf("  Name: %s\n", result.Name)
			}
			fmt.Println()
		}
	}
}

// truncate truncates a string to the specified length and adds ellipsis if needed
func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

func arrayToString(a []string) string {
	return strings.Join(a, ", ")
}

// DB runs command
var dbRunsCmd = &cobra.Command{
	Use:   "runs",
	Short: "Query previous query runs",
	Long:  `Query the database for previous query runs (QueryOptions) based on date range and query content.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Call queryRunsTable directly
		queryRunsTable()
	},
}

// DB creds command
var dbCredsCmd = &cobra.Command{
	Use:   "creds",
	Short: "Query credentials",
	Long:  `Query the database for credentials based on username, email, and password.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Call queryCredsTable directly
		queryCredsTable(cmd)
	},
}

// queryCredsTable queries the credentials table
func queryCredsTable(cmd *cobra.Command) {
	// Create DBOptions with the provided parameters
	options := &sqlite.DBOptions{
		Username:   usernameDBQuery,
		Email:      emailDBQuery,
		Password:   passwordDBQuery,
		Limit:      limitResultsDB,
		ExactMatch: exactMatchDBQuery,
	}

	// Parse non-empty fields if provided
	if nonEmptyFieldsDBQuery != "" {
		options.NonEmptyFields = strings.Split(nonEmptyFieldsDBQuery, ",")
	}

	// Parse display fields if provided
	if displayFieldsDBQuery != "" {
		options.DisplayFields = strings.Split(displayFieldsDBQuery, ",")
	}

	// Check if at least one search parameter is provided
	if options.Username == "" && options.Email == "" && options.Password == "" && len(options.NonEmptyFields) == 0 {
		fmt.Println("Error: At least one search parameter is required.")
		cmd.Help()
		return
	}

	// Get the count of matching credentials
	count, err := sqlite.GetCredsCount(options)
	if err != nil {
		fmt.Printf("Error counting credentials: %v\n", err)
		return
	}

	// Query the database
	creds, err := sqlite.QueryCreds(options)
	if err != nil {
		fmt.Printf("Error querying credentials: %v\n", err)
		return
	}

	// Display the results
	displayCredsResults(count, creds)
}

// displayCredsResults displays the results from the creds table
func displayCredsResults(count int64, creds []sqlite.Creds) {
	// Display the results
	fmt.Printf("Found %d credentials (showing %d):\n", count, len(creds))

	if len(creds) == 0 {
		fmt.Println("No credentials found.")
		return
	}

	// Output results based on format
	switch outputFormatDB {
	case "json":
		data, err := json.MarshalIndent(creds, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting results: %v\n", err)
			return
		}
		fmt.Println(string(data))
	case "table":
		// Define all available fields
		type FieldInfo struct {
			Name   string
			Getter func(cred sqlite.Creds) string
		}

		allFields := []FieldInfo{
			{"ID", func(c sqlite.Creds) string { return fmt.Sprintf("%d", c.ID) }},
			{"Created At", func(c sqlite.Creds) string { return c.CreatedAt.Format("2006-01-02 15:04:05") }},
			{"Email", func(c sqlite.Creds) string { return c.Email }},
			{"Username", func(c sqlite.Creds) string { return c.Username }},
			{"Password", func(c sqlite.Creds) string { return c.Password }},
		}

		// Select fields to display
		var fieldsToDisplay []FieldInfo
		var headers []string

		if len(displayFieldsDBQuery) > 0 {
			// Use specified display fields
			displayFields := strings.Split(displayFieldsDBQuery, ",")
			for _, fieldName := range displayFields {
				fieldName = strings.ToLower(strings.TrimSpace(fieldName))
				for _, field := range allFields {
					if strings.ToLower(field.Name) == fieldName {
						fieldsToDisplay = append(fieldsToDisplay, field)
						headers = append(headers, field.Name)
						break
					}
				}
			}
		} else {
			// Default fields
			fieldsToDisplay = allFields
			for _, field := range fieldsToDisplay {
				headers = append(headers, field.Name)
			}
		}

		// Create rows
		rows := make([][]string, len(creds))
		for i, cred := range creds {
			rowValues := []string{}
			for _, field := range fieldsToDisplay {
				rowValues = append(rowValues, field.Getter(cred))
			}
			rows[i] = rowValues
		}

		pretty.Table(headers, rows)
	default:
		// Simple output
		for _, cred := range creds {
			fmt.Printf("Credential ID: %d\n", cred.ID)
			fmt.Printf("  Created At: %s\n", cred.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("  Email: %s\n", cred.Email)
			fmt.Printf("  Username: %s\n", cred.Username)
			fmt.Printf("  Password: %s\n", cred.Password)
			fmt.Println()
		}
	}
}
