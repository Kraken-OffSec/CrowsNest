package sqlite

import (
	"crowsnest/internal/files"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type QueryOptions struct {
	gorm.Model
	MaxRecords         int            `json:"max_records"`
	MaxRequests        int            `json:"max_requests"`
	StartingPage       int            `json:"starting_page"`
	OutputFormat       files.FileType `json:"output_format"`
	OutputFile         string         `json:"output_file"`
	RegexMatch         bool           `json:"regex_match"`
	WildcardMatch      bool           `json:"wildcard_match"`
	UsernameQuery      string         `json:"username_query"`
	EmailQuery         string         `json:"email_query"`
	IpQuery            string         `json:"ip_query"`
	PassQuery          string         `json:"pass_query"`
	HashQuery          string         `json:"hash_query"`
	NameQuery          string         `json:"name_query"`
	DomainQuery        string         `json:"domain_query"`
	VinQuery           string         `json:"vin_query"`
	LicensePlateQuery  string         `json:"license_plate_query"`
	AddressQuery       string         `json:"address_query"`
	PhoneQuery         string         `json:"phone_query"`
	SocialQuery        string         `json:"social_query"`
	CryptoAddressQuery string         `json:"crypto_address_query"`
	PrintBalance       bool           `json:"print_balance"`
	CredsOnly          bool           `json:"creds_only"`
	Debug              bool           `json:"debug"`
}

func (QueryOptions) TableName() string {
	return "query_options"
}

func NewQueryOptions(maxRecords, maxRequests, startingPage int, outputFormat, outputFile, usernameQuery, emailQuery, ipQuery, passQuery, hashQuery, nameQuery, domainQuery, vinQuery, licensePlateQuery, addressQuery, phoneQuery, socialQuery, cryptoAddressQuery string, regexMatch, wildcardMatch, printBalance, credsOnly, debug bool) *QueryOptions {
	return &QueryOptions{
		MaxRecords:         maxRecords,
		MaxRequests:        maxRequests,
		StartingPage:       startingPage,
		OutputFormat:       files.GetFileType(outputFormat),
		OutputFile:         outputFile,
		PrintBalance:       printBalance,
		CredsOnly:          credsOnly,
		RegexMatch:         regexMatch,
		WildcardMatch:      wildcardMatch,
		UsernameQuery:      usernameQuery,
		EmailQuery:         emailQuery,
		IpQuery:            ipQuery,
		PassQuery:          passQuery,
		HashQuery:          hashQuery,
		NameQuery:          nameQuery,
		DomainQuery:        domainQuery,
		VinQuery:           vinQuery,
		LicensePlateQuery:  licensePlateQuery,
		AddressQuery:       addressQuery,
		PhoneQuery:         phoneQuery,
		SocialQuery:        socialQuery,
		CryptoAddressQuery: cryptoAddressQuery,
		Debug:              debug,
	}
}

type DehashedSearchRequest struct {
	Page     int    `json:"page"`
	Query    string `json:"query"`
	Size     int    `json:"size"`
	Wildcard bool   `json:"wildcard"`
	Regex    bool   `json:"regex"`
	DeDupe   bool   `json:"de_dupe"`
}

type DehashedResponse struct {
	Balance      int      `json:"balance"`
	Entries      []Result `json:"entries"`
	Success      bool     `json:"success"`
	Took         string   `json:"took"`
	TotalResults int      `json:"total"`
}

type Result struct {
	gorm.Model
	DehashedId            string   `json:"id" xml:"id" yaml:"id" gorm:"uniqueIndex"`
	Email                 []string `json:"email,omitempty" xml:"email,omitempty" yaml:"email,omitempty" gorm:"serializer:json"`
	IpAddress             []string `json:"ip_address,omitempty" xml:"ip_address,omitempty" yaml:"ip_address,omitempty" gorm:"serializer:json"`
	Username              []string `json:"username,omitempty" xml:"username,omitempty" yaml:"username,omitempty" gorm:"serializer:json"`
	Password              []string `json:"password,omitempty" xml:"password,omitempty" yaml:"password,omitempty" gorm:"serializer:json"`
	HashedPassword        []string `json:"hashed_password,omitempty" xml:"hashed_password,omitempty" yaml:"hashed_password,omitempty" gorm:"serializer:json"`
	HashType              string   `json:"hash_type,omitempty" xml:"hash_type,omitempty" yaml:"hash_type,omitempty"`
	Name                  []string `json:"name,omitempty" xml:"name,omitempty" yaml:"name,omitempty" gorm:"serializer:json"`
	Vin                   []string `json:"vin,omitempty" xml:"vin,omitempty" yaml:"vin,omitempty" gorm:"serializer:json"`
	LicensePlate          []string `json:"license_plate,omitempty" xml:"license_plate,omitempty" yaml:"license_plate,omitempty" gorm:"serializer:json"`
	Url                   []string `json:"url,omitempty" xml:"url,omitempty" yaml:"url,omitempty" gorm:"serializer:json"`
	Social                []string `json:"social,omitempty" xml:"social,omitempty" yaml:"social,omitempty" gorm:"serializer:json"`
	CryptoCurrencyAddress []string `json:"cryptocurrency_address,omitempty" xml:"cryptocurrency_address,omitempty" yaml:"cryptocurrency_address,omitempty" gorm:"serializer:json"`
	Address               []string `json:"address,omitempty" xml:"address,omitempty" yaml:"address,omitempty" gorm:"serializer:json"`
	Phone                 []string `json:"phone,omitempty" xml:"phone,omitempty" yaml:"phone,omitempty" gorm:"serializer:json"`
	Company               []string `json:"company,omitempty" xml:"company,omitempty" yaml:"company,omitempty" gorm:"serializer:json"`
	DatabaseName          string   `json:"database_name,omitempty" xml:"database_name,omitempty" yaml:"database_name,omitempty"`
}

func (Result) TableName() string {
	return "results"
}

type DehashedResults struct {
	Results []Result `json:"results"`
}

func (dr *DehashedResults) ExtractCredentials() []Creds {
	var creds []Creds

	results := dr.Results

	for _, r := range results {
		if len(r.Password) > 0 {
			// Get first email if available
			email := ""
			if len(r.Email) > 0 {
				email = r.Email[0]
			}

			// Get first password
			password := r.Password[0]

			cred := Creds{Email: email, Password: password}
			creds = append(creds, cred)
		}
	}

	go func() {
		err := StoreDehashedCreds(creds)
		if err != nil {
			zap.L().Error("store_creds",
				zap.String("message", "failed to store creds"),
				zap.Error(err),
			)
			fmt.Printf("Error Storing Results: %v", err)
		}
	}()

	return creds
}

type Creds struct {
	gorm.Model
	Email    string `json:"email" yaml:"email" xml:"email" gorm:"uniqueIndex:idx_email_username_password"`
	Username string `json:"username" yaml:"username" xml:"username" gorm:"uniqueIndex:idx_email_username_password"`
	Password string `json:"password" yaml:"password" xml:"password" gorm:"uniqueIndex:idx_email_username_password"`
}

func (Creds) TableName() string {
	return "creds"
}

func (c Creds) ToString() string {
	return fmt.Sprintf("%s%s%s", c.Username, "%", c.Password)
}

func StoreDehashedResults(results DehashedResults) error {
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

func StoreDehashedCreds(creds []Creds) error {
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

func StoreDehashedQueryOptions(queryOptions *QueryOptions) error {
	db := GetDB()
	return db.Create(queryOptions).Error
}
