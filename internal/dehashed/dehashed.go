package dehashed

import (
	"dehasher/internal/debug"
	"dehasher/internal/export"
	"dehasher/internal/sqlite"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
)

// Dehasher is a struct for querying the Dehashed API
type Dehasher struct {
	options   sqlite.QueryOptions
	nextPage  int
	debug     bool
	balance   int
	request   *DehashedSearchRequest
	client    *DehashedClientV2
	queryPlan []struct{ Page, Size int }
}

// NewDehasher creates a new Dehasher
func NewDehasher(options *sqlite.QueryOptions) *Dehasher {
	dh := &Dehasher{
		options:   *options,
		nextPage:  options.StartingPage + 1,
		debug:     options.Debug,
		balance:   0,
		queryPlan: make([]struct{ Page, Size int }, 0),
	}
	dh.setQueries()
	dh.request = NewDehashedSearchRequest(
		dh.queryPlan[0].Page,
		dh.queryPlan[0].Size,
		dh.options.WildcardMatch,
		dh.options.RegexMatch,
		false,
		options.Debug,
	)

	dh.buildRequest()
	return dh
}

// SetClientCredentials sets the client credentials for the dehasher
func (dh *Dehasher) SetClientCredentials(key string) {
	dh.client = NewDehashedClientV2(key, dh.debug)
}

func (dh *Dehasher) getNextPage() int {
	if dh.debug {
		debug.PrintInfo(fmt.Sprintf("getting next page: %d", dh.nextPage))
	}
	nextPage := dh.nextPage
	dh.nextPage += 1
	return nextPage
}

// generatePagination creates a list of (page, size) tuples such that page * size <= 10000
func generatePagination(maxRecords int) []struct{ Page, Size int } {
	const maxPageProduct = 9500
	var queries []struct{ Page, Size int }

	remaining := maxRecords
	page := 1

	for remaining > 0 {
		size := (maxPageProduct - 1) / page // guarantees page * size < 10000
		if size > remaining {
			size = remaining
		}
		queries = append(queries, struct{ Page, Size int }{page, size})
		remaining -= size
		page++
	}

	return queries
}

// setQueries sets the number of queries to make based on the number of records and requests
func (dh *Dehasher) setQueries() {
	if dh.options.MaxRecords <= 0 {
		dh.options.MaxRecords = 10000
	}

	dh.queryPlan = generatePagination(dh.options.MaxRecords)

	fmt.Printf("Making %d requests to retrieve %d records\n", len(dh.queryPlan), dh.options.MaxRecords)

	if dh.debug {
		for i, q := range dh.queryPlan {
			debug.PrintInfo(fmt.Sprintf("query %d: page=%d, size=%d", i+1, q.Page, q.Size))
		}
	}
}

// Start starts the querying process
func (dh *Dehasher) Start() {
	fmt.Printf("[*] Querying Dehashed API...\n")

	// Make initial request to get total count
	fmt.Printf("   [*] Performing initial request to determine total records...\n")
	totalRecords, balance, err := dh.client.Search(*dh.request)
	if err != nil {
		handleSearchError(dh, err)
		return
	}

	dh.balance = balance
	recordsRetrieved := len(dh.client.results)

	fmt.Printf("      [+] Retrieved %d records\n", recordsRetrieved)
	fmt.Printf("      [*] Total available records: %d\n", totalRecords)

	if dh.options.PrintBalance {
		fmt.Printf("      [*] Balance: %d\n", balance)
	}

	// If we've already got all records or reached our limit, we're done
	if recordsRetrieved >= totalRecords || recordsRetrieved >= dh.options.MaxRecords {
		fmt.Printf("      [*] All requested records retrieved\n")
		dh.parseResults()
		return
	}

	// Calculate remaining records to fetch
	remainingRecords := totalRecords - recordsRetrieved
	if dh.options.MaxRecords > 0 && dh.options.MaxRecords < totalRecords {
		remainingRecords = dh.options.MaxRecords - recordsRetrieved
	}

	// Check if we need user confirmation for large datasets
	if remainingRecords > 30000 {
		tokensRequired := (remainingRecords + 9999) / 10000 // Ceiling division
		fmt.Printf("\n[!] Large dataset detected: %d additional records\n", remainingRecords)
		fmt.Printf("[!] This will require approximately %d API tokens\n", tokensRequired)
		fmt.Printf("[!] Your current balance: %d\n", balance)

		if balance < tokensRequired {
			fmt.Printf("[!] WARNING: Your balance (%d) is less than required tokens (%d)\n", balance, tokensRequired)
		}

		fmt.Printf("[?] Do you want to continue? (y/n): ")
		var response string
		fmt.Scanln(&response)

		if response != "y" && response != "Y" {
			fmt.Println("[*] Operation cancelled by user")
			dh.parseResults()
			return
		}
	}

	// Make additional requests
	for i, q := range dh.queryPlan {
		if i == 0 {
			// We already made the first request before this loop
			continue
		}

		dh.request.Page = q.Page
		dh.request.Size = q.Size

		fmt.Printf("   [*] Performing Request %d of %d (page=%d, size=%d)...\n", i+1, len(dh.queryPlan), q.Page, q.Size)

		_, balance, err := dh.client.Search(*dh.request)
		if err != nil {
			handleSearchError(dh, err)
			break
		}

		dh.balance = balance
		recordsRetrieved += len(dh.client.results)

		fmt.Printf("      [+] Retrieved %d total records so far\n", recordsRetrieved)

		if dh.options.PrintBalance {
			fmt.Printf("      [*] Balance: %d\n", balance)
		}

		if recordsRetrieved >= totalRecords || recordsRetrieved >= dh.options.MaxRecords {
			fmt.Printf("      [*] All requested records retrieved\n")
			break
		}
	}

	dh.parseResults()
}

