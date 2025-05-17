package sqlite

import (
	"crowsnest/internal/files"
	"fmt"
	"gorm.io/gorm"
)

type IString interface {
	String() string
}

type DBOptions struct {
	Username              string
	Email                 string
	IPAddress             string
	Password              string
	HashedPassword        string
	Name                  string
	Vin                   string
	LicensePlate          string
	Address               string
	Phone                 string
	Social                string
	CryptoCurrencyAddress string
	Domain                string
	Limit                 int
	ExactMatch            bool
	NonEmptyFields        []string // Fields that should not be empty
	DisplayFields         []string // Fields to display in output
}

func NewDBOptions() *DBOptions {
	return &DBOptions{
		Limit:          100, // Default limit
		ExactMatch:     false,
		NonEmptyFields: []string{},
		DisplayFields:  []string{},
	}
}

func (o *DBOptions) Empty() bool {
	return o.Username == "" && o.Email == "" && o.IPAddress == "" &&
		o.Password == "" && o.HashedPassword == "" && o.Name == "" &&
		o.Vin == "" && o.LicensePlate == "" && o.Address == "" &&
		o.Phone == "" && o.Social == "" && o.CryptoCurrencyAddress == "" && o.Domain == "" &&
		len(o.NonEmptyFields) == 0
}

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
