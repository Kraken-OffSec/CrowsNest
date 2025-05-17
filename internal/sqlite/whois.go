package sqlite

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

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
	DomainName            string       `json:"domainName" gorm:"unique"`
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

func (w WhoisRecord) String() string {
	var sb strings.Builder

	// Main domain information
	sb.WriteString(fmt.Sprintf("Domain Name: %s\n", w.DomainName))
	sb.WriteString(fmt.Sprintf("Domain Name Ext: %s\n", w.DomainNameExt))
	sb.WriteString(fmt.Sprintf("Registrar Name: %s\n", w.RegistrarName))
	sb.WriteString(fmt.Sprintf("Registrar IANA ID: %s\n", w.RegistrarIANAID))
	sb.WriteString(fmt.Sprintf("Contact HunterEmail: %s\n", w.ContactEmail))
	sb.WriteString(fmt.Sprintf("Estimated Domain Age: %d days\n", w.EstimatedDomainAge))

	// Dates
	sb.WriteString(fmt.Sprintf("Created Date: %s (Normalized: %s)\n", w.CreatedDate, w.CreatedDateNormalized))
	sb.WriteString(fmt.Sprintf("Updated Date: %s (Normalized: %s)\n", w.UpdatedDate, w.UpdatedDateNormalized))
	sb.WriteString(fmt.Sprintf("Expires Date: %s (Normalized: %s)\n", w.ExpiresDate, w.ExpiresDateNormalized))

	// Status
	sb.WriteString(fmt.Sprintf("Status: %s\n", w.Status))

	// Parse code
	sb.WriteString(fmt.Sprintf("Parse Code: %d\n", w.ParseCode))

	// Audit information
	sb.WriteString("\nAudit Information:\n")
	sb.WriteString(fmt.Sprintf("  Created Date: %s\n", w.Audit.CreatedDate))
	sb.WriteString(fmt.Sprintf("  Updated Date: %s\n", w.Audit.UpdatedDate))

	// Name servers
	sb.WriteString("\nName Servers:\n")
	if len(w.NameServers.HostNames) > 0 {
		for i, ns := range w.NameServers.HostNames {
			ip := ""
			if i < len(w.NameServers.IPs) {
				ip = fmt.Sprintf(" (%s)", w.NameServers.IPs[i])
			}
			sb.WriteString(fmt.Sprintf("  %d. %s%s\n", i+1, ns, ip))
		}
	} else {
		sb.WriteString("  None listed\n")
	}

	if w.NameServers.RawText != "" {
		sb.WriteString(fmt.Sprintf("  Raw Text: %s\n", w.NameServers.RawText))
	}

	// Contact information
	sb.WriteString("\nRegistrant Contact:\n")
	formatWhoisContact(&sb, w.Registrant, "  ")

	sb.WriteString("\nTechnical Contact:\n")
	formatWhoisContact(&sb, w.TechnicalContact, "  ")

	// Registry Data
	sb.WriteString("\nRegistry Data:\n")
	if w.RegistryData.DomainName != "" {
		sb.WriteString(fmt.Sprintf("  Domain Name: %s\n", w.RegistryData.DomainName))
		sb.WriteString(fmt.Sprintf("  Registrar Name: %s\n", w.RegistryData.RegistrarName))
		sb.WriteString(fmt.Sprintf("  Registrar IANA ID: %s\n", w.RegistryData.RegistrarIANAID))
		sb.WriteString(fmt.Sprintf("  Whois Server: %s\n", w.RegistryData.WhoisServer))
		sb.WriteString(fmt.Sprintf("  Status: %s\n", w.RegistryData.Status))

		// Registry dates
		sb.WriteString(fmt.Sprintf("  Created Date: %s (Normalized: %s)\n",
			w.RegistryData.CreatedDate, w.RegistryData.CreatedDateNormalized))
		sb.WriteString(fmt.Sprintf("  Updated Date: %s (Normalized: %s)\n",
			w.RegistryData.UpdatedDate, w.RegistryData.UpdatedDateNormalized))
		sb.WriteString(fmt.Sprintf("  Expires Date: %s (Normalized: %s)\n",
			w.RegistryData.ExpiresDate, w.RegistryData.ExpiresDateNormalized))

		// Registry nameservers
		sb.WriteString("  Name Servers:\n")
		if len(w.RegistryData.NameServers.HostNames) > 0 {
			for i, ns := range w.RegistryData.NameServers.HostNames {
				ip := ""
				if i < len(w.RegistryData.NameServers.IPs) {
					ip = fmt.Sprintf(" (%s)", w.RegistryData.NameServers.IPs[i])
				}
				sb.WriteString(fmt.Sprintf("    %d. %s%s\n", i+1, ns, ip))
			}
		} else {
			sb.WriteString("    None listed\n")
		}

		// Registry audit
		sb.WriteString("  Audit Information:\n")
		sb.WriteString(fmt.Sprintf("    Created Date: %s\n", w.RegistryData.Audit.CreatedDate))
		sb.WriteString(fmt.Sprintf("    Updated Date: %s\n", w.RegistryData.Audit.UpdatedDate))
	} else {
		sb.WriteString("  No registry data available\n")
	}

	// Header and footer
	if w.Header != "" {
		headerPreview := w.Header
		if len(headerPreview) > 100 {
			headerPreview = headerPreview[:100] + "... [truncated]"
		}
		sb.WriteString("\nHeader:\n")
		sb.WriteString(headerPreview)
		sb.WriteString("\n")
	}

	if w.Footer != "" {
		footerPreview := w.Footer
		if len(footerPreview) > 100 {
			footerPreview = footerPreview[:100] + "... [truncated]"
		}
		sb.WriteString("\nFooter:\n")
		sb.WriteString(footerPreview)
		sb.WriteString("\n")
	}

	// Raw text (truncated if too long)
	if w.RawText != "" {
		rawTextPreview := w.RawText
		if len(rawTextPreview) > 500 {
			rawTextPreview = rawTextPreview[:500] + "... [truncated]"
		}
		sb.WriteString("\nRaw Text:\n")
		sb.WriteString(rawTextPreview)
		sb.WriteString("\n")
	}

	if w.StrippedText != "" {
		strippedTextPreview := w.StrippedText
		if len(strippedTextPreview) > 500 {
			strippedTextPreview = strippedTextPreview[:500] + "... [truncated]"
		}
		sb.WriteString("\nStripped Text:\n")
		sb.WriteString(strippedTextPreview)
		sb.WriteString("\n")
	}

	return sb.String()
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
	Domain    string `json:"domain" gorm:"unique"`
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
	DomainName            string      `json:"domainName" gorm:"unique"`
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

func (h HistoryRecord) String() string {
	var sb strings.Builder

	// Main domain information
	sb.WriteString(fmt.Sprintf("Domain Name: %s\n", h.DomainName))
	sb.WriteString(fmt.Sprintf("Domain Type: %s\n", h.DomainType))
	sb.WriteString(fmt.Sprintf("Registrar Name: %s\n", h.RegistrarName))
	sb.WriteString(fmt.Sprintf("Whois Server: %s\n", h.WhoisServer))

	// Dates
	sb.WriteString(fmt.Sprintf("Created Date: %s (Raw: %s)\n", h.CreatedDateISO8601, h.CreatedDateRaw))
	sb.WriteString(fmt.Sprintf("Updated Date: %s (Raw: %s)\n", h.UpdatedDateISO8601, h.UpdatedDateRaw))
	sb.WriteString(fmt.Sprintf("Expires Date: %s (Raw: %s)\n", h.ExpiresDateISO8601, h.ExpiresDateRaw))

	// Status
	sb.WriteString("Status: ")
	if len(h.Status) > 0 {
		sb.WriteString(strings.Join(h.Status, ", "))
	} else {
		sb.WriteString("N/A")
	}
	sb.WriteString("\n")

	// Audit information
	sb.WriteString("\nAudit Information:\n")
	sb.WriteString(fmt.Sprintf("  Created Date: %s\n", h.Audit.CreatedDate))
	sb.WriteString(fmt.Sprintf("  Updated Date: %s\n", h.Audit.UpdatedDate))

	// Name servers
	sb.WriteString("\nName Servers:\n")
	if len(h.NameServers) > 0 {
		for i, ns := range h.NameServers {
			sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, ns))
		}
	} else {
		sb.WriteString("  None listed\n")
	}

	// Contact information
	sb.WriteString("\nRegistrant Contact:\n")
	formatContact(&sb, h.RegistrantContact, "  ")

	sb.WriteString("\nAdministrative Contact:\n")
	formatContact(&sb, h.AdministrativeContact, "  ")

	sb.WriteString("\nTechnical Contact:\n")
	formatContact(&sb, h.TechnicalContact, "  ")

	sb.WriteString("\nBilling Contact:\n")
	formatContact(&sb, h.BillingContact, "  ")

	sb.WriteString("\nZone Contact:\n")
	formatContact(&sb, h.ZoneContact, "  ")

	// Raw text (truncated if too long)
	if len(h.RawText) > 0 {
		rawTextPreview := h.RawText
		if len(rawTextPreview) > 500 {
			rawTextPreview = rawTextPreview[:500] + "... [truncated]"
		}
		sb.WriteString("\nRaw Text:\n")
		sb.WriteString(rawTextPreview)
		sb.WriteString("\n")
	}

	if len(h.CleanText) > 0 {
		cleanTextPreview := h.CleanText
		if len(cleanTextPreview) > 500 {
			cleanTextPreview = cleanTextPreview[:500] + "... [truncated]"
		}
		sb.WriteString("\nClean Text:\n")
		sb.WriteString(cleanTextPreview)
		sb.WriteString("\n")
	}

	return sb.String()
}

