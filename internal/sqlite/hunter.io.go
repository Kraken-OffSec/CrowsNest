package sqlite

import (
	"fmt"
	"github.com/charmbracelet/lipgloss/tree"
	"gorm.io/gorm"
)

// HunterDomainSearchResult represents the response from Hunter.io domain search API
type HunterDomainSearchResult struct {
	Data HunterDomainData `json:"data" gorm:"embedded;embeddedPrefix:data_"`
	Meta HunterMeta       `json:"meta" gorm:"embedded;embeddedPrefix:meta_"`
}

// HunterDomainData contains the main domain information
type HunterDomainData struct {
	IString
	gorm.Model
	Domain        string        `json:"domain" gorm:"unique"`
	Disposable    bool          `json:"disposable"`
	Webmail       bool          `json:"webmail"`
	AcceptAll     bool          `json:"accept_all"`
	Pattern       string        `json:"pattern"`
	Organization  string        `json:"organization"`
	Description   string        `json:"description"`
	Industry      string        `json:"industry"`
	Twitter       string        `json:"twitter"`
	Facebook      string        `json:"facebook"`
	Linkedin      string        `json:"linkedin"`
	Instagram     string        `json:"instagram"`
	Youtube       string        `json:"youtube"`
	Technologies  []string      `json:"technologies" gorm:"serializer:json"`
	Country       string        `json:"country"`
	State         string        `json:"state"`
	City          string        `json:"city"`
	PostalCode    string        `json:"postal_code"`
	Street        string        `json:"street"`
	Headcount     string        `json:"headcount"`
	CompanyType   string        `json:"company_type"`
	Emails        []HunterEmail `json:"emails" gorm:"serializer:json"`
	LinkedDomains []string      `json:"linked_domains" gorm:"serializer:json"`
}

func (h HunterDomainData) String() string {
	return fmt.Sprintf("Domain: %s\nDisposable: %t\nWebmail: %t\nAcceptAll: %t\nPattern: %s\nOrganization: %s\nDescription: %s\nIndustry: %s\nTwitter: %s\nFacebook: %s\nLinkedin: %s\nInstagram: %s\nYoutube: %s\nTechnologies: %v\nCountry: %s\nState: %s\nCity: %s\nPostalCode: %s\nStreet: %s\nHeadcount: %s\nCompanyType: %s\nEmails: %v\nLinkedDomains: %v\n",
		h.Domain, h.Disposable, h.Webmail, h.AcceptAll, h.Pattern, h.Organization, h.Description, h.Industry, h.Twitter, h.Facebook, h.Linkedin, h.Instagram, h.Youtube, h.Technologies, h.Country, h.State, h.City, h.PostalCode, h.Street, h.Headcount, h.CompanyType, h.Emails, h.LinkedDomains)
}

func (HunterDomainData) TableName() string {
	return "hunter_domain"
}

// HunterEmail represents an email found for the domain
type HunterEmail struct {
	gorm.Model
	Domain       string             `json:"domain,omitempty"`
	Value        string             `json:"value" gorm:"unique"`
	Type         string             `json:"type"`
	Confidence   int                `json:"confidence"`
	Sources      []HunterSource     `json:"sources" gorm:"serializer:json"`
	FirstName    string             `json:"first_name"`
	LastName     string             `json:"last_name"`
	Position     string             `json:"position"`
	PositionRaw  string             `json:"position_raw"`
	Seniority    string             `json:"seniority"`
	Department   string             `json:"department"`
	Linkedin     string             `json:"linkedin"`
	Twitter      string             `json:"twitter"`
	PhoneNumber  string             `json:"phone_number"`
	Verification HunterVerification `json:"verification" gorm:"embedded;embeddedPrefix:verification_"`
}

func (he *HunterEmail) ToTree() *tree.Tree {
	emailTree := tree.Root(he.Value)
	emailTree.Child("Type: " + he.Type)
	emailTree.Child("Confidence: " + fmt.Sprintf("%d", he.Confidence))
	emailTree.Child("FirstName: " + he.FirstName)
	emailTree.Child("LastName: " + he.LastName)
	emailTree.Child("Position: " + he.Position)
	emailTree.Child("PositionRaw: " + he.PositionRaw)
	emailTree.Child("Seniority: " + he.Seniority)
	emailTree.Child("Department: " + he.Department)
	emailTree.Child("Linkedin: " + he.Linkedin)
	emailTree.Child("Twitter: " + he.Twitter)
	emailTree.Child("PhoneNumber: " + he.PhoneNumber)
	emailTree.Child(he.Verification.ToTree())
	return emailTree
}

