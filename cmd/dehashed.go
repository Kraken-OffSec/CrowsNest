package cmd

import (
	"crowsnest/internal/badger"
	"crowsnest/internal/debug"
	"crowsnest/internal/dehashed"
	"crowsnest/internal/sqlite"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	// Add api command to root command
	rootCmd.AddCommand(dehashedCmd)

	// Add flags specific to api command
	dehashedCmd.Flags().IntVarP(&maxRecords, "max-records", "m", 30000, "Maximum amount of records to return")
	dehashedCmd.Flags().IntVarP(&maxRequests, "max-requests", "r", -1, "Maximum number of requests to make")
	dehashedCmd.Flags().IntVarP(&startingPage, "starting-page", "s", 1, "Starting page for requests")
	dehashedCmd.Flags().BoolVarP(&printBalance, "print-balance", "b", false, "Print remaining balance after requests")
	dehashedCmd.Flags().BoolVarP(&regexMatch, "regex-match", "R", false, "Use regex matching on query fields")
	dehashedCmd.Flags().BoolVarP(&wildcardMatch, "wildcard-match", "W", false, "Use wildcard matching on query fields (Use ? to replace a single character, and * for multiple characters)")
	dehashedCmd.Flags().BoolVarP(&credsOnly, "creds-only", "C", false, "Return credentials only")
	dehashedCmd.Flags().StringVarP(&outputFormat, "format", "f", "json", "Output format (json, yaml, xml, txt)")
	dehashedCmd.Flags().StringVarP(&outputFile, "output", "o", "query", "File to output results to including extension")
	dehashedCmd.Flags().StringVarP(&usernameQuery, "username", "U", "", "Username query")
	dehashedCmd.Flags().StringVarP(&emailQuery, "email-query", "E", "", "HunterEmail query")
	dehashedCmd.Flags().StringVarP(&ipQuery, "ip", "I", "", "IP address query")
	dehashedCmd.Flags().StringVarP(&domainQuery, "domain", "D", "", "Domain query")
	dehashedCmd.Flags().StringVarP(&passwordQuery, "password", "P", "", "Password query")
	dehashedCmd.Flags().StringVarP(&vinQuery, "vin", "V", "", "VIN query")
	dehashedCmd.Flags().StringVarP(&licensePlateQuery, "license", "L", "", "License plate query")
	dehashedCmd.Flags().StringVarP(&addressQuery, "address", "A", "", "Address query")
	dehashedCmd.Flags().StringVarP(&phoneQuery, "phone", "M", "", "Phone query")
	dehashedCmd.Flags().StringVarP(&socialQuery, "social", "S", "", "Social query")
	dehashedCmd.Flags().StringVarP(&cryptoCurrencyAddressQuery, "crypto", "B", "", "Crypto currency address query")
	dehashedCmd.Flags().StringVarP(&hashQuery, "hash", "Q", "", "Hashed password query")
	dehashedCmd.Flags().StringVarP(&nameQuery, "name", "N", "", "Name query")

	// Add mutually exclusive flags to wildcard match and regex match
	dehashedCmd.MarkFlagsMutuallyExclusive("regex-match", "wildcard-match")
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
	dehashedCmd = &cobra.Command{
		Use:   "dehashed",
		Short: "Query the Dehashed API",
		Long:  `Query the Dehashed API for emails, usernames, passwords, hashes, IP addresses, and names.`,
		Run: func(cmd *cobra.Command, args []string) {
			key := getDehashedApiKey()

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
func getDehashedApiKey() string {
	return badger.GetDehashedKey()
}
