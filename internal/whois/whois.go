package whois

import (
	"bytes"
	"dehasher/internal/query"
	"dehasher/internal/sqlite"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type DehashedWHOISSearchRequest struct {
	Include     []string `json:"include,omitempty"`
	Exclude     []string `json:"exclude,omitempty"`
	IPAddress   string   `json:"ip_address,omitempty"`
	ReverseType string   `json:"reverse_type,omitempty"`
	Domain      string   `json:"domain,omitempty"`
	MXAddress   string   `json:"mx_address,omitempty"`
	NSAddress   string   `json:"ns_address,omitempty"`
	SearchType  string   `json:"search_type,omitempty"`
}

func WhoisSearch(domain, apiKey string) (sqlite.WhoIsLookupResult, error) {
	var whois sqlite.WhoIsLookupResult
	whoisSearchRequest := DehashedWHOISSearchRequest{
		Domain:     domain,
		SearchType: "whois",
	}
	reqBody, _ := json.Marshal(whoisSearchRequest)
	req, err := http.NewRequest("POST", "https://api.dehashed.com/v2/whois/search", bytes.NewReader(reqBody))
	if err != nil {
		return whois, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", apiKey)
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		zap.L().Error("whois_search",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return whois, err
	}
	if res == nil {
		zap.L().Error("whois_search",
			zap.String("message", "response was nil"),
		)
		return whois, errors.New("response was nil")
	}

	// Check for HTTP status code errors
	if res.StatusCode != 200 {
		dhErr := query.GetDehashedError(res.StatusCode)
		fmt.Printf("[%d] API Error message: %s\n", res.StatusCode, dhErr.Error())
		zap.L().Error("whois_search",
			zap.String("message", "received error status code"),
			zap.Int("status_code", res.StatusCode),
			zap.String("error", dhErr.Error()),
		)
		return whois, &dhErr
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("whois_search",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return whois, err
	}

	err = json.Unmarshal(b, &whois)
	if err != nil {
		zap.L().Error("whois_search",
			zap.String("message", "failed to unmarshal response body"),
			zap.Error(err),
		)
		fmt.Println("Error unmarshalling response body:", err)
		fmt.Println("Response body:", string(b))
		return whois, err
	}

	return whois, nil
}

func WhoisHistory(domain, apiKey string) (sqlite.WhoIsHistory, error) {
	var whois sqlite.WhoIsHistory
	whoisSearchRequest := DehashedWHOISSearchRequest{
		Domain:     domain,
		SearchType: "whois-history",
	}
	reqBody, _ := json.Marshal(whoisSearchRequest)
	req, err := http.NewRequest("POST", "https://api.dehashed.com/v2/whois/search", bytes.NewReader(reqBody))
	if err != nil {
		return whois, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", apiKey)
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		zap.L().Info("whois_history",
			zap.String("message", "response was not nil"),
		)
		defer res.Body.Close()
	}
	if err != nil {
		zap.L().Error("whois_history",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return whois, err
	}
	if res == nil {
		zap.L().Error("whois_history",
			zap.String("message", "response was nil"),
		)
		return whois, errors.New("response was nil")
	}

	// Check for HTTP status code errors
	if res.StatusCode != 200 {
		dhErr := query.GetDehashedError(res.StatusCode)
		fmt.Printf("[%d] API Error message: %s\n", res.StatusCode, dhErr.Error())
		zap.L().Error("whois_history",
			zap.String("message", "received error status code"),
			zap.Int("status_code", res.StatusCode),
			zap.String("error", dhErr.Error()),
		)
		return whois, &dhErr
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("whois_history",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return whois, err
	}

	err = json.Unmarshal(b, &whois)
	if err != nil {
		zap.L().Error("whois_history",
			zap.String("message", "failed to unmarshal response body"),
			zap.Error(err),
		)
		fmt.Println("Error unmarshalling response body:", err)
		fmt.Println("Response body:", string(b))
		return whois, err
	}

	return whois, nil
}

func ReverseWHOIS(include []string, exclude []string, reverseType, apiKey string) (string, error) {
	whoisSearchRequest := DehashedWHOISSearchRequest{
		Include:     include,
		Exclude:     exclude,
		ReverseType: reverseType,
		SearchType:  "reverse-whois",
	}
	reqBody, _ := json.Marshal(whoisSearchRequest)
	req, err := http.NewRequest("POST", "https://api.dehashed.com/v2/whois/search", bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", apiKey)
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		zap.L().Error("reverse_whois",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return "", err
	}
	if res == nil {
		zap.L().Error("reverse_whois",
			zap.String("message", "response was nil"),
		)
		return "", errors.New("response was nil")
	}

	// Check for HTTP status code errors
	if res.StatusCode != 200 {
		dhErr := query.GetDehashedError(res.StatusCode)
		fmt.Printf("[%d] API Error message: %s\n", res.StatusCode, dhErr.Error())
		zap.L().Error("reverse_whois",
			zap.String("message", "received error status code"),
			zap.Int("status_code", res.StatusCode),
			zap.String("error", dhErr.Error()),
		)
		return "", &dhErr
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("reverse_whois",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return "", err
	}
	return string(b), nil
}

func WhoisIP(ipAddress, apiKey string) ([]byte, error) {
	whoisSearchRequest := DehashedWHOISSearchRequest{
		IPAddress:  ipAddress,
		SearchType: "reverse-ip",
	}
	reqBody, _ := json.Marshal(whoisSearchRequest)
	req, err := http.NewRequest("POST", "https://api.dehashed.com/v2/whois/search", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", apiKey)
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		zap.L().Error("whois_ip",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return nil, err
	}
	if res == nil {
		zap.L().Error("whois_ip",
			zap.String("message", "response was nil"),
		)
		return nil, errors.New("response was nil")
	}

	// Check for HTTP status code errors
	if res.StatusCode != 200 {
		dhErr := query.GetDehashedError(res.StatusCode)
		fmt.Printf("[%d] API Error message: %s\n", res.StatusCode, dhErr.Error())
		zap.L().Error("whois_ip",
			zap.String("message", "received error status code"),
			zap.Int("status_code", res.StatusCode),
			zap.String("error", dhErr.Error()),
		)
		return nil, &dhErr
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("whois_ip",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return nil, err
	}
	return b, nil
}

func WhoisMX(mxAddress, apiKey string) (string, error) {
	whoisSearchRequest := DehashedWHOISSearchRequest{
		MXAddress:  mxAddress,
		SearchType: "reverse-mx",
	}
	reqBody, _ := json.Marshal(whoisSearchRequest)
	req, err := http.NewRequest("POST", "https://api.dehashed.com/v2/whois/search", bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", apiKey)
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		zap.L().Error("whois_mx",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return "", err
	}
	if res == nil {
		zap.L().Error("whois_mx",
			zap.String("message", "response was nil"),
		)
		return "", errors.New("response was nil")
	}

	// Check for HTTP status code errors
	if res.StatusCode != 200 {
		dhErr := query.GetDehashedError(res.StatusCode)
		fmt.Printf("[%d] API Error message: %s\n", res.StatusCode, dhErr.Error())
		zap.L().Error("whois_mx",
			zap.String("message", "received error status code"),
			zap.Int("status_code", res.StatusCode),
			zap.String("error", dhErr.Error()),
		)
		return "", &dhErr
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("whois_mx",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return "", err
	}
	return string(b), nil
}

func WhoisNS(nsAddress, apiKey string) (string, error) {
	whoisSearchRequest := DehashedWHOISSearchRequest{
		NSAddress:  nsAddress,
		SearchType: "reverse-ns",
	}
	reqBody, _ := json.Marshal(whoisSearchRequest)
	req, err := http.NewRequest("POST", "https://api.dehashed.com/v2/whois/search", bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", apiKey)
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		zap.L().Error("whois_ns",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return "", err
	}
	if res == nil {
		zap.L().Error("whois_ns",
			zap.String("message", "response was nil"),
		)
		return "", errors.New("response was nil")
	}

	// Check for HTTP status code errors
	if res.StatusCode != 200 {
		dhErr := query.GetDehashedError(res.StatusCode)
		fmt.Printf("[%d] API Error message: %s\n", res.StatusCode, dhErr.Error())
		zap.L().Error("whois_ns",
			zap.String("message", "received error status code"),
			zap.Int("status_code", res.StatusCode),
			zap.String("error", dhErr.Error()),
		)
		return "", &dhErr
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("whois_ns",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return "", err
	}
	return string(b), nil
}

func WhoisSubdomainScan(domain, apiKey string) (sqlite.WhoIsSubdomainScan, error) {
	var whois sqlite.WhoIsSubdomainScan
	whoisSearchRequest := DehashedWHOISSearchRequest{
		Domain:     domain,
		SearchType: "subdomain-scan",
	}
	reqBody, _ := json.Marshal(whoisSearchRequest)
	req, err := http.NewRequest("POST", "https://api.dehashed.com/v2/whois/search", bytes.NewReader(reqBody))
	if err != nil {
		return whois, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", apiKey)
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		zap.L().Error("whois_subdomain_scan",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return whois, err
	}
	if res == nil {
		zap.L().Error("whois_subdomain_scan",
			zap.String("message", "response was nil"),
		)
		return whois, errors.New("response was nil")
	}

	// Check for HTTP status code errors
	if res.StatusCode != 200 {
		dhErr := query.GetDehashedError(res.StatusCode)
		fmt.Printf("[%d] API Error message: %s\n", res.StatusCode, dhErr.Error())
		zap.L().Error("whois_subdomain_scan",
			zap.String("message", "received error status code"),
			zap.Int("status_code", res.StatusCode),
			zap.String("error", dhErr.Error()),
		)
		return whois, &dhErr
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("whois_subdomain_scan",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return whois, err
	}

	err = json.Unmarshal(b, &whois)
	if err != nil {
		zap.L().Error("whois_subdomain_scan",
			zap.String("message", "failed to unmarshal response body"),
			zap.Error(err),
		)
		fmt.Println("Error unmarshalling response body:", err)
		fmt.Println("Response body:", string(b))
		return whois, err
	}

	return whois, nil
}

func GetWHOISCredits(apiKey string) (sqlite.WhoIsCredits, error) {
	var whoisCredits sqlite.WhoIsCredits

	req, err := http.NewRequest("GET", "https://api.dehashed.com/v2/whois/credits", nil)
	if err != nil {
		return whoisCredits, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", apiKey)
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		zap.L().Error("get_whois_credits",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return whoisCredits, err
	}
	if res == nil {
		zap.L().Error("get_whois_credits",
			zap.String("message", "response was nil"),
		)
		return whoisCredits, errors.New("response was nil")
	}

	// Check for HTTP status code errors
	if res.StatusCode != 200 {
		dhErr := query.GetDehashedError(res.StatusCode)
		fmt.Printf("[%d] API Error message: %s\n", res.StatusCode, dhErr.Error())
		zap.L().Error("get_whois_credits",
			zap.String("message", "received error status code"),
			zap.Int("status_code", res.StatusCode),
			zap.String("error", dhErr.Error()),
		)
		return whoisCredits, &dhErr
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("get_whois_credits",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return whoisCredits, err
	}

	err = json.Unmarshal(b, &whoisCredits)
	if err != nil {
		zap.L().Error("get_whois_credits",
			zap.String("message", "failed to unmarshal response body"),
			zap.Error(err),
		)
		fmt.Println("Error unmarshalling response body:", err)
		fmt.Println("Response body:", string(b))
		return whoisCredits, err
	}

	return whoisCredits, nil
}