func (he *HunterEmail) String() string {
	return fmt.Sprintf("Value: %s\nType: %s\nConfidence: %d\nSources: %v\nFirstName: %s\nLastName: %s\nPosition: %s\nPositionRaw: %s\nSeniority: %s\nDepartment: %s\nLinkedin: %s\nTwitter: %s\nPhoneNumber: %s\nVerification: %v\n",
		he.Value, he.Type, he.Confidence, he.Sources, he.FirstName, he.LastName, he.Position, he.PositionRaw, he.Seniority, he.Department, he.Linkedin, he.Twitter, he.PhoneNumber, he.Verification)
}

func (HunterEmail) TableName() string {
	return "hunter_email"
}

// HunterSource represents where an email was found
type HunterSource struct {
	Domain      string `json:"domain"`
	URI         string `json:"uri"`
	ExtractedOn string `json:"extracted_on"`
	LastSeenOn  string `json:"last_seen_on"`
	StillOnPage bool   `json:"still_on_page"`
}

// HunterVerification represents the verification status of an email
type HunterVerification struct {
	Date   string `json:"date"`
	Status string `json:"status"`
}

func (hv *HunterVerification) ToTree() *tree.Tree {
	verificationTree := tree.Root("Verification")
	verificationTree.Child("Date: " + hv.Date)
	verificationTree.Child("Status: " + hv.Status)
	return verificationTree
}

// HunterMeta contains metadata about the API response
type HunterMeta struct {
	Results int                `json:"results"`
	Limit   int                `json:"limit"`
	Offset  int                `json:"offset"`
	Params  HunterSearchParams `json:"params" gorm:"embedded;embeddedPrefix:params_"`
}

// HunterSearchParams contains the parameters used in the search
type HunterSearchParams struct {
	Domain     string `json:"domain"`
	Company    string `json:"company"`
	Type       string `json:"type"`
	Seniority  string `json:"seniority"`
	Department string `json:"department"`
}

// HunterEmailFinderResponse represents the response from Hunter.io email finder API
type HunterEmailFinderResponse struct {
	Data HunterEmailFinderData `json:"data" gorm:"embedded;embeddedPrefix:data_"`
	Meta EmailFinderMeta       `json:"meta" gorm:"embedded;embeddedPrefix:meta_"`
}

// HunterEmailFinderData contains the main email information
type HunterEmailFinderData struct {
	IString
	FirstName    string             `json:"first_name"`
	LastName     string             `json:"last_name"`
	Email        string             `json:"email"`
	Score        int                `json:"score"`
	Domain       string             `json:"domain"`
	AcceptAll    bool               `json:"accept_all"`
	Position     string             `json:"position"`
	Twitter      string             `json:"twitter"`
	LinkedinURL  string             `json:"linkedin_url"`
	PhoneNumber  string             `json:"phone_number"`
	Company      string             `json:"company"`
	Sources      []HunterSource     `json:"sources" gorm:"serializer:json"`
	Verification HunterVerification `json:"verification" gorm:"embedded;embeddedPrefix:verification_"`
}

func (he HunterEmailFinderData) String() string {
	return fmt.Sprintf("FirstName: %s\nLastName: %s\nEmail: %s\nScore: %d\nDomain: %s\nAcceptAll: %t\nPosition: %s\nTwitter: %s\nLinkedinURL: %s\nPhoneNumber: %s\nCompany: %s\nSources: %v\nVerification: %v\n",
		he.FirstName, he.LastName, he.Email, he.Score, he.Domain, he.AcceptAll, he.Position, he.Twitter, he.LinkedinURL, he.PhoneNumber, he.Company, he.Sources, he.Verification)
}

// EmailFinderMeta contains metadata about the API response
type EmailFinderMeta struct {
	Params EmailFinderParams `json:"params" gorm:"embedded;embeddedPrefix:params_"`
}

// EmailFinderParams contains the parameters used in the search
type EmailFinderParams struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	FullName    string `json:"full_name"`
	Domain      string `json:"domain"`
	Company     string `json:"company"`
	MaxDuration string `json:"max_duration"`
}

func (HunterEmailFinderResponse) TableName() string {
	return "hunter_email_finder"
}

// HunterEmailVerifyResponse represents the response from Hunter.io email verification API
type HunterEmailVerifyResponse struct {
	Data HunterEmailVerifyData `json:"data" gorm:"embedded;embeddedPrefix:data_"`
	Meta EmailVerifyMeta       `json:"meta" gorm:"embedded;embeddedPrefix:meta_"`
}