// Helper function to format contact information
func formatContact(sb *strings.Builder, contact ContactInfo, indent string) {
	if contact.Name == "" && contact.Organization == "" && contact.Email == "" {
		sb.WriteString(indent + "No contact information available\n")
		return
	}

	if contact.Name != "" {
		sb.WriteString(indent + "Name: " + contact.Name + "\n")
	}
	if contact.Organization != "" {
		sb.WriteString(indent + "Organization: " + contact.Organization + "\n")
	}
	if contact.Email != "" {
		sb.WriteString(indent + "HunterEmail: " + contact.Email + "\n")
	}
	if contact.Street != "" {
		sb.WriteString(indent + "Street: " + contact.Street + "\n")
	}
	if contact.City != "" {
		sb.WriteString(indent + "City: " + contact.City + "\n")
	}
	if contact.State != "" {
		sb.WriteString(indent + "State: " + contact.State + "\n")
	}
	if contact.PostalCode != "" {
		sb.WriteString(indent + "Postal Code: " + contact.PostalCode + "\n")
	}
	if contact.Country != "" {
		sb.WriteString(indent + "Country: " + contact.Country + "\n")
	}
	if contact.Telephone != "" {
		phone := contact.Telephone
		if contact.TelephoneExt != "" {
			phone += " ext. " + contact.TelephoneExt
		}
		sb.WriteString(indent + "Telephone: " + phone + "\n")
	}
	if contact.Fax != "" {
		fax := contact.Fax
		if contact.FaxExt != "" {
			fax += " ext. " + contact.FaxExt
		}
		sb.WriteString(indent + "Fax: " + fax + "\n")
	}
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

type WhoIsIPLookup struct {
	RemainingCredits int    `json:"remaining_credits"`
	Data             IPData `json:"data"`
}

type WhoIsMXLookup struct {
	RemainingCredits int    `json:"remaining_credits"`
	Data             IPData `json:"data"`
}

type WhoIsNSLookup struct {
	RemainingCredits int    `json:"remaining_credits"`
	Data             IPData `json:"data"`
}

type IPData struct {
	CurrentPage string         `json:"current_page"`
	Result      []LookupResult `json:"result"`
	Size        int            `json:"size"`
}

type LookupResult struct {
	gorm.Model
	FirstSeen  int64  `json:"first_seen"`
	LastVisit  int64  `json:"last_visit"`
	Name       string `json:"name" gorm:"unique"`
	SearchTerm string `json:"search_term,omitempty"` // For storing the IP address this domain is associated with
	Type       string `json:"type,omitempty"`        // For storing the MX address this domain is associated with
}

func (LookupResult) TableName() string {
	return "lookup"
}

// Helper function to format contact information for WhoisRecord
func formatWhoisContact(sb *strings.Builder, contact Contact, indent string) {
	if contact.Name == "" && contact.Organization == "" {
		sb.WriteString(indent + "No contact information available\n")
		return
	}

	if contact.Name != "" {
		sb.WriteString(indent + "Name: " + contact.Name + "\n")
	}
	if contact.Organization != "" {
		sb.WriteString(indent + "Organization: " + contact.Organization + "\n")
	}
	if contact.Street1 != "" {
		sb.WriteString(indent + "Street: " + contact.Street1 + "\n")
	}
	if contact.City != "" {
		sb.WriteString(indent + "City: " + contact.City + "\n")
	}
	if contact.State != "" {
		sb.WriteString(indent + "State: " + contact.State + "\n")
	}
	if contact.PostalCode != "" {
		sb.WriteString(indent + "Postal Code: " + contact.PostalCode + "\n")
	}
	if contact.Country != "" {
		sb.WriteString(indent + "Country: " + contact.Country + "\n")
	}
	if contact.CountryCode != "" {
		sb.WriteString(indent + "Country Code: " + contact.CountryCode + "\n")
	}
	if contact.Telephone != "" {
		sb.WriteString(indent + "Telephone: " + contact.Telephone + "\n")
	}
	if contact.RawText != "" {
		rawTextPreview := contact.RawText
		if len(rawTextPreview) > 100 {
			rawTextPreview = rawTextPreview[:100] + "... [truncated]"
		}
		sb.WriteString(indent + "Raw Text: " + rawTextPreview + "\n")
	}
}
