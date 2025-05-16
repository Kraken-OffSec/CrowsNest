package sqlite

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

// QueryResults queries the database for results based on the provided options
func QueryResults(options *DBOptions) ([]Result, error) {
	db := GetDB()
	var results []Result
	query := db.Model(&Result{})

	// Apply filters based on the provided options
	query = applyFilters(query, options)

	// Apply limit
	if options.Limit > 0 {
		query = query.Limit(options.Limit)
	}

	// Execute the query
	if err := query.Find(&results).Error; err != nil {
		zap.L().Error("query_results",
			zap.String("message", "failed to query results"),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to query results: %w", err)
	}

	return results, nil
}

// applyFilters applies filters to the query based on the provided options
func applyFilters(query *gorm.DB, options *DBOptions) *gorm.DB {
	// Helper function to apply filter based on exact match setting
	applyFilter := func(field, value string) *gorm.DB {
		if value == "" {
			return query
		}

		if options.ExactMatch {
			return query.Where(field+" = ?", value)
		} else {
			return query.Where(field+" LIKE ?", "%"+value+"%")
		}
	}

	// Apply filters for each field if provided
	if options.Email != "" {
		query = applyFilter("email", options.Email)
	}

	if options.Username != "" {
		query = applyFilter("username", options.Username)
	}

	if options.IPAddress != "" {
		query = applyFilter("ip_address", options.IPAddress)
	}

	if options.Password != "" {
		query = applyFilter("password", options.Password)
	}

	if options.HashedPassword != "" {
		query = applyFilter("hashed_password", options.HashedPassword)
	}

	if options.Name != "" {
		query = applyFilter("name", options.Name)
	}

	if options.Vin != "" {
		query = applyFilter("vin", options.Vin)
	}

	if options.LicensePlate != "" {
		query = applyFilter("license_plate", options.LicensePlate)
	}

	if options.Address != "" {
		query = applyFilter("address", options.Address)
	}

	if options.Phone != "" {
		query = applyFilter("phone", options.Phone)
	}

	if options.Social != "" {
		query = applyFilter("social", options.Social)
	}

	if options.CryptoCurrencyAddress != "" {
		query = applyFilter("cryptocurrency_address", options.CryptoCurrencyAddress)
	}

	if options.Domain != "" {
		query = applyFilter("url", options.Domain)
	}

	// Apply non-empty field filters
	for _, field := range options.NonEmptyFields {
		switch field {
		case "username":
			query = query.Where("JSON_ARRAY_LENGTH(username) > 0")
		case "email":
			query = query.Where("JSON_ARRAY_LENGTH(email) > 0")
		case "ip_address", "ipaddress", "ip":
			query = query.Where("JSON_ARRAY_LENGTH(ip_address) > 0")
		case "password":
			query = query.Where("JSON_ARRAY_LENGTH(password) > 0")
		case "hashed_password", "hash":
			query = query.Where("JSON_ARRAY_LENGTH(hashed_password) > 0")
		case "name":
			query = query.Where("JSON_ARRAY_LENGTH(name) > 0")
		case "vin":
			query = query.Where("JSON_ARRAY_LENGTH(vin) > 0")
		case "license_plate", "license":
			query = query.Where("JSON_ARRAY_LENGTH(license_plate) > 0")
		case "address":
			query = query.Where("JSON_ARRAY_LENGTH(address) > 0")
		case "phone":
			query = query.Where("JSON_ARRAY_LENGTH(phone) > 0")
		case "social":
			query = query.Where("JSON_ARRAY_LENGTH(social) > 0")
		case "cryptocurrency_address", "crypto":
			query = query.Where("JSON_ARRAY_LENGTH(cryptocurrency_address) > 0")
		case "url", "domain":
			query = query.Where("JSON_ARRAY_LENGTH(url) > 0")
		}
	}

	return query
}

// GetResultsCount returns the count of results matching the provided options
func GetResultsCount(options *DBOptions) (int64, error) {
	db := GetDB()
	var count int64
	query := db.Model(&Result{})

	// Apply filters based on the provided options
	query = applyFilters(query, options)

	// Count the results
	if err := query.Count(&count).Error; err != nil {
		zap.L().Error("get_results_count",
			zap.String("message", "failed to count results"),
			zap.Error(err),
		)
		return 0, fmt.Errorf("failed to count results: %w", err)
	}

	return count, nil
}

// QueryRuns queries the database for previous query runs (QueryOptions) based on the provided filters
func QueryRuns(limit, lastXRuns int, startDate, endDate time.Time, containsQuery string) ([]QueryOptions, error) {
	db := GetDB()
	var runs []QueryOptions
	query := db.Model(&QueryOptions{})

	// Apply date range filter if provided
	if lastXRuns > 0 {
		query = query.Order("created_at DESC").Limit(lastXRuns)
	} else if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	} else if !startDate.IsZero() {
		query = query.Where("created_at >= ?", startDate)
	} else if !endDate.IsZero() {
		query = query.Where("created_at <= ?", endDate)
	}

	// Apply query filter if provided
	if containsQuery != "" {
		// SearchTerm in all query fields
		query = query.Where(
			"username_query LIKE ? OR "+
				"email_query LIKE ? OR "+
				"ip_query LIKE ? OR "+
				"pass_query LIKE ? OR "+
				"hash_query LIKE ? OR "+
				"name_query LIKE ? OR "+
				"domain_query LIKE ? OR "+
				"vin_query LIKE ? OR "+
				"license_plate_query LIKE ? OR "+
				"address_query LIKE ? OR "+
				"phone_query LIKE ? OR "+
				"social_query LIKE ? OR "+
				"crypto_address_query LIKE ?",
			"%"+containsQuery+"%", "%"+containsQuery+"%", "%"+containsQuery+"%",
			"%"+containsQuery+"%", "%"+containsQuery+"%", "%"+containsQuery+"%",
			"%"+containsQuery+"%", "%"+containsQuery+"%", "%"+containsQuery+"%",
			"%"+containsQuery+"%", "%"+containsQuery+"%", "%"+containsQuery+"%",
			"%"+containsQuery+"%",
		)
	}

	// Apply limit
	if limit > 0 {
		query = query.Limit(limit)
	}

	// Order by most recent first
	query = query.Order("created_at DESC")

	// Execute the query
	if err := query.Find(&runs).Error; err != nil {
		zap.L().Error("query_runs",
			zap.String("message", "failed to query runs"),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to query runs: %w", err)
	}

	return runs, nil
}

// GetRunsCount returns the count of runs matching the provided filters
func GetRunsCount(lastXRuns int, startDate, endDate time.Time, containsQuery string) (int64, error) {
	db := GetDB()
	var count int64
	query := db.Model(&QueryOptions{})

	// Apply date range filter if provided
	if lastXRuns > 0 {
		query = query.Order("created_at DESC").Limit(lastXRuns)
	} else if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	} else if !startDate.IsZero() {
		query = query.Where("created_at >= ?", startDate)
	} else if !endDate.IsZero() {
		query = query.Where("created_at <= ?", endDate)
	}

	// Apply query filter if provided
	if containsQuery != "" {
		// SearchTerm in all query fields
		query = query.Where(
			"username_query LIKE ? OR "+
				"email_query LIKE ? OR "+
				"ip_query LIKE ? OR "+
				"pass_query LIKE ? OR "+
				"hash_query LIKE ? OR "+
				"name_query LIKE ? OR "+
				"domain_query LIKE ? OR "+
				"vin_query LIKE ? OR "+
				"license_plate_query LIKE ? OR "+
				"address_query LIKE ? OR "+
				"phone_query LIKE ? OR "+
				"social_query LIKE ? OR "+
				"crypto_address_query LIKE ?",
			"%"+containsQuery+"%", "%"+containsQuery+"%", "%"+containsQuery+"%",
			"%"+containsQuery+"%", "%"+containsQuery+"%", "%"+containsQuery+"%",
			"%"+containsQuery+"%", "%"+containsQuery+"%", "%"+containsQuery+"%",
			"%"+containsQuery+"%", "%"+containsQuery+"%", "%"+containsQuery+"%",
			"%"+containsQuery+"%",
		)
	}

	// Count the results
	if err := query.Count(&count).Error; err != nil {
		zap.L().Error("get_runs_count",
			zap.String("message", "failed to count runs"),
			zap.Error(err),
		)
		return 0, fmt.Errorf("failed to count runs: %w", err)
	}

	return count, nil
}