// HunterEmailVerifyData contains the email verification information
type HunterEmailVerifyData struct {
	IString
	Status            string         `json:"status"`
	Result            string         `json:"result"`
	DeprecationNotice string         `json:"_deprecation_notice"`
	Score             int            `json:"score"`
	Email             string         `json:"email"`
	Regexp            bool           `json:"regexp"`
	Gibberish         bool           `json:"gibberish"`
	Disposable        bool           `json:"disposable"`
	Webmail           bool           `json:"webmail"`
	MXRecords         bool           `json:"mx_records"`
	SMTPServer        bool           `json:"smtp_server"`
	SMTPCheck         bool           `json:"smtp_check"`
	AcceptAll         bool           `json:"accept_all"`
	Block             bool           `json:"block"`
	Sources           []HunterSource `json:"sources" gorm:"serializer:json"`
}

func (ev HunterEmailVerifyData) String() string {
	return fmt.Sprintf("Status: %s\nResult: %s\nDeprecationNotice: %s\nScore: %d\nEmail: %s\nRegexp: %t\nGibberish: %t\nDisposable: %t\nWebmail: %t\nMXRecords: %t\nSMTPServer: %t\nSMTPCheck: %t\nAcceptAll: %t\nBlock: %t\nSources: %v\n",
		ev.Status, ev.Result, ev.DeprecationNotice, ev.Score, ev.Email, ev.Regexp, ev.Gibberish, ev.Disposable, ev.Webmail, ev.MXRecords, ev.SMTPServer, ev.SMTPCheck, ev.AcceptAll, ev.Block, ev.Sources)
}

// EmailVerifyMeta contains metadata about the API response
type EmailVerifyMeta struct {
	Params EmailVerifyParams `json:"params" gorm:"embedded;embeddedPrefix:params_"`
}

// EmailVerifyParams contains the parameters used in the verification
type EmailVerifyParams struct {
	Email string `json:"email"`
}

// HunterCompanyEnrichmentResponse represents the response from Hunter.io company enrichment API
type HunterCompanyEnrichmentResponse struct {
	Data CompanyData `json:"data" gorm:"embedded;embeddedPrefix:data_"`
	Meta CompanyMeta `json:"meta" gorm:"embedded;embeddedPrefix:meta_"`
}

// CompanyData contains the detailed company information
type CompanyData struct {
	IString
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	LegalName      string        `json:"legalName"`
	Domain         string        `json:"domain"`
	DomainAliases  []string      `json:"domainAliases" gorm:"serializer:json"`
	Site           CompanySite   `json:"site" gorm:"embedded;embeddedPrefix:site_"`
	Category       Category      `json:"category" gorm:"embedded;embeddedPrefix:category_"`
	Tags           []string      `json:"tags" gorm:"serializer:json"`
	Description    string        `json:"description"`
	FoundedYear    int           `json:"foundedYear"`
	Location       string        `json:"location"`
	TimeZone       string        `json:"timeZone"`
	UTCOffset      int           `json:"utcOffset"`
	Geo            Geography     `json:"geo" gorm:"embedded;embeddedPrefix:geo_"`
	Logo           string        `json:"logo"`
	Facebook       Facebook      `json:"facebook" gorm:"embedded;embeddedPrefix:facebook_"`
	LinkedIn       LinkedIn      `json:"linkedin" gorm:"embedded;embeddedPrefix:linkedin_"`
	Twitter        Twitter       `json:"twitter" gorm:"embedded;embeddedPrefix:twitter_"`
	Crunchbase     Crunchbase    `json:"crunchbase" gorm:"embedded;embeddedPrefix:crunchbase_"`
	YouTube        YouTube       `json:"youtube" gorm:"embedded;embeddedPrefix:youtube_"`
	EmailProvider  string        `json:"emailProvider"`
	Type           string        `json:"type"`
	Ticker         string        `json:"ticker"`
	Identifiers    Identifiers   `json:"identifiers" gorm:"embedded;embeddedPrefix:identifiers_"`
	Phone          string        `json:"phone"`
	Metrics        Metrics       `json:"metrics" gorm:"embedded;embeddedPrefix:metrics_"`
	IndexedAt      string        `json:"indexedAt"`
	Tech           []string      `json:"tech" gorm:"serializer:json"`
	TechCategories []string      `json:"techCategories" gorm:"serializer:json"`
	Parent         ParentCompany `json:"parent" gorm:"embedded;embeddedPrefix:parent_"`
	UltimateParent ParentCompany `json:"ultimateParent" gorm:"embedded;embeddedPrefix:ultimate_parent_"`
}

