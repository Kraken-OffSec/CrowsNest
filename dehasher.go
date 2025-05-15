package main

import (
	"dehasher/cmd"
	"dehasher/internal/badger"
	"dehasher/internal/sqlite"
	"fmt"
	"github.com/winking324/rzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
)

var (
	basePath  string
	logPath   string
	storePath string
	dbPath    string
)

func init() {
	basePath = filepath.Join(os.Getenv("HOME"), ".local", "share", "Dehasher")
	logPath = filepath.Join(basePath, "logs")
	storePath = filepath.Join(basePath, "keystore")
	// dbPath will be set in main() after badger is initialized
}

func createDirectories() {
	var err error

	if _, err = os.Stat(basePath); os.IsNotExist(err) {
		err = os.MkdirAll(basePath, 0755)
		if err != nil {
			zap.L().Error("Error creating directory", zap.Error(err))
			fmt.Printf("[!] Error creating base directory: %v", err)
			os.Exit(-1)
		}
	}

	for _, dir := range []string{"logs", "keystore", "db"} {
		if _, err := os.Stat(filepath.Join(basePath, dir)); os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Join(basePath, dir), 0755)
			if err != nil {
				zap.L().Error("Error creating directory", zap.Error(err), zap.String("directory", dir))
				fmt.Printf("[!] Error creating directory: %v", err)
				os.Exit(-1)
			}
		}
	}
}

func initializeLogger() {
	rzap.NewGlobalLogger([]zapcore.Core{
		rzap.NewCore(&lumberjack.Logger{
			Filename: filepath.Join(logPath, "info.log"),
		}, zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level <= zap.InfoLevel
		})),
		rzap.NewCore(&lumberjack.Logger{
			Filename: filepath.Join(logPath, "error.log"),
		}, zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level > zap.InfoLevel
		})),
	})

	zap.L().Info("logger initialized", zap.Int("status", 0))
}

func main() {
	initializeLogger()

	zap.L().Info("creating_directories")
	createDirectories()

	zap.L().Info("starting_badger")
	db := badger.Start(storePath)
	defer db.Close()

	// Set database path based on useLocalDatabase flag
	useLocalDB := badger.GetUseLocalDB()
	if useLocalDB {
		// Use local database in current directory
		dbPath = "./dehasher.sqlite"
		zap.L().Info("Using local database", zap.String("path", dbPath))
	} else {
		// Use default database path
		dbPath = filepath.Join(basePath, "db")
		zap.L().Info("Using default database path", zap.String("path", dbPath))
	}

	zap.L().Info("initializing_database")
	_, err := sqlite.InitDB(dbPath)
	if err != nil {
		zap.L().Error("init_db",
			zap.String("message", "failed to initialize database"),
			zap.Error(err),
		)
		fmt.Printf("[!] Error initializing database: %v", err)
		os.Exit(1)
	}

	zap.L().Info("executing_command")
	cmd.Execute()
}
