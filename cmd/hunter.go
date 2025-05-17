package cmd

import (
	"crowsnest/internal/badger"
	"crowsnest/internal/debug"
	"crowsnest/internal/export"
	"crowsnest/internal/files"
	hunter "crowsnest/internal/hunter.io"
	"crowsnest/internal/pretty"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"time"
)

func init() {
	// Add hunter command to root command
	rootCmd.AddCommand(hunterCmd)

	// Add flags specific to hunter command
	hunterCmd.Flags().StringVarP(&hunterDomain, "domain", "d", "", "Domain to query")
	hunterCmd.Flags().StringVarP(&hunterEmail, "email", "e", "", "Email to query")
	hunterCmd.Flags().StringVarP(&hunterFirstName, "first-name", "F", "", "First name to query")
	hunterCmd.Flags().StringVarP(&hunterLastName, "last-name", "L", "", "Last name to query")
	hunterCmd.Flags().BoolVarP(&hunterDomainSearch, "domain-search", "D", false, "Search for domain")
	hunterCmd.Flags().BoolVarP(&hunterEmailFind, "email-find", "E", false, "Find emails for user using domain, first name, and last name")
	hunterCmd.Flags().BoolVarP(&hunterEmailVerify, "email-verify", "V", false, "Verify email")
	hunterCmd.Flags().BoolVarP(&hunterCompanyEnrichmentDomain, "company-enrichment", "C", false, "Company enrichment for domain")
	hunterCmd.Flags().BoolVarP(&hunterPersonEnrichmentEmail, "person-enrichment", "P", false, "Person enrichment for email")
	hunterCmd.Flags().BoolVarP(&hunterCombinedEnrichmentEmail, "combined-enrichment", "B", false, "Combined Company and Person enrichment for email")
	hunterCmd.Flags().StringVarP(&hunterOutputFormat, "format", "f", "json", "Output format (json, yaml, xml, txt)")
	hunterCmd.Flags().StringVarP(&hunterOutputFile, "output", "o", "hunter", "File to output results to including extension")

	// Add mutually exclusive flags to hunter command
	hunterCmd.MarkFlagsMutuallyExclusive("email-find")

}