func (cd CompanyData) String() string {
	return fmt.Sprintf("ID: %s\nName: %s\nLegalName: %s\nDomain: %s\nDomainAliases: %v\nSite: %v\nCategory: %v\nTags: %v\nDescription: %s\nFoundedYear: %d\nLocation: %s\nTimeZone: %s\nUTCOffset: %d\nGeo: %v\nLogo: %s\nFacebook: %v\nLinkedIn: %v\nTwitter: %v\nCrunchbase: %v\nYouTube: %v\nEmailProvider: %s\nType: %s\nTicker: %s\nIdentifiers: %v\nPhone: %s\nMetrics: %v\nIndexedAt: %s\nTech: %v\nTechCategories: %v\nParent: %v\nUltimateParent: %v\n",
		cd.ID, cd.Name, cd.LegalName, cd.Domain, cd.DomainAliases, cd.Site, cd.Category, cd.Tags, cd.Description, cd.FoundedYear, cd.Location, cd.TimeZone, cd.UTCOffset, cd.Geo, cd.Logo, cd.Facebook, cd.LinkedIn, cd.Twitter, cd.Crunchbase, cd.YouTube, cd.EmailProvider, cd.Type, cd.Ticker, cd.Identifiers, cd.Phone, cd.Metrics, cd.IndexedAt, cd.Tech, cd.TechCategories, cd.Parent, cd.UltimateParent)
}

func (cd *CompanyData) DomainAliasesTree() *tree.Tree {
	domainAliasesTree := tree.Root("Domain Aliases")
	for _, domainAlias := range cd.DomainAliases {
		domainAliasesTree.Child(domainAlias)
	}
	return domainAliasesTree
}

func (cd *CompanyData) SiteTree() *tree.Tree {
	siteTree := tree.Root("Site")
	phoneTree := tree.Root("Phone Numbers")
	for _, phoneNumber := range cd.Site.PhoneNumbers {
		phoneTree.Child(phoneNumber)
	}
	emailTree := tree.Root("Email Addresses")
	for _, emailAddress := range cd.Site.EmailAddresses {
		emailTree.Child(emailAddress)
	}
	siteTree.Child(phoneTree)
	siteTree.Child(emailTree)
	return siteTree
}

func (cd *CompanyData) CategoryTree() *tree.Tree {
	categoryTree := tree.Root("Category")
	categoryTree.Child("Sector: " + cd.Category.Sector)
	categoryTree.Child("Industry Group: " + cd.Category.IndustryGroup)
	categoryTree.Child("Industry: " + cd.Category.Industry)
	categoryTree.Child("Sub Industry: " + cd.Category.SubIndustry)
	categoryTree.Child("GICS Code: " + cd.Category.GICSCode)
	categoryTree.Child("SIC Code: " + cd.Category.SICCode)

	sic4CodesTree := tree.Root("SIC 4 Codes")
	for _, sic4Code := range cd.Category.SIC4Codes {
		sic4CodesTree.Child(sic4Code)
	}
	categoryTree.Child(sic4CodesTree)

	categoryTree.Child("NAICS Code: " + cd.Category.NAICSCode)

	naics6CodesTree := tree.Root("NAICS 6 Codes")
	for _, naics6Code := range cd.Category.NAICS6Codes {
		naics6CodesTree.Child(naics6Code)
	}
	categoryTree.Child(naics6CodesTree)

	naics6Codes2022Tree := tree.Root("NAICS 6 Codes 2022")
	for _, naics6Code2022 := range cd.Category.NAICS6Codes2022 {
		naics6Codes2022Tree.Child(naics6Code2022)
	}
	categoryTree.Child(naics6Codes2022Tree)
	return categoryTree
}

func (cd *CompanyData) GeoTree() *tree.Tree {
	geoTree := tree.Root("Geo")
	geoTree.Child("Street Number: " + cd.Geo.StreetNumber)
	geoTree.Child("Street Name: " + cd.Geo.StreetName)
	geoTree.Child("Sub Premise: " + cd.Geo.SubPremise)
	geoTree.Child("Street Address: " + cd.Geo.StreetAddress)
	geoTree.Child("City: " + cd.Geo.City)
	geoTree.Child("Postal Code: " + cd.Geo.PostalCode)
	geoTree.Child("State: " + cd.Geo.State)
	geoTree.Child("State Code: " + cd.Geo.StateCode)
	geoTree.Child("Country: " + cd.Geo.Country)
	geoTree.Child("Country Code: " + cd.Geo.CountryCode)
	geoTree.Child("Latitude: " + fmt.Sprintf("%f", cd.Geo.Lat))
	geoTree.Child("Longitude: " + fmt.Sprintf("%f", cd.Geo.Lng))
	return geoTree
}

