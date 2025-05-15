package cmd

import (
	"dehasher/internal/badger"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"strings"
)

var (
	// Global Flags
	useLocalDatabase bool

	// rootCmd is the base command for the CLI.
	rootCmd = &cobra.Command{
		Use:   "dehasher",
		Short: `Dehasher is a cli tool for querying query.`,
		Long: fmt.Sprintf(
			"%s\n%s",
			`
 ______   _______           _______  _______           _______  _______ 
(  __  \ (  ____ \|\     /|(  ___  )(  ____ \|\     /|(  ____ \(  ____ )
| (  \  )| (    \/| )   ( || (   ) || (    \/| )   ( || (    \/| (    )|
| |   ) || (__    | (___) || (___) || (_____ | (___) || (__    | (____)|
| |   | ||  __)   |  ___  ||  ___  |(_____  )|  ___  ||  __)   |     __)
| |   ) || (      | (   ) || (   ) |      ) || (   ) || (      | (\ (   
| (__/  )| (____/\| )   ( || )   ( |/\____) || )   ( || (____/\| ) \ \__
(______/ (_______/|/     \||/     \|\_______)|/     \|(_______/|/   \__/
						     An Ar1ste1a Project                                                                        
`,
			`––•–√\/––√\/––•––––•–√\/––√\/––•––––•–√\/––√\/––•––√\/––•––––•–√\/––√\/––•––
  Dehasher can query the query API for:
  - Emails		- Usernames 		- Password
  - Hashes 		- IP Addresses		- Names
  - VINs		- License Plates	- Addresses
  - Phones		- Social Media		- Crypto Currency Addresses
  Dehasher supports:
  - Regex Matching
  - Exact Matching
––•–√\/––√\/––•––––•–√\/––√\/––•––––•–√\/––√\/––•––√\/––•––––•–√\/––√\/––•––
`,
		),
		Version: "v1.0",
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		zap.L().Fatal("execute_root_command",
			zap.String("message", "failed to execute root command"),
			zap.Error(err),
		)
		fmt.Printf("[!] %v", err)
		os.Exit(1)
	}
}

func init() {
	// Attempt to retreive the useLocalDB
	useLocalDatabase = badger.GetUseLocalDB()

	// Hide the default help command
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	// Add subcommands
	rootCmd.AddCommand(setKeyCmd)
	rootCmd.AddCommand(setLocalDb)
}

// Command to set API key
var setKeyCmd = &cobra.Command{
	Use:   "set-key [key]",
	Short: "Set and store API key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		// Store key in badger DB
		err := storeApiKey(key)
		if err != nil {
			fmt.Printf("Error storing API key: %v\n", err)
			return
		}
		fmt.Println("API key stored successfully")
	},
}

var setLocalDb = &cobra.Command{
	Use:   "set-local-db [true|false]",
	Short: "Set dehasher to use a local database path instead of the default (must be unset to use default)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		useLocal := strings.ToLower(args[0])

		if useLocal == "true" {
			useLocalDatabase = true
		} else if useLocal == "false" {
			useLocalDatabase = false
		} else {
			fmt.Println("Invalid argument. Please use 'true' or 'false'.")
			return
		}

		// Store useLocal in badger DB
		err := badger.StoreUseLocalDB(useLocalDatabase)
		if err != nil {
			fmt.Printf("Error storing local database useLocal: %v\n", err)
			return
		}
		fmt.Println("Local database useLocal stored successfully")
	},
}

// Helper functions to store API credentials
func storeApiKey(key string) error {
	err := badger.StoreKey(key)
	if err != nil {
		fmt.Printf("Error storing API key: %v\n", err)
		return err
	}
	return nil
}
