package sqlite

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Subdomain struct {
	gorm.Model
	Domain    string `json:"domain" yaml:"domain" xml:"domain"`
	Subdomain string `json:"subdomain" yaml:"subdomain" xml:"subdomain" gorm:"uniqueIndex:idx_subdomain"`
}

func StoreSubdomains(subs []Subdomain) error {
	if len(subs) == 0 {
		return nil
	}

	zap.L().Info("Storing subdomains", zap.Int("count", len(subs)))
	db := GetDB()

	// Use batch insert with conflict handling
	// This will insert records in batches and continue even if some fail
	const batchSize = 100
	var lastErr error

	for i := 0; i < len(subs); i += batchSize {
		end := i + batchSize
		if end > len(subs) {
			end = len(subs)
		}

		batch := subs[i:end]
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