func (cd *CompanyData) FacebookTree() *tree.Tree {
	facebookTree := tree.Root("Facebook")
	facebookTree.Child("Handle: " + cd.Facebook.Handle)
	facebookTree.Child("Likes: " + fmt.Sprintf("%d", cd.Facebook.Likes))
	return facebookTree
}

func (cd *CompanyData) LinkedInTree() *tree.Tree {
	linkedinTree := tree.Root("LinkedIn")
	linkedinTree.Child("Handle: " + cd.LinkedIn.Handle)
	return linkedinTree
}

func (cd *CompanyData) TwitterTree() *tree.Tree {
	twitterTree := tree.Root("Twitter")
	twitterTree.Child("Handle: " + cd.Twitter.Handle)
	twitterTree.Child("ID: " + cd.Twitter.ID)
	twitterTree.Child("Bio: " + cd.Twitter.Bio)
	twitterTree.Child("Followers: " + fmt.Sprintf("%d", cd.Twitter.Followers))
	twitterTree.Child("Following: " + fmt.Sprintf("%d", cd.Twitter.Following))
	twitterTree.Child("Location: " + cd.Twitter.Location)
	twitterTree.Child("Site: " + cd.Twitter.Site)
	twitterTree.Child("Avatar" + cd.Twitter.Avatar)
	return twitterTree
}

func (cd *CompanyData) CrunchbaseTree() *tree.Tree {
	crunchbaseTree := tree.Root("Crunchbase")
	crunchbaseTree.Child("Handle: " + cd.Crunchbase.Handle)
	return crunchbaseTree
}

func (cd *CompanyData) YouTubeTree() *tree.Tree {
	youtubeTree := tree.Root("YouTube")
	youtubeTree.Child("Handle: " + cd.YouTube.Handle)
	return youtubeTree
}

func (cd *CompanyData) IdentifiersTree() *tree.Tree {
	identifiersTree := tree.Root("Identifiers")
	identifiersTree.Child("UsEIN: " + cd.Identifiers.UsEIN)
	return identifiersTree
}

func (cd *CompanyData) MetricsTree() *tree.Tree {
	metricsTree := tree.Root("Metrics")
	metricsTree.Child("Alexa Us Rank: " + fmt.Sprintf("%d", cd.Metrics.AlexaUsRank))
	metricsTree.Child("Alexa Global Rank: " + fmt.Sprintf("%d", cd.Metrics.AlexaGlobalRank))
	metricsTree.Child("Traffic Rank: " + cd.Metrics.TrafficRank)
	metricsTree.Child("Employees: " + cd.Metrics.Employees)
	metricsTree.Child("Market Cap: " + cd.Metrics.MarketCap)
	metricsTree.Child("Raised: " + cd.Metrics.Raised)
	metricsTree.Child("Annual Revenue: " + cd.Metrics.AnnualRevenue)
	metricsTree.Child("Estimated Annual Revenue: " + cd.Metrics.EstimatedAnnualRevenue)
	metricsTree.Child("Fiscal Year End: " + cd.Metrics.FiscalYearEnd)
	return metricsTree
}

func (cd *CompanyData) TagsTree() *tree.Tree {
	tagsTree := tree.Root("Tags")
	for _, tag := range cd.Tags {
		tagsTree.Child(tag)
	}
	return tagsTree
}

func (cd *CompanyData) TechTree() *tree.Tree {
	techTree := tree.Root("Tech")
	for _, tech := range cd.Tech {
		techTree.Child(tech)
	}
	return techTree
}

func (cd *CompanyData) TechCategoriesTree() *tree.Tree {
	techCategoriesTree := tree.Root("Tech Categories")
	for _, techCategory := range cd.TechCategories {
		techCategoriesTree.Child(techCategory)
	}
	return techCategoriesTree
}

func (cd *CompanyData) ParentTree() *tree.Tree {
	parentTree := tree.Root("Parent")
	parentTree.Child("Domain: " + cd.Parent.Domain)
	return parentTree
}