// Helper function to handle search errors
func handleSearchError(dh *Dehasher, err error) {
	if dh.debug {
		debug.PrintInfo("error performing request")
		debug.PrintError(err)
	}

	// Check if it's a DehashError
	if dhErr, ok := err.(*DehashError); ok {
		fmt.Printf("      [!] Dehashed API Error: %s (Code: %d)\n", dhErr.Message, dhErr.Code)
		zap.L().Error("dehashed_api_error",
			zap.String("message", dhErr.Message),
			zap.Int("code", dhErr.Code),
		)
	} else {
		fmt.Printf("   [!] Error performing request: %v\n", err)
		zap.L().Error("request_error",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
	}

	if len(dh.client.results) > 0 {
		fmt.Printf("   [!] Partial results retrieved. Storing Results...\n")
		err := sqlite.StoreResults(dh.client.GetResults())
		if err != nil {
			zap.L().Error("store_results",
				zap.String("message", "failed to store results"),
				zap.Error(err),
			)
			fmt.Printf("   [!] Error storing results: %v\n", err)
		}
	}
}

// buildRequest constructs the query map
func (dh *Dehasher) buildRequest() {
	if len(dh.options.UsernameQuery) > 0 {
		dh.request.AddUsernameQuery(dh.options.UsernameQuery)
	}
	if len(dh.options.EmailQuery) > 0 {
		dh.request.AddEmailQuery(dh.options.EmailQuery)
	}
	if len(dh.options.IpQuery) > 0 {
		dh.request.AddIpAddressQuery(dh.options.IpQuery)
	}
	if len(dh.options.HashQuery) > 0 {
		dh.request.AddHashedPasswordQuery(dh.options.HashQuery)
	}
	if len(dh.options.PassQuery) > 0 {
		dh.request.AddPasswordQuery(dh.options.PassQuery)
	}
	if len(dh.options.NameQuery) > 0 {
		dh.request.AddNameQuery(dh.options.NameQuery)
	}
	if len(dh.options.DomainQuery) > 0 {
		dh.request.AddDomainQuery(dh.options.DomainQuery)
	}
	if len(dh.options.VinQuery) > 0 {
		dh.request.AddVinQuery(dh.options.VinQuery)
	}
	if len(dh.options.LicensePlateQuery) > 0 {
		dh.request.AddLicensePlateQuery(dh.options.LicensePlateQuery)
	}
	if len(dh.options.AddressQuery) > 0 {
		dh.request.AddAddressQuery(dh.options.AddressQuery)
	}
	if len(dh.options.PhoneQuery) > 0 {
		dh.request.AddPhoneQuery(dh.options.PhoneQuery)
	}
	if len(dh.options.SocialQuery) > 0 {
		dh.request.AddSocialQuery(dh.options.SocialQuery)
	}
	if len(dh.options.CryptoAddressQuery) > 0 {
		dh.request.AddCryptoAddressQuery(dh.options.CryptoAddressQuery)
	}
}

// parseResults parses the results and writes them to a file
func (dh *Dehasher) parseResults() {
	var data []byte

	zap.L().Info("extracting_credentials")
	results := dh.client.GetResults()
	creds := results.ExtractCredentials()
	fmt.Printf("\n\t[+] Discovered %d Credentials", len(creds))
	err := sqlite.StoreCreds(creds)
	if err != nil {
		zap.L().Error("store_creds",
			zap.String("message", "failed to store creds"),
			zap.Error(err),
		)
	}
	zap.L().Info("creds_stored", zap.Int("count", len(creds)))

	zap.L().Info("storing_results")
	err = sqlite.StoreResults(results)
	if err != nil {
		zap.L().Error("store_results",
			zap.String("message", "failed to store results"),
			zap.Error(err),
		)
	}
	zap.L().Info("results_stored", zap.Int("count", len(results.Results)))

	if len(results.Results) > 0 {
		fmt.Printf("\n\t[*] Writing entries to file: %s.%s", dh.options.OutputFile, dh.options.OutputFormat.String())
		if !dh.options.CredsOnly {
			err := export.WriteToFile(results, dh.options.OutputFile, dh.options.OutputFormat)
			if err != nil {
				fmt.Printf("\n[!] Error Writing to file: %v\n\tOutputting to terminal.", err)
				data, err = json.MarshalIndent(results, "", "  ")
				fmt.Println(string(data))
				os.Exit(0)
			} else {
				fmt.Println("\n\t\t[*] Success\n")
			}
		} else {
			creds := results.ExtractCredentials()
			err := export.WriteCredsToFile(creds, dh.options.OutputFile, dh.options.OutputFormat)
			if err != nil {
				fmt.Printf("\n[!] Error Writing to file: %v\n\tOutputting to terminal.", err)
				data, err = json.MarshalIndent(creds, "", "  ")
				fmt.Println(string(data))
				os.Exit(0)
			} else {
				fmt.Println("\n\t\t[*] Success\n")
			}
		}
	}
}
