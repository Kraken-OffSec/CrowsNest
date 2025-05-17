package hunter_io

import (
	"dehasher/internal/debug"
	"dehasher/internal/sqlite"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

const (
	DOMAIN_SEARCH       = "https://api.hunter.io/v2/domain-search?domain={{domain}}&api_key={{apikey}}"
	EMAIL_FINDER        = "https://api.hunter.io/v2/email-finder?domain={{domain}}&first_name={{first_name}}&last_name={{last_name}}&api_key={{apikey}}"
	EMAIL_VERIFICATION  = "https://api.hunter.io/v2/email-verifier?email={{email}}&api_key={{apikey}}"
	COMPANY_ENRICHMENT  = "https://api.hunter.io/v2/companies/find?domain={{domain}}&api_key={{apikey}}"
	PERSON_ENRICHMENT   = "https://api.hunter.io/v2/people/find?email={{email}}&api_key={{apikey}}"
	COMBINED_ENRICHMENT = "https://api.hunter.io/v2/combined/find?email={{email}}&api_key={{apikey}}"
)

type HunterIO struct {
	apiKey string
	debug  bool
}

func NewHunterIO(apiKey string, debugEnabled bool) *HunterIO {
	return &HunterIO{apiKey: apiKey, debug: debugEnabled}
}

func (h *HunterIO) DomainSearch(domain string) (sqlite.HunterDomainData, error) {
	var hunterDomainData sqlite.HunterDomainData

	if h.debug {
		debug.PrintInfo("performing domain search")
		zap.L().Info("hunter_domain_search_debug",
			zap.String("message", "performing domain search"),
		)
	}

	url := DOMAIN_SEARCH
	url = strings.Replace(url, "{{domain}}", domain, -1)
	url = strings.Replace(url, "{{apikey}}", h.apiKey, -1)

	if h.debug {
		debug.PrintInfo("performing request")
		debug.PrintInfo(fmt.Sprintf("URL: %s\n", url))
		zap.L().Info("hunter_domain_search_debug",
			zap.String("message", "performing request"),
			zap.String("url", url),
		)
	}

	resp, err := http.Get(url)
	if err != nil {
		if h.debug {
			debug.PrintInfo("failed to perform request")
			debug.PrintError(err)
		}
		zap.L().Error("hunter_domain_search",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return hunterDomainData, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		if h.debug {
			debug.PrintInfo("failed to read response body")
			debug.PrintError(err)
		}
		zap.L().Error("hunter_domain_search",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return hunterDomainData, err
	}

	if resp.StatusCode != 200 {
		if h.debug {
			debug.PrintInfo("received error status code")
			debug.PrintJson(fmt.Sprintf("Status Code: %d\n", resp.StatusCode))
			debug.PrintJson(fmt.Sprintf("Body: %s\n", string(b)))
			zap.L().Info("hunter_domain_search_debug",
				zap.String("message", "received error status code"),
				zap.Int("status_code", resp.StatusCode),
				zap.String("body_error", string(b)),
			)

		}
		zap.L().Error("hunter_domain_search",
			zap.String("message", "received error status code"),
			zap.Int("status_code", resp.StatusCode),
			zap.String("body_error", string(b)),
		)
		return hunterDomainData, fmt.Errorf("received error status code: %d", resp.StatusCode)
	}

	if h.debug {
		debug.PrintInfo("unmarshalled response body")
		debug.PrintJson(fmt.Sprintf("Body: %s\n", string(b)))
		zap.L().Info("hunter_domain_search_debug",
			zap.String("message", "unmarshalled response body"),
			zap.String("body", string(b)),
		)
	}

	var hunterDomainSearchResult sqlite.HunterDomainSearchResult
	err = json.Unmarshal(b, &hunterDomainSearchResult)
	if err != nil {
		if h.debug {
			debug.PrintInfo("failed to unmarshal response body")
			debug.PrintError(err)
		}
		zap.L().Error("hunter_domain_search",
			zap.String("message", "failed to unmarshal response body"),
			zap.Error(err),
		)
		return hunterDomainData, err
	}

	hunterDomainData = hunterDomainSearchResult.Data

	// Create a list of email object associated with the domain
	var emails []sqlite.HunterEmail
	for _, email := range hunterDomainData.Emails {
		emails = append(emails, sqlite.HunterEmail{
			Domain:       domain,
			Value:        email.Value,
			Type:         email.Type,
			Confidence:   email.Confidence,
			Sources:      email.Sources,
			FirstName:    email.FirstName,
			LastName:     email.LastName,
			Position:     email.Position,
			PositionRaw:  email.PositionRaw,
			Seniority:    email.Seniority,
			Department:   email.Department,
			Linkedin:     email.Linkedin,
			Twitter:      email.Twitter,
			PhoneNumber:  email.PhoneNumber,
			Verification: email.Verification,
		})
	}

	err = sqlite.StoreHunterEmails(emails)
	if err != nil {
		if h.debug {
			debug.PrintInfo("failed to store hunter emails")
			debug.PrintError(err)
		}
		zap.L().Error("store_hunter_emails",
			zap.String("message", "failed to store hunter emails"),
			zap.Error(err),
		)
		return hunterDomainData, err
	}

	err = sqlite.StoreHunterDomain(hunterDomainData)
	if err != nil {
		if h.debug {
			debug.PrintInfo("failed to store hunter domain")
			debug.PrintError(err)
		}
		zap.L().Error("store_hunter_domain",
			zap.String("message", "failed to store hunter domain"),
			zap.Error(err),
		)
		return hunterDomainData, err
	}

	return hunterDomainData, nil
}

func (h *HunterIO) EmailFinder(domain, firstName, lastName string) (sqlite.HunterEmailFinderData, error) {
	var hunterEmailFinderData sqlite.HunterEmailFinderData

	if h.debug {
		debug.PrintInfo("performing email find")
		zap.L().Info("hunter_email_find_debug",
			zap.String("message", "performing email find"),
		)
	}

	url := EMAIL_FINDER
	url = strings.Replace(url, "{{domain}}", domain, -1)
	url = strings.Replace(url, "{{first_name}}", firstName, -1)
	url = strings.Replace(url, "{{last_name}}", lastName, -1)
	url = strings.Replace(url, "{{apikey}}", h.apiKey, -1)

	if h.debug {
		debug.PrintInfo("performing request")
		debug.PrintInfo(fmt.Sprintf("URL: %s\n", url))
		zap.L().Info("hunter_email_find_debug",
			zap.String("message", "performing request"),
			zap.String("url", url),
		)
	}

	resp, err := http.Get(url)
	if err != nil {
		if h.debug {
			debug.PrintInfo("failed to perform request")
			debug.PrintError(err)
		}
		zap.L().Error("hunter_email_find",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return hunterEmailFinderData, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		if h.debug {
			debug.PrintInfo("failed to read response body")
			debug.PrintError(err)
		}
		zap.L().Error("hunter_email_find",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return hunterEmailFinderData, err
	}

	if resp.StatusCode != 200 {
		if h.debug {
			debug.PrintInfo("received error status code")
			debug.PrintJson(fmt.Sprintf("Status Code: %d\n", resp.StatusCode))
			debug.PrintJson(fmt.Sprintf("Body: %s\n", string(b)))
			zap.L().Info("hunter_email_find_debug",
				zap.String("message", "received error status code"),
				zap.Int("status_code", resp.StatusCode),
				zap.String("body_error", string(b)),
			)
		}
		zap.L().Error("hunter_email_find",
			zap.String("message", "received error status code"),
			zap.Int("status_code", resp.StatusCode),
			zap.String("body_error", string(b)),
		)
		return hunterEmailFinderData, fmt.Errorf("received error status code: %d", resp.StatusCode)
	}

	if h.debug {
		debug.PrintInfo("unmarshalled response body")
		debug.PrintJson(fmt.Sprintf("Body: %s\n", string(b)))
		zap.L().Info("hunter_email_find_debug",
			zap.String("message", "unmarshalled response body"),
			zap.String("body", string(b)),
		)
	}

	var hunterEmailFinderResult sqlite.HunterEmailFinderResponse
	err = json.Unmarshal(b, &hunterEmailFinderResult)
	if err != nil {
		if h.debug {
			debug.PrintInfo("failed to unmarshal response body")
			debug.PrintError(err)
		}
		zap.L().Error("hunter_email_find",
			zap.String("message", "failed to unmarshal response body"),
			zap.Error(err),
		)
		return hunterEmailFinderData, err
	}

	hunterEmailFinderData = hunterEmailFinderResult.Data

	var hunterEmails []sqlite.HunterEmail

	hunterEmails = append(hunterEmails, sqlite.HunterEmail{
		Domain:       hunterEmailFinderData.Domain,
		Value:        hunterEmailFinderData.Email,
		Type:         "personal",
		Confidence:   100,
		Sources:      hunterEmailFinderData.Sources,
		FirstName:    hunterEmailFinderData.FirstName,
		LastName:     hunterEmailFinderData.LastName,
		Position:     hunterEmailFinderData.Position,
		PositionRaw:  "",
		Seniority:    "",
		Department:   "",
		Linkedin:     hunterEmailFinderData.LinkedinURL,
		Twitter:      hunterEmailFinderData.Twitter,
		PhoneNumber:  hunterEmailFinderData.PhoneNumber,
		Verification: hunterEmailFinderData.Verification,
	})

	err = sqlite.StoreHunterEmails(hunterEmails)
	if err != nil {
		if h.debug {
			debug.PrintInfo("failed to store hunter email finder")
			debug.PrintError(err)
		}
		zap.L().Error("store_hunter_email_finder",
			zap.String("message", "failed to store hunter email finder"),
			zap.Error(err),
		)
		return hunterEmailFinderData, err
	}

	return hunterEmailFinderData, nil
}

func (h *HunterIO) EmailVerification(email string) (sqlite.HunterEmailVerifyData, error) {
	var hunterEmailVerifyData sqlite.HunterEmailVerifyData

	if h.debug {
		debug.PrintInfo("performing email verification")
		zap.L().Info("hunter_email_verification_debug",
			zap.String("message", "performing email verification"),
		)
	}

	url := EMAIL_VERIFICATION
	url = strings.Replace(url, "{{email}}", email, -1)
	url = strings.Replace(url, "{{apikey}}", h.apiKey, -1)

	if h.debug {
		debug.PrintInfo("performing request")
		debug.PrintInfo(fmt.Sprintf("URL: %s\n", url))
		zap.L().Info("hunter_email_verification_debug",
			zap.String("message", "performing request"),
			zap.String("url", url),
		)
	}

	resp, err := http.Get(url)
	if err != nil {
		if h.debug {
			debug.PrintInfo("failed to perform request")
			debug.PrintError(err)
		}
		zap.L().Error("hunter_email_verification",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return hunterEmailVerifyData, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		if h.debug {
			debug.PrintInfo("failed to read response body")
			debug.PrintError(err)
		}
		zap.L().Error("hunter_email_verification",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return hunterEmailVerifyData, err
	}

	if resp.StatusCode != 200 {
		if h.debug {
			debug.PrintInfo("received error status code")
			debug.PrintJson(fmt.Sprintf("Status Code: %d\n", resp.StatusCode))
			debug.PrintJson(fmt.Sprintf("Body: %s\n", string(b)))
			zap.L().Info("hunter_email_verification_debug",
				zap.String("message", "received error status code"),
				zap.Int("status_code", resp.StatusCode),
				zap.String("body_error", string(b)),
			)
		}
		zap.L().Error("hunter_email_verification",
			zap.String("message", "received error status code"),
			zap.Int("status_code", resp.StatusCode),
			zap.String("body_error", string(b)),
		)
		return hunterEmailVerifyData, fmt.Errorf("received error status code: %d", resp.StatusCode)
	}

	if h.debug {
		debug.PrintInfo("unmarshalled response body")
		debug.PrintJson(fmt.Sprintf("Body: %s\n", string(b)))
		zap.L().Info("hunter_email_verification_debug",
			zap.String("message", "unmarshalled response body"),
			zap.String("body", string(b)),
		)
	}

	var hunterEmailVerifyResult sqlite.HunterEmailVerifyResponse
	err = json.Unmarshal(b, &hunterEmailVerifyResult)
	if err != nil {
		if h.debug {
			debug.PrintInfo("failed to unmarshal response body")
			debug.PrintError(err)
		}
		zap.L().Error("hunter_email_verification",
			zap.String("message", "failed to unmarshal response body"),
			zap.Error(err),
		)
		return hunterEmailVerifyData, err
	}

	hunterEmailVerifyData = hunterEmailVerifyResult.Data

	return hunterEmailVerifyData, nil
}

func (h *HunterIO) CompanyEnrichment(domain string) (sqlite.CompanyData, error) {
	var companyData sqlite.CompanyData

	if h.debug {
		debug.PrintInfo("performing company enrichment")
		zap.L().Info("hunter_company_enrichment_debug",
			zap.String("message", "performing company enrichment"),
		)
	}

	url := COMPANY_ENRICHMENT
	url = strings.Replace(url, "{{domain}}", domain, -1)
	url = strings.Replace(url, "{{apikey}}", h.apiKey, -1)

	if h.debug {
		debug.PrintInfo("performing request")
		debug.PrintInfo(fmt.Sprintf("URL: %s\n", url))
		zap.L().Info("hunter_company_enrichment_debug",
			zap.String("message", "performing request"),
			zap.String("url", url),
		)
	}

	resp, err := http.Get(url)
	if err != nil {
		if h.debug {
			debug.PrintInfo("failed to perform request")
			debug.PrintError(err)
		}
		zap.L().Error("hunter_company_enrichment",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return companyData, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		if h.debug {
			debug.PrintInfo("failed to read response body")
			debug.PrintError(err)
		}
		zap.L().Error("hunter_company_enrichment",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return companyData, err
	}

	if resp.StatusCode != 200 {
		if h.debug {
			debug.PrintInfo("received error status code")
			debug.PrintJson(fmt.Sprintf("Status Code: %d\n", resp.StatusCode))
			debug.PrintJson(fmt.Sprintf("Body: %s\n", string(b)))
			zap.L().Info("hunter_company_enrichment_debug",
				zap.String("message", "received error status code"),
				zap.Int("status_code", resp.StatusCode),
				zap.String("body_error", string(b)),
			)
		}
		zap.L().Error("hunter_company_enrichment",
			zap.String("message", "received error status code"),
			zap.Int("status_code", resp.StatusCode),
			zap.String("body_error", string(b)),
		)
		return companyData, fmt.Errorf("received error status code: %d", resp.StatusCode)
	}

	if h.debug {
		debug.PrintInfo("unmarshalled response body")
		debug.PrintJson(fmt.Sprintf("Body: %s\n", string(b)))
		zap.L().Info("hunter_company_enrichment_debug",
			zap.String("message", "unmarshalled response body"),
			zap.String("body", string(b)),
		)
	}

	var hunterCompanyEnrichmentResult sqlite.HunterCompanyEnrichmentResponse
	err = json.Unmarshal(b, &hunterCompanyEnrichmentResult)
	if err != nil {
		if h.debug {
			debug.PrintInfo("failed to unmarshal response body")
			debug.PrintError(err)
		}
		zap.L().Error("hunter_company_enrichment",
			zap.String("message", "failed to unmarshal response body"),
			zap.Error(err),
		)
		return companyData, err
	}

	companyData = hunterCompanyEnrichmentResult.Data

	return companyData, nil
}

func (h *HunterIO) PersonEnrichment(email string) (sqlite.PersonData, error) {
	var personData sqlite.PersonData

	if h.debug {
		debug.PrintInfo("performing person enrichment")
		zap.L().Info("hunter_person_enrichment_debug",
			zap.String("message", "performing person enrichment"),
		)
	}

	url := PERSON_ENRICHMENT
	url = strings.Replace(url, "{{email}}", email, -1)
	url = strings.Replace(url, "{{apikey}}", h.apiKey, -1)

	if h.debug {
		debug.PrintInfo("performing request")
		debug.PrintInfo(fmt.Sprintf("URL: %s\n", url))
		zap.L().Info("hunter_person_enrichment_debug",
			zap.String("message", "performing request"),
			zap.String("url", url),
		)
	}

	resp, err := http.Get(url)
	if err != nil {
		if h.debug {
			debug.PrintInfo("failed to perform request")
			debug.PrintError(err)
		}
		zap.L().Error("hunter_person_enrichment",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return personData, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		if h.debug {
			debug.PrintInfo("failed to read response body")
			debug.PrintError(err)
		}
		zap.L().Error("hunter_person_enrichment",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return personData, err
	}

	if resp.StatusCode != 200 {
		if h.debug {
			debug.PrintInfo("received error status code")
			debug.PrintJson(fmt.Sprintf("Status Code: %d\n", resp.StatusCode))
			debug.PrintJson(fmt.Sprintf("Body: %s\n", string(b)))
			zap.L().Info("hunter_person_enrichment_debug",
				zap.String("message", "received error status code"),
				zap.Int("status_code", resp.StatusCode),
				zap.String("body_error", string(b)),
			)
		}
		zap.L().Error("hunter_person_enrichment",
			zap.String("message", "received error status code"),
			zap.Int("status_code", resp.StatusCode),
			zap.String("body_error", string(b)),
		)
		return personData, fmt.Errorf("received error status code: %d", resp.StatusCode)
	}

	if h.debug {
		debug.PrintInfo("unmarshalled response body")
		debug.PrintJson(fmt.Sprintf("Body: %s\n", string(b)))
		zap.L().Info("hunter_person_enrichment_debug",
			zap.String("message", "unmarshalled response body"),
			zap.String("body", string(b)),
		)
	}

	var hunterPersonEnrichmentResult sqlite.HunterPersonEnrichmentResponse
	err = json.Unmarshal(b, &hunterPersonEnrichmentResult)
	if err != nil {
		if h.debug {
			debug.PrintInfo("failed to unmarshal response body")
			debug.PrintError(err)
		}
		zap.L().Error("hunter_person_enrichment",
			zap.String("message", "failed to unmarshal response body"),
			zap.Error(err),
		)
		return personData, err
	}

	personData = hunterPersonEnrichmentResult.Data

	return personData, nil
}

func (h *HunterIO) CombinedEnrichment(email string) (sqlite.CombinedData, error) {
	var combinedData sqlite.CombinedData

	if h.debug {
		debug.PrintInfo("performing combined enrichment")
		zap.L().Info("hunter_combined_enrichment_debug",
			zap.String("message", "performing combined enrichment"),
		)
	}

	url := COMBINED_ENRICHMENT
	url = strings.Replace(url, "{{email}}", email, -1)
	url = strings.Replace(url, "{{apikey}}", h.apiKey, -1)

	if h.debug {
		debug.PrintInfo("performing request")
		debug.PrintInfo(fmt.Sprintf("URL: %s\n", url))
		zap.L().Info("hunter_combined_enrichment_debug",
			zap.String("message", "performing request"),
			zap.String("url", url),
		)
	}

	resp, err := http.Get(url)
	if err != nil {
		if h.debug {
			debug.PrintInfo("failed to perform request")
			debug.PrintError(err)
		}
		zap.L().Error("hunter_combined_enrichment",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return combinedData, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		if h.debug {
			debug.PrintInfo("failed to read response body")
			debug.PrintError(err)
		}
		zap.L().Error("hunter_combined_enrichment",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return combinedData, err
	}

	if resp.StatusCode != 200 {
		if h.debug {
			debug.PrintInfo("received error status code")
			debug.PrintJson(fmt.Sprintf("Status Code: %d\n", resp.StatusCode))
			debug.PrintJson(fmt.Sprintf("Body: %s\n", string(b)))
			zap.L().Info("hunter_combined_enrichment_debug",
				zap.String("message", "received error status code"),
				zap.Int("status_code", resp.StatusCode),
				zap.String("body_error", string(b)),
			)
		}
		zap.L().Error("hunter_combined_enrichment",
			zap.String("message", "received error status code"),
			zap.Int("status_code", resp.StatusCode),
			zap.String("body_error", string(b)),
		)
		return combinedData, fmt.Errorf("received error status code: %d", resp.StatusCode)
	}

	if h.debug {
		debug.PrintInfo("unmarshalled response body")
		debug.PrintJson(fmt.Sprintf("Body: %s\n", string(b)))
		zap.L().Info("hunter_combined_enrichment_debug",
			zap.String("message", "unmarshalled response body"),
			zap.String("body", string(b)),
		)
	}

	var hunterCombinedEnrichmentResult sqlite.HunterCombinedEnrichmentResponse
	err = json.Unmarshal(b, &hunterCombinedEnrichmentResult)
	if err != nil {
		if h.debug {
			debug.PrintInfo("failed to unmarshal response body")
			debug.PrintError(err)
		}
		zap.L().Error("hunter_combined_enrichment",
			zap.String("message", "failed to unmarshal response body"),
			zap.Error(err),
		)
		return combinedData, err
	}

	combinedData = hunterCombinedEnrichmentResult.Data

	return combinedData, nil
}