func (cd *CompanyData) UltimateParentTree() *tree.Tree {
	ultimateParentTree := tree.Root("Ultimate Parent")
	ultimateParentTree.Child("Domain: " + cd.UltimateParent.Domain)
	return ultimateParentTree
}

// CompanySite contains contact information from the company website
type CompanySite struct {
	PhoneNumbers   []string `json:"phoneNumbers" gorm:"serializer:json"`
	EmailAddresses []string `json:"emailAddresses" gorm:"serializer:json"`
}

// Category contains industry classification information
type Category struct {
	Sector          string   `json:"sector"`
	IndustryGroup   string   `json:"industryGroup"`
	Industry        string   `json:"industry"`
	SubIndustry     string   `json:"subIndustry"`
	GICSCode        string   `json:"gicsCode"`
	SICCode         string   `json:"sicCode"`
	SIC4Codes       []string `json:"sic4Codes" gorm:"serializer:json"`
	NAICSCode       string   `json:"naicsCode"`
	NAICS6Codes     []string `json:"naics6Codes" gorm:"serializer:json"`
	NAICS6Codes2022 []string `json:"naics6Codes2022" gorm:"serializer:json"`
}

// Geography contains location information
type Geography struct {
	StreetNumber  string  `json:"streetNumber"`
	StreetName    string  `json:"streetName"`
	SubPremise    string  `json:"subPremise"`
	StreetAddress string  `json:"streetAddress"`
	City          string  `json:"city"`
	PostalCode    string  `json:"postalCode"`
	State         string  `json:"state"`
	StateCode     string  `json:"stateCode"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"countryCode"`
	Lat           float64 `json:"lat"`
	Lng           float64 `json:"lng"`
}

// Identifiers contains company identification numbers
type Identifiers struct {
	UsEIN string `json:"usEIN"`
}

// Metrics contains company performance metrics
type Metrics struct {
	AlexaUsRank            int    `json:"alexaUsRank"`
	AlexaGlobalRank        int    `json:"alexaGlobalRank"`
	TrafficRank            string `json:"trafficRank"`
	Employees              string `json:"employees"`
	MarketCap              string `json:"marketCap"`
	Raised                 string `json:"raised"`
	AnnualRevenue          string `json:"annualRevenue"`
	EstimatedAnnualRevenue string `json:"estimatedAnnualRevenue"`
	FiscalYearEnd          string `json:"fiscalYearEnd"`
}

// ParentCompany contains information about parent companies
type ParentCompany struct {
	Domain string `json:"domain"`
}

// CompanyMeta contains metadata about the API response
type CompanyMeta struct {
	Domain string `json:"domain"`
}

func (HunterCompanyEnrichmentResponse) TableName() string {
	return "hunter_company_enrichment"
}

// HunterPersonEnrichmentResponse represents the response from Hunter.io person enrichment API
type HunterPersonEnrichmentResponse struct {
	Data PersonData `json:"data" gorm:"embedded;embeddedPrefix:data_"`
	Meta PersonMeta `json:"meta" gorm:"embedded;embeddedPrefix:meta_"`
}

// PersonData contains the detailed person information
type PersonData struct {
	IString
	gorm.Model
	ID            string     `json:"id"`
	Name          PersonName `json:"name" gorm:"embedded;embeddedPrefix:name_"`
	Email         string     `json:"email" gorm:"unique"`
	Location      string     `json:"location"`
	TimeZone      string     `json:"timeZone"`
	UTCOffset     int        `json:"utcOffset"`
	Geo           PersonGeo  `json:"geo" gorm:"embedded;embeddedPrefix:geo_"`
	Bio           string     `json:"bio"`
	Site          string     `json:"site"`
	Avatar        string     `json:"avatar"`
	Employment    Employment `json:"employment" gorm:"embedded;embeddedPrefix:employment_"`
	Facebook      Facebook   `json:"facebook" gorm:"embedded;embeddedPrefix:facebook_"`
	GitHub        GitHub     `json:"github" gorm:"embedded;embeddedPrefix:github_"`
	Twitter       Twitter    `json:"twitter" gorm:"embedded;embeddedPrefix:twitter_"`
	LinkedIn      LinkedIn   `json:"linkedin" gorm:"embedded;embeddedPrefix:linkedin_"`
	GooglePlus    GooglePlus `json:"googleplus" gorm:"embedded;embeddedPrefix:googleplus_"`
	Gravatar      Gravatar   `json:"gravatar" gorm:"embedded;embeddedPrefix:gravatar_"`
	Fuzzy         bool       `json:"fuzzy"`
	EmailProvider string     `json:"emailProvider"`
	IndexedAt     string     `json:"indexedAt"`
	Phone         string     `json:"phone"`
	ActiveAt      string     `json:"activeAt"`
	InactiveAt    string     `json:"inactiveAt"`
}