var (
	// Hunter Commands Flags
	hunterDomain                  string
	hunterEmail                   string
	hunterFirstName               string
	hunterLastName                string
	hunterDomainSearch            bool
	hunterEmailFind               bool
	hunterEmailVerify             bool
	hunterCompanyEnrichmentDomain bool
	hunterPersonEnrichmentEmail   bool
	hunterCombinedEnrichmentEmail bool
	hunterOutputFormat            string
	hunterOutputFile              string

	hunterCmd = &cobra.Command{
		Use:   "hunter",
		Short: "Hunter.io API interaction",
		Long:  `Interact with the Hunter.io API for email and domain information.`,
		Run: func(cmd *cobra.Command, args []string) {
			if debugGlobal {
				debug.PrintInfo("debug mode enabled")
				zap.L().Info("hunter_debug",
					zap.String("message", "debug mode enabled"),
				)
			}

			// Flag Checks
			if !hunterFlagCheck() {
				return
			}

			if hunterOutputFile == "" {
				if debugGlobal {
					debug.PrintInfo("output file not specified, using default")
				}
				hunterOutputFile = "hunter_" + time.Now().Format("05_04_05")
			}

			if hunterOutputFormat == "" {
				if debugGlobal {
					debug.PrintInfo("output format not specified, using default")
				}
				hunterOutputFormat = "json"
			}

			fType := files.GetFileType(hunterOutputFormat)
			if fType == files.UNKNOWN {
				fmt.Println("[!] Error: Invalid output format. Must be 'json', 'xml', 'yaml', or 'txt'.")
				return
			}
			if debugGlobal {
				debug.PrintInfo("using output format: " + hunterOutputFormat)
			}

			fmt.Println("[*] Hunter.io API interaction [Beta]")

			h := hunter.NewHunterIO(getHunterApiKey(), debugGlobal)

			if hunterDomainSearch {
				fmt.Println("[*] Performing domain search search...")
				result, err := h.DomainSearch(hunterDomain)
				if err != nil {
					if debugGlobal {
						debug.PrintInfo("failed to perform domain search")
						debug.PrintError(err)
					}
					zap.L().Error("hunter_domain_search",
						zap.String("message", "failed to perform domain search"),
						zap.Error(err),
					)
					fmt.Printf("Error performing domain search: %v\n", err)
					return
				}

				// Write Hunter.io Domain Search Result to file
				fmt.Printf("[*] Writing Hunter.io Domain Search Result to file: %s%s\n", hunterOutputFile, fType.Extension())
				err = export.WriteIStringToFile(result, hunterOutputFile, fType)
				if err != nil {
					if debugGlobal {
						debug.PrintInfo("failed to write hunter domain search to file")
						debug.PrintError(err)
					}
					zap.L().Error("write_hunter_domain_search",
						zap.String("message", "failed to write hunter domain search to file"),
						zap.Error(err),
					)
					fmt.Printf("Error writing Hunter.io Domain Search Result to file: %v\n", err)
				}

				// Pretty Print Hunter.io Domain Search Result
				fmt.Println("Domain Search Result:")
				pretty.HunterDomainTree(hunterDomain, result)
				return
			}

			if hunterEmailFind {
				fmt.Println("[*] Performing email find search...")
				result, err := h.EmailFinder(hunterDomain, hunterFirstName, hunterLastName)
				if err != nil {
					if debugGlobal {
						debug.PrintInfo("failed to perform email find")
						debug.PrintError(err)
					}
					zap.L().Error("hunter_email_find",
						zap.String("message", "failed to perform email find"),
						zap.Error(err),
					)
					fmt.Printf("Error performing email find: %v\n", err)
					return
				}

				// Write Hunter.io Email Finder Result to file
				fmt.Printf("[*] Writing Hunter.io Email Finder Result to file: %s%s\n", hunterOutputFile, fType.Extension())
				err = export.WriteIStringToFile(result, hunterOutputFile, fType)
				if err != nil {
					if debugGlobal {
						debug.PrintInfo("failed to write hunter email find to file")
						debug.PrintError(err)
					}
					zap.L().Error("write_hunter_email_find",
						zap.String("message", "failed to write hunter email find to file"),
						zap.Error(err),
					)
					fmt.Printf("Error writing Hunter.io Email Finder Result to file: %v\n", err)
				}

				fmt.Println("Email Find Result:")

				var (
					headers = []string{"Email", "Score", "Domain", "Accept All", "Position", "Twitter", "Linkedin", "Phone Number", "Company", "Verification"}
					rows    [][]string
				)

				rows = append(rows, []string{
					result.Email,
					fmt.Sprintf("%d", result.Score),
					result.Domain,
					fmt.Sprintf("%t", result.AcceptAll),
					result.Position,
					result.Twitter,
					result.LinkedinURL,
					result.PhoneNumber,
					result.Company,
					fmt.Sprintf("%v", result.Verification),
				})

				pretty.Table(headers, rows)

				return
			}

			if hunterEmailVerify {
				fmt.Println("[*] Performing email verification search...")
				result, err := h.EmailVerification(hunterEmail)
				if err != nil {
					if debugGlobal {
						debug.PrintInfo("failed to perform email verification")
						debug.PrintError(err)
					}
					zap.L().Error("hunter_email_verification",
						zap.String("message", "failed to perform email verification"),
						zap.Error(err),
					)
					fmt.Printf("Error performing email verification: %v\n", err)
					return
				}
				// Write Hunter.io Email Verification Result to file
				fmt.Printf("[*] Writing Hunter.io Email Verification Result to file: %s%s\n", hunterOutputFile, fType.Extension())
				err = export.WriteIStringToFile(result, hunterOutputFile, fType)
				if err != nil {
					if debugGlobal {
						debug.PrintInfo("failed to write hunter email verification to file")
						debug.PrintError(err)
					}
					zap.L().Error("write_hunter_email_verification",
						zap.String("message", "failed to write hunter email verification to file"),
						zap.Error(err),
					)
					fmt.Printf("Error writing Hunter.io Email Verification Result to file: %v\n", err)
				}

				// Pretty Print Hunter.io Email Verification Result
				var (
					headers = []string{"Email", "Status", "Result", "Score", "Regexp", "Gibberish", "Disposable", "Webmail", "MX Records", "SMTP Server", "SMTP Check", "Accept All", "Block", "Sources"}
					rows    [][]string
				)
				rows = append(rows, []string{
					result.Email,
					result.Status,
					result.Result,
					fmt.Sprintf("%d", result.Score),
					fmt.Sprintf("%t", result.Regexp),
					fmt.Sprintf("%t", result.Gibberish),
					fmt.Sprintf("%t", result.Disposable),
					fmt.Sprintf("%t", result.Webmail),
					fmt.Sprintf("%t", result.MXRecords),
					fmt.Sprintf("%t", result.SMTPServer),
					fmt.Sprintf("%t", result.SMTPCheck),
					fmt.Sprintf("%t", result.AcceptAll),
					fmt.Sprintf("%t", result.Block),
					fmt.Sprintf("%v", result.Sources),
				})

				fmt.Println("Email Verification Result:")
				pretty.Table(headers, rows)
				return
			}

			if hunterCompanyEnrichmentDomain {
				fmt.Println("[*] Performing company enrichment search...")
				result, err := h.CompanyEnrichment(hunterDomain)
				if err != nil {
					if debugGlobal {
						debug.PrintInfo("failed to perform company enrichment")
						debug.PrintError(err)
					}
					zap.L().Error("hunter_company_enrichment",
						zap.String("message", "failed to perform company enrichment"),
						zap.Error(err),
					)
					fmt.Printf("Error performing company enrichment: %v\n", err)
					return
				}

				// Write to file
				fmt.Printf("[*] Writing Hunter.io Company Enrichment Result to file: %s%s\n", hunterOutputFile, fType.Extension())
				err = export.WriteIStringToFile(result, hunterOutputFile, fType)
				if err != nil {
					if debugGlobal {
						debug.PrintInfo("failed to write hunter company enrichment to file")
						debug.PrintError(err)
					}
					zap.L().Error("write_hunter_company_enrichment",
						zap.String("message", "failed to write hunter company enrichment to file"),
						zap.Error(err),
					)
					fmt.Printf("Error writing Hunter.io Company Enrichment Result to file: %v\n", err)
				}

				// Pretty Print Hunter.io Company Enrichment Result
				fmt.Println("Company Enrichment Result:")
				pretty.HunterCompanyEnrichmentTree(hunterDomain, result)

				return
			}

			if hunterPersonEnrichmentEmail {
				fmt.Println("[*] Performing person enrichment search...")
				result, err := h.PersonEnrichment(hunterEmail)
				if err != nil {
					if debugGlobal {
						debug.PrintInfo("failed to perform person enrichment")
						debug.PrintError(err)
					}
					zap.L().Error("hunter_person_enrichment",
						zap.String("message", "failed to perform person enrichment"),
						zap.Error(err),
					)
					fmt.Printf("Error performing person enrichment: %v\n", err)
					return
				}

				// Write to file
				fmt.Printf("[*] Writing Hunter.io Person Enrichment Result to file: %s%s\n", hunterOutputFile, fType.Extension())
				err = export.WriteIStringToFile(result, hunterOutputFile, fType)
				if err != nil {
					if debugGlobal {
						debug.PrintInfo("failed to write hunter person enrichment to file")
						debug.PrintError(err)
					}
					zap.L().Error("write_hunter_person_enrichment",
						zap.String("message", "failed to write hunter person enrichment to file"),
						zap.Error(err),
					)
					fmt.Printf("Error writing Hunter.io Person Enrichment Result to file: %v\n", err)
				}

				// Pretty Print Hunter.io Person Enrichment Result
				fmt.Println("Person Enrichment Result:")
				pretty.HunterPersonEnrichmentTree(hunterEmail, result)
				return
			}

			if hunterCombinedEnrichmentEmail {
				fmt.Println("[*] Performing combined enrichment search...")
				result, err := h.CombinedEnrichment(hunterEmail)
				if err != nil {
					if debugGlobal {
						debug.PrintInfo("failed to perform combined enrichment")
						debug.PrintError(err)
					}
					zap.L().Error("hunter_combined_enrichment",
						zap.String("message", "failed to perform combined enrichment"),
						zap.Error(err),
					)
					fmt.Printf("Error performing combined enrichment: %v\n", err)
					return
				}

				// Write to file
				fmt.Printf("[*] Writing Hunter.io Combined Enrichment Result to file: %s%s\n", hunterOutputFile, fType.Extension())
				err = export.WriteIStringToFile(result, hunterOutputFile, fType)
				if err != nil {
					if debugGlobal {
						debug.PrintInfo("failed to write hunter combined enrichment to file")
						debug.PrintError(err)
					}
					zap.L().Error("write_hunter_combined_enrichment",
						zap.String("message", "failed to write hunter combined enrichment to file"),
						zap.Error(err),
					)
					fmt.Printf("Error writing Hunter.io Combined Enrichment Result to file: %v\n", err)
				}

				fmt.Println("Combined Enrichment Result:")
				pretty.HunterCombinedEnrichmentTree(hunterEmail, result)
				return
			}

		},
	}
)

