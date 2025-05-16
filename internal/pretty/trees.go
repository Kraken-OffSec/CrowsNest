package pretty

import (
	"dehasher/internal/sqlite"
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
	rootTree.Child("Contact Email: " + record.ContactEmail)
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
