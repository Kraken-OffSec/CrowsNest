package sqlite

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB(dbPath string) (*gorm.DB, error) {
	zap.L().Info("Initializing database", zap.String("path", dbPath))

	// Check if the path is a file or directory
	fileInfo, err := os.Stat(dbPath)
	var finalDbPath string

	// If path doesn't exist or is a directory
	if os.IsNotExist(err) || (err == nil && fileInfo.IsDir()) {
		// Treat as directory path
		if err := os.MkdirAll(dbPath, 0755); err != nil {
			zap.L().Error("Failed to create database directory", zap.Error(err))
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
		finalDbPath = filepath.Join(dbPath, "dehashed.sqlite")
	} else {
		// Treat as file path
		// Ensure the directory exists
		dir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			zap.L().Error("Failed to create parent directory for database", zap.Error(err))
			return nil, fmt.Errorf("failed to create parent directory for database: %w", err)
		}
		finalDbPath = dbPath
	}

	zap.L().Info("Opening database", zap.String("finalPath", finalDbPath))
	db, err := gorm.Open(sqlite.Open(finalDbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		zap.L().Error("Failed to connect to database", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate your models
	err = db.AutoMigrate(&Result{}, &Creds{}, &QueryOptions{}, &Creds{}, &WhoisRecord{}, &SubdomainRecord{}, &HistoryRecord{}, &LookupResult{})
	if err != nil {
		zap.L().Error("Failed to migrate database", zap.Error(err))
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	DB = db
	return db, nil
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	if DB == nil {
		zap.L().Error("database not initialized")
		fmt.Println("sqlite database not initialized")
		os.Exit(1)
	}
	return DB
}

func StoreResults(results DehashedResults) error {
	if len(results.Results) == 0 {
		return nil
	}

	zap.L().Info("Storing results", zap.Int("count", len(results.Results)))
	db := GetDB()

	// Use batch insert with conflict handling
	const batchSize = 100
	var lastErr error

	// Extract the slice of results
	resultSlice := results.Results

	for i := 0; i < len(resultSlice); i += batchSize {
		end := i + batchSize
		if end > len(resultSlice) {
			end = len(resultSlice)
		}

		batch := resultSlice[i:end]
		// Use Clauses with OnConflict DoNothing to skip conflicts
		err := db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&batch, batchSize).Error
		if err != nil {
			zap.L().Warn("Error storing some results", zap.Error(err))
			lastErr = err
			// Continue with next batch despite error
		}
	}

	return lastErr
}

func StoreCreds(creds []Creds) error {
	if len(creds) == 0 {
		return nil
	}

	zap.L().Info("Storing credentials", zap.Int("count", len(creds)))
	db := GetDB()

	// Use batch insert with conflict handling
	// This will insert records in batches and continue even if some fail
	const batchSize = 100
	var lastErr error

	for i := 0; i < len(creds); i += batchSize {
		end := i + batchSize
		if end > len(creds) {
			end = len(creds)
		}

		batch := creds[i:end]
		// Use Clauses with OnConflict DoNothing to skip conflicts
		err := db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&batch, batchSize).Error
		if err != nil {
			zap.L().Warn("Error storing some credentials", zap.Error(err))
			lastErr = err
			// Continue with next batch despite error
		}
	}

	return lastErr
}

func StoreQueryOptions(queryOptions *QueryOptions) error {
	db := GetDB()
	return db.Create(queryOptions).Error
}

func StoreWhoisRecord(whoisRecord WhoisRecord) error {
	// Create a pointer to the record to make it addressable
	recordPtr := &whoisRecord

	zap.L().Info("Storing WHOIS record",
		zap.String("domain", whoisRecord.DomainName))

	db := GetDB()

	// Use OnConflict clause to handle duplicates
	err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(recordPtr).Error
	if err != nil {
		zap.L().Error("store_whois_record",
			zap.String("message", "failed to store whois record"),
			zap.Error(err))
		return err
	}

	return nil
}

func StoreSubdomainRecords(subdomainRecords []SubdomainRecord) error {
	if len(subdomainRecords) == 0 {
		return nil
	}

	zap.L().Info("Storing subdomain records", zap.Int("count", len(subdomainRecords)))
	db := GetDB()

	// Use batch insert with conflict handling
	const batchSize = 100
	var lastErr error

	for i := 0; i < len(subdomainRecords); i += batchSize {
		end := i + batchSize
		if end > len(subdomainRecords) {
			end = len(subdomainRecords)
		}

		batch := subdomainRecords[i:end]
		// Use Clauses with OnConflict DoNothing to skip conflicts
		err := db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&batch, batchSize).Error
		if err != nil {
			zap.L().Warn("Error storing some subdomain records", zap.Error(err))
			lastErr = err
			// Continue with next batch despite error
		}
	}

	return lastErr
}

func StoreHistoryRecord(historyRecords []HistoryRecord) error {
	if len(historyRecords) == 0 {
		return nil
	}

	zap.L().Info("Storing history records", zap.Int("count", len(historyRecords)))
	db := GetDB()

	// Use batch insert with conflict handling
	const batchSize = 100
	var lastErr error

	for i := 0; i < len(historyRecords); i += batchSize {
		end := i + batchSize
		if end > len(historyRecords) {
			end = len(historyRecords)
		}

		batch := historyRecords[i:end]
		// Use Clauses with OnConflict DoNothing to skip conflicts
		err := db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&batch, batchSize).Error
		if err != nil {
			zap.L().Warn("Error storing some history records", zap.Error(err))
			lastErr = err
			// Continue with next batch despite error
		}
	}

	return lastErr
}

func StoreIPLookup(ipLookup []LookupResult) error {
	if len(ipLookup) == 0 {
		return nil
	}

	zap.L().Info("Storing IP lookup records", zap.Int("count", len(ipLookup)))
	db := GetDB()

	// Use batch insert with conflict handling
	const batchSize = 100
	var lastErr error

	for i := 0; i < len(ipLookup); i += batchSize {
		end := i + batchSize
		if end > len(ipLookup) {
			end = len(ipLookup)
		}

		batch := ipLookup[i:end]
		// Use Clauses with OnConflict DoNothing to skip conflicts
		err := db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&batch, batchSize).Error
		if err != nil {
			zap.L().Warn("Error storing some IP lookup records", zap.Error(err))
			lastErr = err
			// Continue with next batch despite error
		}
	}

	return lastErr
}
