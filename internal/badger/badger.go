package badger

import (
	"crypto/sha256"
	"errors"
	"github.com/dgraph-io/badger/v4"
	"go.uber.org/zap"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	encryptionKey []byte // must be 32 bytes
	db            *badger.DB
	rootDir       string
	once          sync.Once
)

func GetHardwareEntropy() []byte {
	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-host"
		log.Printf("Error getting hostname: %v", err)
	}

	// Get username
	currentUser, err := user.Current()
	username := "unknown-user"
	if err == nil && currentUser != nil {
		username = currentUser.Username
	}

	// Get OS and architecture info
	osInfo := runtime.GOOS + "-" + runtime.GOARCH

	// Combine all information for a unique but consistent fingerprint
	fingerprint := strings.Join([]string{
		hostname,
		username,
		osInfo,
		// You could add a static salt here for additional security
		"CrowsNest-static-salt-value",
	}, ":")

	// Hash the fingerprint to get a 32-byte key
	sum := sha256.Sum256([]byte(fingerprint))
	return sum[:]
}

func Start(dirPath string) *badger.DB {
	var err error

	zap.L().Info("Starting Badger DB", zap.String("directory", dirPath))
	zap.L().Info("Badger DB Directory Path", zap.String("directory", dirPath))

	once.Do(func() {
		if !strings.HasSuffix(dirPath, "db") {
			dirPath = filepath.Join(dirPath, "db")
		}
		rootDir = dirPath

		encryptionKey = GetHardwareEntropy()
		if err != nil {
			zap.L().Fatal("get_encryption_key",
				zap.String("message", "failed to get encryption key"),
				zap.Error(err),
			)
		}

		badgerDB := filepath.Join(rootDir, "badger.db")
		opts := badger.DefaultOptions(badgerDB).
			WithEncryptionKey(encryptionKey).
			WithIndexCacheSize(10 << 20). // 10MB
			WithLoggingLevel(badger.ERROR)
		db, err = badger.Open(opts)
		if err != nil {
			zap.L().Fatal("new_badger_db",
				zap.String("message", "failed to open badger database"),
				zap.Error(err),
			)
		}
	})

	return db
}

func Close() {
	err := db.Close()
	if err != nil {
		zap.L().Fatal("new_badger_db",
			zap.String("message", "failed to close badger database"),
			zap.Error(err),
		)
	}
}

func GetDehashedKey() string {
	var apiKey string

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("cfg:api_key"))
		if err != nil {
			return err // could be ErrKeyNotFound
		}
		return item.Value(func(val []byte) error {
			apiKey = string(val)
			return nil
		})
	})

	if err != nil {
		zap.L().Error("get_api_key",
			zap.String("message", "failed to get api_key"),
			zap.Error(err),
		)
	}

	return apiKey
}

func GetHunterKey() string {
	var apiKey string

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("cfg:hunter_api_key"))
		if err != nil {
			return err // could be ErrKeyNotFound
		}
		return item.Value(func(val []byte) error {
			apiKey = string(val)
			return nil
		})
	})

	if err != nil {
		zap.L().Error("get_hunter_api_key",
			zap.String("message", "failed to get hunter_api_key"),
			zap.Error(err),
		)
	}
	return apiKey
}

func GetUseLocalDB() bool {
	var useLocal bool

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("cfg:use_local_db"))
		if err != nil {
			// If key not found, set default to false and return nil
			if errors.Is(err, badger.ErrKeyNotFound) {
				// Store the default value for future use
				err = StoreUseLocalDB(false)
				if err != nil {
					zap.L().Error("store_use_local_db",
						zap.String("message", "failed to store use_local_db"),
						zap.Error(err),
					)
					return err
				}
				return nil
			}
			// Return other errors
			return err
		}
		return item.Value(func(val []byte) error {
			useLocal = val[0] == 1
			return nil
		})
	})

	if err != nil {
		zap.L().Error("get_use_local_db",
			zap.String("message", "failed to get use_local_db"),
			zap.Error(err),
		)
	}

	return useLocal
}

func StoreDehashedKey(apiKey string) error {
	err := db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("cfg:api_key"), []byte(apiKey))
	})
	if err != nil {
		zap.L().Error("set_api_key",
			zap.String("message", "failed to set dehashed api_key"),
			zap.Error(err),
		)
	}
	return err
}

func StoreHunterKey(apiKey string) error {
	err := db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("cfg:hunter_api_key"), []byte(apiKey))
	})
	if err != nil {
		zap.L().Error("set_api_key",
			zap.String("message", "failed to set hunter api_key"),
			zap.Error(err),
		)
	}
	return err
}

func StoreUseLocalDB(useLocal bool) error {
	var local byte
	if useLocal {
		local = 1
	} else {
		local = 0
	}

	err := db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("cfg:use_local_db"), []byte{local})
	})
	if err != nil {
		zap.L().Error("set_use_local_db",
			zap.String("message", "failed to set use_local_db"),
			zap.Error(err),
		)
	}
	return err
}