// QueryCreds queries the database for credentials based on the provided filters
func QueryCreds(options *DBOptions) ([]Creds, error) {
	db := GetDB()
	var creds []Creds
	query := db.Model(&Creds{})

	// Apply filters based on the provided options
	if options.Username != "" {
		if options.ExactMatch {
			query = query.Where("username = ?", options.Username)
		} else {
			query = query.Where("username LIKE ?", "%"+options.Username+"%")
		}
	}

	if options.Email != "" {
		if options.ExactMatch {
			query = query.Where("email = ?", options.Email)
		} else {
			query = query.Where("email LIKE ?", "%"+options.Email+"%")
		}
	}

	if options.Password != "" {
		if options.ExactMatch {
			query = query.Where("password = ?", options.Password)
		} else {
			query = query.Where("password LIKE ?", "%"+options.Password+"%")
		}
	}

	// Apply limit
	if options.Limit > 0 {
		query = query.Limit(options.Limit)
	}

	// Execute the query
	if err := query.Find(&creds).Error; err != nil {
		zap.L().Error("query_creds",
			zap.String("message", "failed to query credentials"),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to query credentials: %w", err)
	}

	return creds, nil
}

// GetCredsCount returns the count of credentials matching the provided filters
func GetCredsCount(options *DBOptions) (int64, error) {
	db := GetDB()
	var count int64
	query := db.Model(&Creds{})

	// Apply filters based on the provided options
	if options.Username != "" {
		if options.ExactMatch {
			query = query.Where("username = ?", options.Username)
		} else {
			query = query.Where("username LIKE ?", "%"+options.Username+"%")
		}
	}

	if options.Email != "" {
		if options.ExactMatch {
			query = query.Where("email = ?", options.Email)
		} else {
			query = query.Where("email LIKE ?", "%"+options.Email+"%")
		}
	}

	if options.Password != "" {
		if options.ExactMatch {
			query = query.Where("password = ?", options.Password)
		} else {
			query = query.Where("password LIKE ?", "%"+options.Password+"%")
		}
	}

	// Count the results
	if err := query.Count(&count).Error; err != nil {
		zap.L().Error("get_creds_count",
			zap.String("message", "failed to count credentials"),
			zap.Error(err),
		)
		return 0, fmt.Errorf("failed to count credentials: %w", err)
	}

	return count, nil
}

// ExecuteRawQuery executes a raw SQL query and returns the results as a slice of maps
func ExecuteRawQuery(query string) ([]map[string]interface{}, error) {
	db := GetDB()
	rows, err := db.Raw(query).Rows()
	if err != nil {
		zap.L().Error("raw_query",
			zap.String("message", "failed to execute raw query"),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to execute raw query: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		zap.L().Error("raw_query",
			zap.String("message", "failed to get columns from raw query"),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get columns from raw query: %w", err)
	}

	var results []map[string]interface{}

	for rows.Next() {
		// Create a slice of interface{} to hold the values
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		// Scan the result into the pointers
		if err := rows.Scan(pointers...); err != nil {
			zap.L().Error("raw_query",
				zap.String("message", "failed to scan row from raw query"),
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to scan row from raw query: %w", err)
		}

		// Create a map for this row
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			rowMap[col] = val
		}

		results = append(results, rowMap)
	}

	return results, nil
}