func (pd PersonData) String() string {
	return fmt.Sprintf("ID: %s\nName: %v\nEmail: %s\nLocation: %s\nTimeZone: %s\nUTCOffset: %d\nGeo: %v\nBio: %s\nSite: %s\nAvatar: %s\nEmployment: %v\nFacebook: %v\nGitHub: %v\nTwitter: %v\nLinkedIn: %v\nGooglePlus: %v\nGravatar: %v\nFuzzy: %t\nEmailProvider: %s\nIndexedAt: %s\nPhone: %s\nActiveAt: %s\nInactiveAt: %s\n",
		pd.ID, pd.Name, pd.Email, pd.Location, pd.TimeZone, pd.UTCOffset, pd.Geo, pd.Bio, pd.Site, pd.Avatar, pd.Employment, pd.Facebook, pd.GitHub, pd.Twitter, pd.LinkedIn, pd.GooglePlus, pd.Gravatar, pd.Fuzzy, pd.EmailProvider, pd.IndexedAt, pd.Phone, pd.ActiveAt, pd.InactiveAt)
}

func (pd *PersonData) NameTree() *tree.Tree {
	nameTree := tree.Root("Name")
	nameTree.Child("Full Name: " + pd.Name.FullName)
	nameTree.Child("Given Name: " + pd.Name.GivenName)
	nameTree.Child("Family Name: " + pd.Name.FamilyName)
	return nameTree
}

func (pd *PersonData) GeoTree() *tree.Tree {
	geoTree := tree.Root("Geo")
	geoTree.Child("City: " + pd.Geo.City)
	geoTree.Child("State: " + pd.Geo.State)
	geoTree.Child("State Code: " + pd.Geo.StateCode)
	geoTree.Child("Country: " + pd.Geo.Country)
	geoTree.Child("Country Code: " + pd.Geo.CountryCode)
	geoTree.Child("Latitude: " + fmt.Sprintf("%f", pd.Geo.Lat))
	geoTree.Child("Longitude: " + fmt.Sprintf("%f", pd.Geo.Lng))
	return geoTree
}

func (pd *PersonData) EmploymentTree() *tree.Tree {
	employmentTree := tree.Root("Employment")
	employmentTree.Child("Domain: " + pd.Employment.Domain)
	employmentTree.Child("Name: " + pd.Employment.Name)
	employmentTree.Child("Title: " + pd.Employment.Title)
	employmentTree.Child("Role: " + pd.Employment.Role)
	employmentTree.Child("Sub Role: " + pd.Employment.SubRole)
	employmentTree.Child("Seniority: " + pd.Employment.Seniority)
	return employmentTree
}

func (pd *PersonData) FacebookTree() *tree.Tree {
	facebookTree := tree.Root("Facebook")
	facebookTree.Child("Handle: " + pd.Facebook.Handle)
	facebookTree.Child("Likes: " + fmt.Sprintf("%d", pd.Facebook.Likes))
	return facebookTree
}

func (pd *PersonData) GitHubTree() *tree.Tree {
	githubTree := tree.Root("GitHub")
	githubTree.Child("Handle: " + pd.GitHub.Handle)
	githubTree.Child("ID: " + pd.GitHub.ID)
	githubTree.Child("Avatar: " + pd.GitHub.Avatar)
	githubTree.Child("Company: " + pd.GitHub.Company)
	githubTree.Child("Blog: " + pd.GitHub.Blog)
	githubTree.Child("Followers: " + fmt.Sprintf("%d", pd.GitHub.Followers))
	githubTree.Child("Following: " + fmt.Sprintf("%d", pd.GitHub.Following))
	return githubTree
}

func (pd *PersonData) TwitterTree() *tree.Tree {
	twitterTree := tree.Root("Twitter")
	twitterTree.Child("Handle: " + pd.Twitter.Handle)
	twitterTree.Child("ID: " + pd.Twitter.ID)
	twitterTree.Child("Bio: " + pd.Twitter.Bio)
	twitterTree.Child("Followers: " + fmt.Sprintf("%d", pd.Twitter.Followers))
	twitterTree.Child("Following: " + fmt.Sprintf("%d", pd.Twitter.Following))
	twitterTree.Child("Location: " + pd.Twitter.Location)
	twitterTree.Child("Site: " + pd.Twitter.Site)
	twitterTree.Child("Avatar: " + pd.Twitter.Avatar)
	return twitterTree
}

