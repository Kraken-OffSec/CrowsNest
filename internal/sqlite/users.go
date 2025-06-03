package sqlite

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	gorm.Model
	Company     string `json:"company" yaml:"company" xml:"company"`
	Position    string `json:"position" yaml:"position" xml:"position"`
	Department  string `json:"department" yaml:"department" xml:"department"`
	PhoneNumber string `json:"phone_number" yaml:"phone_number" xml:"phone_number"`
	FullName    string `json:"full_name" yaml:"full_name" xml:"full_name"`
	Phone       string `json:"phone" yaml:"phone" xml:"phone"`
	Linkedin    string `json:"linkedin" yaml:"linkedin" xml:"linkedin"`
	Twitter     string `json:"twitter" yaml:"twitter" xml:"twitter"`
	Facebook    string `json:"facebook" yaml:"facebook" xml:"facebook"`
	Instagram   string `json:"instagram" yaml:"instagram" xml:"instagram"`
	Youtube     string `json:"youtube" yaml:"youtube" xml:"youtube"`
	Gravatar    string `json:"gravatar" yaml:"gravatar" xml:"gravatar"`
	Email       string `json:"email" yaml:"email" xml:"email" gorm:"uniqueIndex:idx_email_username_password"`
	Username    string `json:"username" yaml:"username" xml:"username" gorm:"uniqueIndex:idx_email_username_password"`
	Password    string `json:"password" yaml:"password" xml:"password" gorm:"uniqueIndex:idx_email_username_password"`
}

func StoreUsers(users []User) error {
	if len(users) == 0 {
		return nil
	}

	zap.L().Info("Storing credentials", zap.Int("count", len(users)))
	db := GetDB()

	// Use batch insert with conflict handling
	// This will insert records in batches and continue even if some fail
	const batchSize = 100
	var lastErr error

	for i := 0; i < len(users); i += batchSize {
		end := i + batchSize
		if end > len(users) {
			end = len(users)
		}

		batch := users[i:end]
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
