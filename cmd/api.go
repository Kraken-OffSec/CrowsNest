package cmd

import (
	"dehasher/internal/badger"
	"dehasher/internal/debug"
	"dehasher/internal/dehashed"
	"dehasher/internal/sqlite"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	// Add query command to root command
	rootCmd.AddCommand(apiCmd)

	// Add flags specific to api command
	apiCmd.Flags().IntVarP(&maxRecords, "max-records", "m", 30000, "Maximum amount of records to return")
	apiCmd.Flags().IntVarP(&maxRequests, "max-requests", "r", -1, "Maximum number of requests to make")
	apiCmd.Flags().IntVarP(&startingPage, "starting-page", "s", 1, "Starting page for requests")
	apiCmd.Flags().BoolVarP(&printBalance, "print-balance", "b", false, "Print remaining balance after requests")
	apiCmd.Flags().BoolVarP(&regexMatch, "regex-match", "R", false, "Use regex matching on query fields")
	apiCmd.Flags().BoolVarP(&wildcardMatch, "wildcard-match", "W", false, "Use wildcard matching on query fields (Use ? to replace a single character, and * for multiple characters)")
	apiCmd.Flags().BoolVarP(&credsOnly, "creds-only", "C", false, "Return credentials only")
	apiCmd.Flags().StringVarP(&outputFormat, "format", "f", "json", "Output format (json, yaml, xml, txt)")
	apiCmd.Flags().StringVarP(&outputFile, "output", "o", "query", "File to output results to including extension")
	apiCmd.Flags().StringVarP(&usernameQuery, "username", "U", "", "Username query")
	apiCmd.Flags().StringVarP(&emailQuery, "email-query", "E", "", "Email query")
	apiCmd.Flags().StringVarP(&ipQuery, "ip", "I", "", "IP address query")
	apiCmd.Flags().StringVarP(&domainQuery, "domain", "D", "", "Domain query")
	apiCmd.Flags().StringVarP(&passwordQuery, "password", "P", "", "Password query")
	apiCmd.Flags().StringVarP(&vinQuery, "vin", "V", "", "VIN query")
	apiCmd.Flags().StringVarP(&licensePlateQuery, "license", "L", "", "License plate query")
	apiCmd.Flags().StringVarP(&addressQuery, "address", "A", "", "Address query")
	apiCmd.Flags().StringVarP(&phoneQuery, "phone", "M", "", "Phone query")
	apiCmd.Flags().StringVarP(&socialQuery, "social", "S", "", "Social query")
	apiCmd.Flags().StringVarP(&cryptoCurrencyAddressQuery, "crypto", "B", "", "Crypto currency address query")
	apiCmd.Flags().StringVarP(&hashQuery, "hash", "Q", "", "Hashed password query")
	apiCmd.Flags().StringVarP(&nameQuery, "name", "N", "", "Name query")

	// Add mutually exclusive flags to wildcard match and regex match
	apiCmd.MarkFlagsMutuallyExclusive("regex-match", "wildcard-match")
}

var (
	// Query command flags
	maxRecords                 int
	maxRequests                int
	startingPage               int
	credsOnly                  bool
	printBalance               bool
	regexMatch                 bool
	wildcardMatch              bool
	outputFormat               string
	outputFile                 string
	usernameQuery              string
	emailQuery                 string
	ipQuery                    string
	passwordQuery              string
	hashQuery                  string
	nameQuery                  string
	domainQuery                string
	vinQuery                   string
	licensePlateQuery          string
	addressQuery               string
	phoneQuery                 string
	socialQuery                string
	cryptoCurrencyAddressQuery string

	// Query command
	apiCmd = &cobra.Command{
		Use:   "api",
		Short: "Query the Dehashed API",
		Long:  `Query the Dehashed API for emails, usernames, passwords, hashes, IP addresses, and names.`,
		Run: func(cmd *cobra.Command, args []string) {
			key := getStoredApiKey()

			// Validate credentials
			if key == "" {
				fmt.Println("API key is required. Set the key with the \"set-key\" command. [dehasher set-key <api_key>]")
				return
			}

			// Create new QueryOptions
			queryOptions := sqlite.NewQueryOptions(
				maxRecords,
				maxRequests,
				startingPage,
				outputFormat,
				outputFile,
				usernameQuery,
				emailQuery,
				ipQuery,
				passwordQuery,
				hashQuery,
				nameQuery,
				domainQuery,
				vinQuery,
				licensePlateQuery,
				addressQuery,
				phoneQuery,
				socialQuery,
				cryptoCurrencyAddressQuery,
				regexMatch,
				wildcardMatch,
				printBalance,
				credsOnly,
				debugGlobal,
			)

			// Create new Dehasher
			dehasher := dehashed.NewDehasher(queryOptions)
			dehasher.SetClientCredentials(
				key,
			)

			// Start querying
			dehasher.Start()
			fmt.Println("\n[*] Completing Process")

			err := sqlite.StoreQueryOptions(queryOptions)
			if err != nil {
				if debugGlobal {
					debug.PrintInfo("failed to store query options")
					debug.PrintError(err)
				}
				zap.L().Error("store_query_options",
					zap.String("message", "failed to store query options"),
					zap.Error(err),
				)
				fmt.Printf("Error storing query options: %v\n", err)
			}
		},
	}
)

// Helper functions to get stored API credentials
func getStoredApiKey() string {
	return badger.GetKey()
}
