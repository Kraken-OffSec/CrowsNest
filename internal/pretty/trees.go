package pretty

import (
	"crowsnest/internal/sqlite"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

func WhoIsTree(root string, record sqlite.WhoisRecord) {
	enumeratorStyle := lipgloss.NewStyle().Foreground(purple).MarginRight(1)
	rootStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	itemStyle := lipgloss.NewStyle().Foreground(gray)

	rootTree := tree.Root(root)

	// Child Trees
	// Root Audit Tree
	auditTree := tree.Root("Audit")
	auditTree.Child(fmt.Sprintf("Created Date: %s", record.Audit.CreatedDate))
	auditTree.Child(fmt.Sprintf("Updated Date: %s", record.Audit.UpdatedDate))
	rootTree.Child(auditTree)

	// Root Name Servers Tree
	nameServersTree := tree.Root("Name Servers")
	nameServersTree.Child("Host Names: " + fmt.Sprintf("%v", record.NameServers.HostNames))
	nameServersTree.Child("IPs: " + fmt.Sprintf("%v", record.NameServers.IPs))
	nameServersTree.Child("Raw Text: " + record.NameServers.RawText)

	// Root Registry Data Tree
	registryDataTree := tree.Root("Registry Data")
	registryDataTree.Child("Audit: " + fmt.Sprintf("%v", record.RegistryData.Audit))
	registryDataTree.Child("Created Date: " + record.RegistryData.CreatedDate)
	registryDataTree.Child("Created Date Normalized: " + record.RegistryData.CreatedDateNormalized)
	registryDataTree.Child("Domain Name: " + record.RegistryData.DomainName)
	registryDataTree.Child("Expires Date: " + record.RegistryData.ExpiresDate)
	registryDataTree.Child("Expires Date Normalized: " + record.RegistryData.ExpiresDateNormalized)
	registryDataTree.Child("Footer: " + record.RegistryData.Footer)
	registryDataTree.Child("Header: " + record.RegistryData.Header)

	// Registry Data Name Servers Tree
	registryNameServersTree := tree.Root("Name Servers")
	registryNameServersTree.Child("Host Names: " + fmt.Sprintf("%v", record.RegistryData.NameServers.HostNames))
	registryNameServersTree.Child("IPs: " + fmt.Sprintf("%v", record.RegistryData.NameServers.IPs))
	registryNameServersTree.Child("Raw Text: " + record.RegistryData.NameServers.RawText)
	registryDataTree.Child(registryNameServersTree)

	// Root Registry Data Tree
	registryDataTree.Child("Parse Code: " + fmt.Sprintf("%d", record.RegistryData.ParseCode))
	registryDataTree.Child("Raw Text: " + record.RegistryData.RawText)
	registryDataTree.Child("Registrar IANA ID: " + record.RegistryData.RegistrarIANAID)
	registryDataTree.Child("Registrar Name: " + record.RegistryData.RegistrarName)
	registryDataTree.Child("Status: " + record.RegistryData.Status)
	registryDataTree.Child("Stripped Text: " + record.RegistryData.StrippedText)
	registryDataTree.Child("Updated Date: " + record.RegistryData.UpdatedDate)
	registryDataTree.Child("Updated Date Normalized: " + record.RegistryData.UpdatedDateNormalized)
	registryDataTree.Child("Whois Server: " + record.RegistryData.WhoisServer)

	// Root Contract Tree
	technicalContactTree := tree.Root("Technical Contact")
	technicalContactTree.Child("City: " + record.TechnicalContact.City)
	technicalContactTree.Child("Country: " + record.TechnicalContact.Country)
	technicalContactTree.Child("Country Code: " + record.TechnicalContact.CountryCode)
	technicalContactTree.Child("Name: " + record.TechnicalContact.Name)
	technicalContactTree.Child("Organization: " + record.TechnicalContact.Organization)
	technicalContactTree.Child("Postal Code: " + record.TechnicalContact.PostalCode)
	technicalContactTree.Child("Raw Text: " + record.TechnicalContact.RawText)
	technicalContactTree.Child("State: " + record.TechnicalContact.State)
	technicalContactTree.Child("Street 1: " + record.TechnicalContact.Street1)
	technicalContactTree.Child("Telephone: " + record.TechnicalContact.Telephone)

	// Root Tree Children
	rootTree.Child("Contact HunterEmail: " + record.ContactEmail)
	rootTree.Child("Created Date: " + record.CreatedDate)
	rootTree.Child("Created Date Normalized: " + record.CreatedDateNormalized)
	rootTree.Child("Domain Name: " + record.DomainName)
	rootTree.Child("Domain Name Ext: " + record.DomainNameExt)
	rootTree.Child("Estimated Domain Age: " + fmt.Sprintf("%d", record.EstimatedDomainAge))
	rootTree.Child("Expires Date: " + record.ExpiresDate)
	rootTree.Child("Expires Date Normalized: " + record.ExpiresDateNormalized)
	rootTree.Child("Footer: " + record.Footer)
	rootTree.Child("Header: " + record.Header)
	rootTree.Child(nameServersTree)
	rootTree.Child("Parse Code: " + fmt.Sprintf("%d", record.ParseCode))
	rootTree.Child("Raw Text: " + record.RawText)
	rootTree.Child("Registrant: " + fmt.Sprintf("%v", record.Registrant))
	rootTree.Child("Registrar IANA ID: " + record.RegistrarIANAID)
	rootTree.Child("Registrar Name: " + record.RegistrarName)
	rootTree.Child(registryDataTree)
	rootTree.Child("Status: " + record.Status)
	rootTree.Child("Stripped Text: " + record.StrippedText)
	rootTree.Child(technicalContactTree)
	rootTree.Child("Updated Date: " + record.UpdatedDate)
	rootTree.Child("Updated Date Normalized: " + record.UpdatedDateNormalized)

	// Styles
	rootTree.Enumerator(tree.RoundedEnumerator)
	rootTree.EnumeratorStyle(enumeratorStyle)
	rootTree.RootStyle(rootStyle)
	rootTree.ItemStyle(itemStyle)

	// Print Tree
	fmt.Println(rootTree)
}

func HunterDomainTree(root string, record sqlite.HunterDomainData) {
	enumeratorStyle := lipgloss.NewStyle().Foreground(purple).MarginRight(1)
	rootStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	itemStyle := lipgloss.NewStyle().Foreground(gray)

	rootTree := tree.Root(root)

	// Root Tree Children
	rootTree.Child("Domain: " + record.Domain)
	rootTree.Child("Disposable: " + fmt.Sprintf("%t", record.Disposable))
	rootTree.Child("Webmail: " + fmt.Sprintf("%t", record.Webmail))
	rootTree.Child("Accept All: " + fmt.Sprintf("%t", record.AcceptAll))
	rootTree.Child("Pattern: " + record.Pattern)
	rootTree.Child("Organization: " + record.Organization)
	rootTree.Child("Description: " + record.Description)
	rootTree.Child("Industry: " + record.Industry)
	rootTree.Child("Twitter: " + record.Twitter)
	rootTree.Child("Facebook: " + record.Facebook)
	rootTree.Child("Linkedin: " + record.Linkedin)
	rootTree.Child("Instagram: " + record.Instagram)
	rootTree.Child("Youtube: " + record.Youtube)

	techTree := tree.Root("Technologies")
	for _, tech := range record.Technologies {
		techTree.Child(tech)
	}
	rootTree.Child(techTree)

	rootTree.Child("Country: " + record.Country)
	rootTree.Child("State: " + record.State)
	rootTree.Child("City: " + record.City)
	rootTree.Child("Postal Code: " + record.PostalCode)
	rootTree.Child("Street: " + record.Street)
	rootTree.Child("Headcount: " + record.Headcount)
	rootTree.Child("Company Type: " + record.CompanyType)

	emailTree := tree.Root("Emails")
	for _, email := range record.Emails {
		emailTree.Child(email.ToTree())
	}
	rootTree.Child(emailTree)

	linkedDomainTree := tree.Root("Linked Domains")
	for _, domain := range record.LinkedDomains {
		linkedDomainTree.Child(domain)
	}
	rootTree.Child(linkedDomainTree)

	// Styles
	rootTree.Enumerator(tree.RoundedEnumerator)
	rootTree.EnumeratorStyle(enumeratorStyle)
	rootTree.RootStyle(rootStyle)
	rootTree.ItemStyle(itemStyle)

	// Print Tree
	fmt.Println(rootTree)
}

func HunterCompanyEnrichmentTree(root string, record sqlite.CompanyData) {
	enumeratorStyle := lipgloss.NewStyle().Foreground(purple).MarginRight(1)
	rootStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	itemStyle := lipgloss.NewStyle().Foreground(gray)

	rootTree := tree.Root(root)

	// Root Tree Children
	rootTree.Child("ID: " + record.ID)
	rootTree.Child("Name: " + record.Name)
	rootTree.Child("Legal Name: " + record.LegalName)
	rootTree.Child("Domain: " + record.Domain)
	rootTree.Child(record.DomainAliasesTree())
	rootTree.Child(record.SiteTree())
	rootTree.Child(record.CategoryTree())
	rootTree.Child(record.TagsTree())
	rootTree.Child("Description: " + record.Description)
	rootTree.Child("Founded Year: " + fmt.Sprintf("%d", record.FoundedYear))
	rootTree.Child("Location: " + record.Location)
	rootTree.Child("Time Zone: " + record.TimeZone)
	rootTree.Child("UTC Offset: " + fmt.Sprintf("%d", record.UTCOffset))
	rootTree.Child(record.GeoTree())
	rootTree.Child("Logo: " + record.Logo)
	rootTree.Child(record.FacebookTree())
	rootTree.Child(record.LinkedInTree())
	rootTree.Child(record.TwitterTree())
	rootTree.Child(record.CrunchbaseTree())
	rootTree.Child(record.YouTubeTree())
	rootTree.Child("Email Provider: " + record.EmailProvider)
	rootTree.Child("Type: " + record.Type)
	rootTree.Child("Ticker: " + record.Ticker)
	rootTree.Child(record.IdentifiersTree())
	rootTree.Child("Phone: " + record.Phone)
	rootTree.Child(record.MetricsTree())
	rootTree.Child("Indexed At: " + record.IndexedAt)
	rootTree.Child(record.TechTree())
	rootTree.Child(record.TechCategoriesTree())
	rootTree.Child(record.ParentTree())
	rootTree.Child(record.UltimateParentTree())

	// Styles
	rootTree.Enumerator(tree.RoundedEnumerator)
	rootTree.EnumeratorStyle(enumeratorStyle)
	rootTree.RootStyle(rootStyle)
	rootTree.ItemStyle(itemStyle)

	// Print Tree
	fmt.Println(rootTree)
}

func HunterPersonEnrichmentTree(root string, record sqlite.PersonData) {
	enumeratorStyle := lipgloss.NewStyle().Foreground(purple).MarginRight(1)
	rootStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	itemStyle := lipgloss.NewStyle().Foreground(gray)

	rootTree := tree.Root(root)

	// Root Tree Children
	rootTree.Child("ID: " + record.ID)
	rootTree.Child(record.NameTree())
	rootTree.Child("Email: " + record.Email)
	rootTree.Child("Location: " + record.Location)
	rootTree.Child("Time Zone: " + record.TimeZone)
	rootTree.Child("UTC Offset: " + fmt.Sprintf("%d", record.UTCOffset))
	rootTree.Child(record.GeoTree())
	rootTree.Child("Bio: " + record.Bio)
	rootTree.Child("Site: " + record.Site)
	rootTree.Child("Avatar: " + record.Avatar)
	rootTree.Child(record.EmploymentTree())
	rootTree.Child(record.FacebookTree())
	rootTree.Child(record.GitHubTree())
	rootTree.Child(record.TwitterTree())
	rootTree.Child(record.LinkedInTree())
	rootTree.Child(record.GooglePlusTree())
	rootTree.Child(record.GravatarTree())
	rootTree.Child("Fuzzy: " + fmt.Sprintf("%t", record.Fuzzy))
	rootTree.Child("Email Provider: " + record.EmailProvider)
	rootTree.Child("Indexed At: " + record.IndexedAt)
	rootTree.Child("Phone: " + record.Phone)
	rootTree.Child("Active At: " + record.ActiveAt)
	rootTree.Child("Inactive At: " + record.InactiveAt)

	// Styles
	rootTree.Enumerator(tree.RoundedEnumerator)
	rootTree.EnumeratorStyle(enumeratorStyle)
	rootTree.RootStyle(rootStyle)
	rootTree.ItemStyle(itemStyle)

	// Print Tree
	fmt.Println(rootTree)
}

func HunterCombinedEnrichmentTree(root string, record sqlite.CombinedData) {
	enumeratorStyle := lipgloss.NewStyle().Foreground(purple).MarginRight(1)
	rootStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	itemStyle := lipgloss.NewStyle().Foreground(gray)

	rootTree := tree.Root(root)

	// Root Tree Children
	rootTree.Child(personTree(record.Person))
	rootTree.Child(companyTree(record.Company))

	// Styles
	rootTree.Enumerator(tree.RoundedEnumerator)
	rootTree.EnumeratorStyle(enumeratorStyle)
	rootTree.RootStyle(rootStyle)
	rootTree.ItemStyle(itemStyle)

	// Print Tree
	fmt.Println(rootTree)
}

func companyTree(record sqlite.CompanyData) *tree.Tree {
	companyTree := tree.Root("Company")

	// Company Tree Children
	companyTree.Child("ID: " + record.ID)
	companyTree.Child("Name: " + record.Name)
	companyTree.Child("Legal Name: " + record.LegalName)
	companyTree.Child("Domain: " + record.Domain)
	companyTree.Child(record.DomainAliasesTree())
	companyTree.Child(record.SiteTree())
	companyTree.Child(record.CategoryTree())
	companyTree.Child(record.TagsTree())
	companyTree.Child("Description: " + record.Description)
	companyTree.Child("Founded Year: " + fmt.Sprintf("%d", record.FoundedYear))
	companyTree.Child("Location: " + record.Location)
	companyTree.Child("Time Zone: " + record.TimeZone)
	companyTree.Child("UTC Offset: " + fmt.Sprintf("%d", record.UTCOffset))
	companyTree.Child(record.GeoTree())
	companyTree.Child("Logo: " + record.Logo)
	companyTree.Child(record.FacebookTree())
	companyTree.Child(record.LinkedInTree())
	companyTree.Child(record.TwitterTree())
	companyTree.Child(record.CrunchbaseTree())
	companyTree.Child(record.YouTubeTree())
	companyTree.Child("Email Provider: " + record.EmailProvider)
	companyTree.Child("Type: " + record.Type)
	companyTree.Child("Ticker: " + record.Ticker)
	companyTree.Child(record.IdentifiersTree())
	companyTree.Child("Phone: " + record.Phone)
	companyTree.Child(record.MetricsTree())
	companyTree.Child("Indexed At: " + record.IndexedAt)
	companyTree.Child(record.TechTree())
	companyTree.Child(record.TechCategoriesTree())
	companyTree.Child(record.ParentTree())
	companyTree.Child(record.UltimateParentTree())

	return companyTree
}

func personTree(record sqlite.PersonData) *tree.Tree {
	personTree := tree.Root("Person")

	// Person Tree Children
	personTree.Child("ID: " + record.ID)
	personTree.Child(record.NameTree())
	personTree.Child("Email: " + record.Email)
	personTree.Child("Location: " + record.Location)
	personTree.Child("Time Zone: " + record.TimeZone)
	personTree.Child("UTC Offset: " + fmt.Sprintf("%d", record.UTCOffset))
	personTree.Child(record.GeoTree())
	personTree.Child("Bio: " + record.Bio)
	personTree.Child("Site: " + record.Site)
	personTree.Child("Avatar: " + record.Avatar)
	personTree.Child(record.EmploymentTree())
	personTree.Child(record.FacebookTree())
	personTree.Child(record.GitHubTree())
	personTree.Child(record.TwitterTree())
	personTree.Child(record.LinkedInTree())
	personTree.Child(record.GooglePlusTree())
	personTree.Child(record.GravatarTree())
	personTree.Child("Fuzzy: " + fmt.Sprintf("%t", record.Fuzzy))
	personTree.Child("Email Provider: " + record.EmailProvider)
	personTree.Child("Indexed At: " + record.IndexedAt)
	personTree.Child("Phone: " + record.Phone)
	personTree.Child("Active At: " + record.ActiveAt)
	personTree.Child("Inactive At: " + record.InactiveAt)

	return personTree
}