func hunterFlagCheck() bool {
	if debugGlobal {
		debug.PrintInfo("checking flags")
	}

	var optionSet bool

	if hunterDomainSearch {
		if hunterDomain == "" {
			fmt.Println("Domain is required for domain search")
			return false
		}
		optionSet = true
	}
	if hunterEmailVerify {
		if hunterEmail == "" {
			fmt.Println("Email is required for email verification")
			return false
		}
		optionSet = true
	}
	if hunterCompanyEnrichmentDomain {
		if hunterDomain == "" {
			fmt.Println("Domain is required for company enrichment")
			return false
		}
		optionSet = true
	}
	if hunterPersonEnrichmentEmail {
		if hunterEmail == "" {
			fmt.Println("Email is required for person enrichment")
			return false
		}
		optionSet = true
	}
	if hunterCombinedEnrichmentEmail {
		if hunterEmail == "" {
			fmt.Println("Email is required for combined enrichment")
			return false
		}
		optionSet = true
	}
	if hunterEmailFind {
		if hunterFirstName == "" || hunterLastName == "" {
			fmt.Println("First name and last name are required for email find")
			return false
		}
		if hunterDomain == "" {
			fmt.Println("Domain is required for email find")
			return false
		}
		optionSet = true
	}

	if !optionSet {
		fmt.Println("[!] No options selected")
		return false
	}

	return true
}

// Helper functions to get stored API credentials
func getHunterApiKey() string {
	return badger.GetHunterKey()
}
