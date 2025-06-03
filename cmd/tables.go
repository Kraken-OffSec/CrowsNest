package cmd

import (
	"crowsnest/internal/pretty"
	"fmt"
	"strings"
)

// Map of available tables and their columns
var availableTables = map[string][]string{
	"users": {
		"id", "created_at", "updated_at", "deleted_at", "email", "username", "password",
	},
	//"history": {
	//	"id", "created_at", "updated_at", "deleted_at", "domain_name", "domain_type",
	//	"registrar_name", "whois_server", "created_date_iso8601", "updated_date_iso8601", "expires_date_iso8601",
	//},
	"lookup": {
		"id", "created_at", "updated_at", "deleted_at", "search_term", "type", "first_seen", "last_visit",
		"name",
	},
	// Query Options
	"runs": {
		"id", "created_at", "updated_at", "deleted_at", "max_records", "max_requests", "starting_page",
		"output_format", "output_file", "regex_match", "wildcard_match", "username_query", "email_query",
		"ip_query", "pass_query", "hash_query", "name_query", "domain_query", "vin_query", "license_plate_query",
		"address_query", "phone_query", "social_query", "crypto_address_query", "print_balance", "creds_only",
	},
	"dehashed": {
		"id", "created_at", "updated_at", "deleted_at", "dehashed_id", "email", "ip_address", "username",
		"password", "hashed_password", "hash_type", "name", "vin", "license_plate", "url", "social",
		"cryptocurrency_address", "address", "phone", "company", "database_name",
	},
	"subdomains": {
		"id", "created_at", "updated_at", "deleted_at", "domain", "subdomain",
	},
	"whois": {
		"id", "created_at", "updated_at", "deleted_at", "audit", "contact_email", "created_date", "created_date_normalized",
		"domain_name", "domain_name_ext", "estimated_domain_age", "expires_date", "expires_date_normalized", "footer", "header",
		"name_servers", "parse_code", "raw_text", "registrant", "registrar_iana_id", "registrar_name", "registry_data",
		"status", "stripped_text", "updated_date", "updated_date_normalized",
	},
	"hunter_domain": {
		"id", "created_at", "updated_at", "deleted_at", "domain", "disposable", "webmail", "accept_all", "pattern",
		"organization", "description", "industry", "twitter", "facebook", "linkedin", "instagram", "youtube",
		"technologies", "country", "state", "city", "postal_code", "street", "headcount", "company_type", "emails", "linked_domains",
	},
	"hunter_email": {
		"id", "created_at", "updated_at", "deleted_at", "value", "type", "confidence", "sources", "first_name", "last_name",
		"position", "position_raw", "seniority", "department", "linkedin", "twitter", "phone_number", "verification_date", "verification_status",
	},
}

// Function to list available tables and their columns
func listAvailableTables() {
	fmt.Println("Available tables and columns:")

	// Prepare data for pretty.Table
	headers := []string{"Table", "Columns"}
	var tableRows [][]string

	// Sort tables alphabetically for consistent output
	var tableNames []string
	for tableName := range availableTables {
		tableNames = append(tableNames, tableName)
	}

	// Simple bubble sort for table names
	for i := 0; i < len(tableNames)-1; i++ {
		for j := 0; j < len(tableNames)-i-1; j++ {
			if tableNames[j] > tableNames[j+1] {
				tableNames[j], tableNames[j+1] = tableNames[j+1], tableNames[j]
			}
		}
	}

	// Create rows for the table
	for _, tableName := range tableNames {
		columns := availableTables[tableName]

		// Format columns with line breaks for better readability
		var formattedColumns string
		for i := 0; i < len(columns); i += 5 {
			end := i + 5
			if end > len(columns) {
				end = len(columns)
			}
			if i > 0 {
				formattedColumns += "\n"
			}
			formattedColumns += strings.Join(columns[i:end], ", ")
		}

		tableRows = append(tableRows, []string{tableName, formattedColumns})
	}

	// Display the table
	pretty.Table(headers, tableRows)
}

// Function to validate table name
func isValidTable(tableName string) bool {
	_, exists := availableTables[tableName]
	return exists
}

// Function to validate column names for a specific table
func validateColumns(tableName string, columns []string) []string {
	if tableName == "" || columns == nil || len(columns) == 0 || columns[0] == "*" {
		return nil
	}

	tableColumns, exists := availableTables[tableName]
	if !exists {
		return []string{fmt.Sprintf("Table '%s' does not exist", tableName)}
	}

	var invalidColumns []string
	for _, col := range columns {
		valid := false
		for _, tableCol := range tableColumns {
			if col == tableCol {
				valid = true
				break
			}
		}
		if !valid {
			invalidColumns = append(invalidColumns, col)
		}
	}

	return invalidColumns
}
