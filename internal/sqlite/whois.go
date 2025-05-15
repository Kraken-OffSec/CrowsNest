package sqlite

import "gorm.io/gorm"

type WhoIsLookupResult struct {
	RemainingCredits int  `json:"remaining_credits"`
	Data             Data `json:"data"`
}

type Data struct {
	WhoisRecord WhoisRecord `json:"WhoisRecord"`
}

type WhoisRecord struct {
	gorm.Model
	Audit                 Audit        `json:"audit" gorm:"serializer:json"`
	ContactEmail          string       `json:"contactEmail"`
	CreatedDate           string       `json:"createdDate"`
	CreatedDateNormalized string       `json:"createdDateNormalized"`
	DomainName            string       `json:"domainName"`
	DomainNameExt         string       `json:"domainNameExt"`
	EstimatedDomainAge    int          `json:"estimatedDomainAge"`
	ExpiresDate           string       `json:"expiresDate"`
	ExpiresDateNormalized string       `json:"expiresDateNormalized"`
	Footer                string       `json:"footer"`
	Header                string       `json:"header"`
	NameServers           NameServers  `json:"nameServers" gorm:"serializer:json"`
	ParseCode             int          `json:"parseCode"`
	RawText               string       `json:"rawText"`
	Registrant            Contact      `json:"registrant" gorm:"serializer:json"`
	RegistrarIANAID       string       `json:"registrarIANAID"`
	RegistrarName         string       `json:"registrarName"`
	RegistryData          RegistryData `json:"registryData" gorm:"serializer:json"`
	Status                string       `json:"status"`
	StrippedText          string       `json:"strippedText"`
	TechnicalContact      Contact      `json:"technicalContact" gorm:"serializer:json"`
	UpdatedDate           string       `json:"updatedDate"`
	UpdatedDateNormalized string       `json:"updatedDateNormalized"`
}

func (WhoisRecord) TableName() string {
	return "whois"
}

type Audit struct {
	CreatedDate string `json:"createdDate"`
	UpdatedDate string `json:"updatedDate"`
}

type NameServers struct {
	HostNames []string `json:"hostNames"`
	IPs       []string `json:"ips"`
	RawText   string   `json:"rawText"`
}

type Contact struct {
	City         string `json:"city"`
	Country      string `json:"country"`
	CountryCode  string `json:"countryCode"`
	Name         string `json:"name"`
	Organization string `json:"organization"`
	PostalCode   string `json:"postalCode"`
	RawText      string `json:"rawText"`
	State        string `json:"state"`
	Street1      string `json:"street1"`
	Telephone    string `json:"telephone"`
}

type RegistryData struct {
	Audit                 Audit       `json:"audit"`
	CreatedDate           string      `json:"createdDate"`
	CreatedDateNormalized string      `json:"createdDateNormalized"`
	DomainName            string      `json:"domainName"`
	ExpiresDate           string      `json:"expiresDate"`
	ExpiresDateNormalized string      `json:"expiresDateNormalized"`
	Footer                string      `json:"footer"`
	Header                string      `json:"header"`
	NameServers           NameServers `json:"nameServers"`
	ParseCode             int         `json:"parseCode"`
	RawText               string      `json:"rawText"`
	RegistrarIANAID       string      `json:"registrarIANAID"`
	RegistrarName         string      `json:"registrarName"`
	Status                string      `json:"status"`
	StrippedText          string      `json:"strippedText"`
	UpdatedDate           string      `json:"updatedDate"`
	UpdatedDateNormalized string      `json:"updatedDateNormalized"`
	WhoisServer           string      `json:"whoisServer"`
}

type WhoIsSubdomainScan struct {
	RemainingCredits int      `json:"remaining_credits"`
	Data             ScanData `json:"data"`
}

type ScanData struct {
	Result ScanResult `json:"result"`
	Search string     `json:"search"`
}

type ScanResult struct {
	Count   int               `json:"count"`
	Records []SubdomainRecord `json:"records"`
}

type SubdomainRecord struct {
	gorm.Model
	Domain    string `json:"domain"`
	FirstSeen int64  `json:"firstSeen"`
	LastSeen  int64  `json:"lastSeen"`
}

func (SubdomainRecord) TableName() string {
	return "subdomains"
}

type WhoIsHistory struct {
	RemainingCredits int         `json:"remaining_credits"`
	Data             HistoryData `json:"data"`
}

type HistoryData struct {
	Records      []HistoryRecord `json:"records"`
	RecordsCount int             `json:"recordsCount"`
}

type HistoryRecord struct {
	gorm.Model
	AdministrativeContact ContactInfo `json:"administrativeContact" gorm:"serializer:json"`
	Audit                 Audit       `json:"audit" gorm:"serializer:json"`
	BillingContact        ContactInfo `json:"billingContact" gorm:"serializer:json"`
	CleanText             string      `json:"cleanText"`
	CreatedDateISO8601    string      `json:"createdDateISO8601"`
	CreatedDateRaw        string      `json:"createdDateRaw"`
	DomainName            string      `json:"domainName"`
	DomainType            string      `json:"domainType"`
	ExpiresDateISO8601    string      `json:"expiresDateISO8601"`
	ExpiresDateRaw        string      `json:"expiresDateRaw"`
	NameServers           []string    `json:"nameServers" gorm:"serializer:json"`
	RawText               string      `json:"rawText"`
	RegistrantContact     ContactInfo `json:"registrantContact" gorm:"serializer:json"`
	RegistrarName         string      `json:"registrarName"`
	Status                []string    `json:"status" gorm:"serializer:json"`
	TechnicalContact      ContactInfo `json:"technicalContact" gorm:"serializer:json"`
	UpdatedDateISO8601    string      `json:"updatedDateISO8601"`
	UpdatedDateRaw        string      `json:"updatedDateRaw"`
	WhoisServer           string      `json:"whoisServer"`
	ZoneContact           ContactInfo `json:"zoneContact" gorm:"serializer:json"`
}

func (HistoryRecord) TableName() string {
	return "history"
}

type ContactInfo struct {
	City         string `json:"city"`
	Country      string `json:"country"`
	Email        string `json:"email"`
	Fax          string `json:"fax"`
	FaxExt       string `json:"faxExt"`
	Name         string `json:"name"`
	Organization string `json:"organization"`
	PostalCode   string `json:"postalCode"`
	RawText      string `json:"rawText"`
	State        string `json:"state"`
	Street       string `json:"street"`
	Telephone    string `json:"telephone"`
	TelephoneExt string `json:"telephoneExt"`
}

type WhoIsCredits struct {
	WhoisCredits int `json:"whois_credits"`
}