func (pd *PersonData) LinkedInTree() *tree.Tree {
	linkedinTree := tree.Root("LinkedIn")
	linkedinTree.Child("Handle: " + pd.LinkedIn.Handle)
	return linkedinTree
}

func (pd *PersonData) GooglePlusTree() *tree.Tree {
	googlePlusTree := tree.Root("GooglePlus")
	googlePlusTree.Child("Handle: " + pd.GooglePlus.Handle)
	return googlePlusTree
}

func (pd *PersonData) GravatarTree() *tree.Tree {
	gravatarTree := tree.Root("Gravatar")
	gravatarTree.Child("Handle: " + pd.Gravatar.Handle)
	gravatarTree.Child("Avatar: " + pd.Gravatar.Avatar)
	return gravatarTree
}

func (PersonData) TableName() string {
	return "person"
}

// PersonName contains the person's name components
type PersonName struct {
	FullName   string `json:"fullName"`
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
}

// PersonGeo contains location information for a person
type PersonGeo struct {
	City        string  `json:"city"`
	State       string  `json:"state"`
	StateCode   string  `json:"stateCode"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
}

// Employment contains employment information
type Employment struct {
	Domain    string `json:"domain"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	Role      string `json:"role"`
	SubRole   string `json:"subRole"`
	Seniority string `json:"seniority"`
}

// GitHub contains GitHub profile information
type GitHub struct {
	gorm.Model
	Handle    string `json:"handle"`
	ID        string `json:"id"`
	Avatar    string `json:"avatar"`
	Company   string `json:"company"`
	Blog      string `json:"blog"`
	Followers int    `json:"followers"`
	Following int    `json:"following"`
	PersonID  int    `json:"person_id,omitempty"`
}

// GooglePlus contains Google+ profile information
type GooglePlus struct {
	Handle   string `json:"handle"`
	PersonID int    `json:"person_id,omitempty"`
}

// Gravatar contains Gravatar profile information
type Gravatar struct {
	gorm.Model
	Handle   string   `json:"handle"`
	URLs     []string `json:"urls" gorm:"serializer:json"`
	Avatar   string   `json:"avatar"`
	Avatars  []string `json:"avatars" gorm:"serializer:json"`
	PersonID int      `json:"person_id,omitempty"`
}

// Facebook contains Facebook profile information
type Facebook struct {
	Handle string `json:"handle"`
	Likes  int    `json:"likes"`
}

// LinkedIn contains LinkedIn profile information
type LinkedIn struct {
	Handle string `json:"handle"`
}

// Twitter contains Twitter profile information
type Twitter struct {
	Handle    string `json:"handle"`
	ID        string `json:"id"`
	Bio       string `json:"bio"`
	Followers int    `json:"followers"`
	Following int    `json:"following"`
	Location  string `json:"location"`
	Site      string `json:"site"`
	Avatar    string `json:"avatar"`
}

// Crunchbase contains Crunchbase profile information
type Crunchbase struct {
	Handle string `json:"handle"`
}

// YouTube contains YouTube profile information
type YouTube struct {
	Handle string `json:"handle"`
}

// PersonMeta contains metadata about the API response
type PersonMeta struct {
	Email string `json:"email"`
}

// HunterCombinedEnrichmentResponse represents the response from Hunter.io combined enrichment API
type HunterCombinedEnrichmentResponse struct {
	Data CombinedData `json:"data" gorm:"embedded;embeddedPrefix:data_"`
	Meta CombinedMeta `json:"meta" gorm:"embedded;embeddedPrefix:meta_"`
}

// CombinedData contains both person and company information
type CombinedData struct {
	IString
	Person  PersonData  `json:"person" gorm:"embedded;embeddedPrefix:person_"`
	Company CompanyData `json:"company" gorm:"embedded;embeddedPrefix:company_"`
}

func (cbd CombinedData) String() string {
	return fmt.Sprintf("Person: %s\nCompany: %s",
		cbd.Person.String(),
		cbd.Company.String())
}

// CombinedMeta contains metadata about the API response
type CombinedMeta struct {
	Email string `json:"email"`
}

// String returns a string representation of the combined enrichment response
func (c *HunterCombinedEnrichmentResponse) String() string {
	return fmt.Sprintf("Person:\n%s\n\nCompany:\n%s",
		c.Data.Person.String(),
		c.Data.Company.String())
}
