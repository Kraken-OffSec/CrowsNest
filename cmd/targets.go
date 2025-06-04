package cmd

import (
	"crowsnest/internal/sqlite"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"strings"
)

func init() {
	// Add targets command to root command
	rootCmd.AddCommand(targetsCmd)

	// Add flags specific to targets command
	targetsCmd.Flags().StringVarP(&targetsOutputFile, "output", "o", "targets", "Output file name (required)")
	targetsCmd.Flags().BoolVarP(&targetsExternal, "external", "e", false, "Output external format (email:password)")
	targetsCmd.Flags().BoolVarP(&targetsInternal, "internal", "i", false, "Output internal format (username:password)")
	targetsCmd.Flags().BoolVarP(&targetsSubdomains, "subdomains", "s", false, "Output subdomains")
	targetsCmd.Flags().BoolVarP(&targetsEmails, "emails", "E", false, "Output emails only (no passwords)")
	targetsCmd.Flags().StringVarP(&targetsDomain, "domain", "d", "", "Filter by domain (for emails and subdomains)")

	// Add mutually exclusive flags to targets command
	targetsCmd.MarkFlagsMutuallyExclusive("external", "internal", "subdomains", "emails")
}

var (
	// Targets command flags
	targetsOutputFile string
	targetsExternal   bool
	targetsInternal   bool
	targetsSubdomains bool
	targetsEmails     bool
	targetsDomain     string

	// Targets command
	targetsCmd = &cobra.Command{
		Use:   "targets",
		Short: "Export users and subdomains in formats suitable for external tools",
		Long: `Export users and subdomains from the database in easily digestible formats for tools like sprays or other security testing tools.

Formats:
  --external (-e): Output in email:password format
  --internal (-i): Output in username:password format
  --emails (-E): Output emails only (no passwords)
  --subdomains (-s): Output subdomains only

Options:
  --domain (-d): Filter results by domain (applies to emails and subdomains)
  --output (-o): Specify output file name (required)

Examples:
  # Export all external credentials (email:password)
  crowsnest targets -e -o external_creds

  # Export internal credentials for a specific domain
  crowsnest targets -i -d example.com -o internal_creds

  # Export all emails
  crowsnest targets -E -o all_emails

  # Export emails for a specific domain
  crowsnest targets -E -d example.com -o domain_emails

  # Export subdomains for a specific domain
  crowsnest targets -s -d example.com -o subdomains

  # Export all subdomains
  crowsnest targets -s -o all_subdomains`,
		Run: func(cmd *cobra.Command, args []string) {
			// Validate that at least one format is specified
			if !targetsExternal && !targetsInternal && !targetsSubdomains && !targetsEmails {
				fmt.Println("[!] Error: You must specify at least one output format:")
				fmt.Println("    --external (-e) for email:password format")
				fmt.Println("    --internal (-i) for username:password format")
				fmt.Println("    --emails (-E) for emails only")
				fmt.Println("    --subdomains (-s) for subdomains")
				return
			}

			if debugGlobal {
				zap.L().Info("targets_debug",
					zap.String("message", "targets command started"),
					zap.Bool("external", targetsExternal),
					zap.Bool("internal", targetsInternal),
					zap.Bool("subdomains", targetsSubdomains),
					zap.Bool("emails", targetsEmails),
					zap.String("domain", targetsDomain),
					zap.String("output_file", targetsOutputFile),
				)
			}

			// Execute the targets export
			err := executeTargetsExport()
			if err != nil {
				fmt.Printf("[!] Error: %v\n", err)
				return
			}

			fmt.Printf("[+] Successfully exported targets to: %s\n", targetsOutputFile)
		},
	}
)

// executeTargetsExport performs the main logic for exporting targets
func executeTargetsExport() error {
	var outputLines []string

	// Export external credentials (email:password)
	if targetsExternal {
		if debugGlobal {
			fmt.Println("[*] Exporting external credentials (email:password)...")
		}

		externalCreds, err := getExternalCredentials()
		if err != nil {
			return fmt.Errorf("failed to get external credentials: %v", err)
		}

		for _, cred := range externalCreds {
			if cred.Email != "" && cred.Password != "" {
				outputLines = append(outputLines, fmt.Sprintf("%s:%s", cred.Email, cred.Password))
			}
		}

		if debugGlobal {
			fmt.Printf("[*] Found %d external credentials\n", len(externalCreds))
		}
	}

	// Export internal credentials (username:password)
	if targetsInternal {
		if debugGlobal {
			fmt.Println("[*] Exporting internal credentials (username:password)...")
		}

		internalCreds, err := getInternalCredentials()
		if err != nil {
			return fmt.Errorf("failed to get internal credentials: %v", err)
		}

		for _, cred := range internalCreds {
			if cred.Username != "" && cred.Password != "" {
				outputLines = append(outputLines, fmt.Sprintf("%s:%s", cred.Username, cred.Password))
			}
		}

		if debugGlobal {
			fmt.Printf("[*] Found %d internal credentials\n", len(internalCreds))
		}
	}

	// Export emails only
	if targetsEmails {
		if debugGlobal {
			fmt.Println("[*] Exporting emails only...")
		}

		emails, err := getEmailsOnly()
		if err != nil {
			return fmt.Errorf("failed to get emails: %v", err)
		}

		for _, email := range emails {
			if email.Email != "" {
				outputLines = append(outputLines, email.Email)
			}
		}

		if debugGlobal {
			fmt.Printf("[*] Found %d emails\n", len(emails))
		}
	}

	// Export subdomains
	if targetsSubdomains {
		if debugGlobal {
			fmt.Println("[*] Exporting subdomains...")
		}

		subdomains, err := getSubdomains()
		if err != nil {
			return fmt.Errorf("failed to get subdomains: %v", err)
		}

		for _, subdomain := range subdomains {
			if subdomain.Subdomain != "" {
				outputLines = append(outputLines, subdomain.Subdomain)
			}
		}

		if debugGlobal {
			fmt.Printf("[*] Found %d subdomains\n", len(subdomains))
		}
	}

	// Write to file
	if len(outputLines) == 0 {
		return fmt.Errorf("no data found to export")
	}

	// Join all lines with newlines and add a single newline at the end
	content := strings.Join(outputLines, "\n") + "\n"

	err := os.WriteFile(targetsOutputFile, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	if debugGlobal {
		fmt.Printf("[*] Wrote %d lines to %s\n", len(outputLines), targetsOutputFile)
	}

	return nil
}

// getExternalCredentials retrieves credentials for external format (email:password)
func getExternalCredentials() ([]sqlite.User, error) {
	db := sqlite.GetDB()
	var users []sqlite.User

	query := db.Where("email IS NOT NULL AND email != '' AND password IS NOT NULL AND password != ''")

	// Apply domain filter if specified
	if targetsDomain != "" {
		query = query.Where("email LIKE ?", "%@"+targetsDomain)
	}

	err := query.Find(&users).Error
	if err != nil {
		zap.L().Error("get_external_credentials",
			zap.String("message", "failed to query external credentials"),
			zap.Error(err),
		)
		return nil, err
	}

	return users, nil
}

// getInternalCredentials retrieves credentials for internal format (username:password)
func getInternalCredentials() ([]sqlite.User, error) {
	db := sqlite.GetDB()
	var users []sqlite.User

	query := db.Where("username IS NOT NULL AND username != '' AND password IS NOT NULL AND password != ''")

	// Apply domain filter if specified (filter usernames that might contain domain info)
	if targetsDomain != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+targetsDomain+"%", "%@"+targetsDomain)
	}

	err := query.Find(&users).Error
	if err != nil {
		zap.L().Error("get_internal_credentials",
			zap.String("message", "failed to query internal credentials"),
			zap.Error(err),
		)
		return nil, err
	}

	return users, nil
}

// getEmailsOnly retrieves emails only (no passwords required)
func getEmailsOnly() ([]sqlite.User, error) {
	db := sqlite.GetDB()
	var users []sqlite.User

	query := db.Where("email IS NOT NULL AND email != ''")

	// Apply domain filter if specified
	if targetsDomain != "" {
		query = query.Where("email LIKE ?", "%@"+targetsDomain)
	}

	err := query.Find(&users).Error
	if err != nil {
		zap.L().Error("get_emails_only",
			zap.String("message", "failed to query emails"),
			zap.Error(err),
		)
		return nil, err
	}

	return users, nil
}

// getSubdomains retrieves subdomains from the database
func getSubdomains() ([]sqlite.Subdomain, error) {
	db := sqlite.GetDB()
	var subdomains []sqlite.Subdomain

	query := db.Where("subdomain IS NOT NULL AND subdomain != ''")

	// Apply domain filter if specified
	if targetsDomain != "" {
		query = query.Where("domain = ? OR subdomain LIKE ?", targetsDomain, "%."+targetsDomain)
	}

	err := query.Find(&subdomains).Error
	if err != nil {
		zap.L().Error("get_subdomains",
			zap.String("message", "failed to query subdomains"),
			zap.Error(err),
		)
		return nil, err
	}

	return subdomains, nil
}
