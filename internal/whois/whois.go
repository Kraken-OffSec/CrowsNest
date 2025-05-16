package whois

import (
	"bytes"
	"dehasher/internal/debug"
	"dehasher/internal/dehashed"
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

type DehashedWhoIs struct {
	balance int
	debug   bool
	apiKey  string
}

func NewWhoIs(apiKey string, debug bool) *DehashedWhoIs {
	return &DehashedWhoIs{apiKey: apiKey, debug: debug, balance: -1}
}

func (w *DehashedWhoIs) WhoisSearch(domain string) (sqlite.WhoisRecord, error) {
	var whois sqlite.WhoIsLookupResult
	var whoisRecord sqlite.WhoisRecord

	if w.debug {
		debug.PrintInfo("performing whois search")
		zap.L().Info("whois_search_debug",
			zap.String("message", "performing whois search"),
		)
	}

	whoisSearchRequest := DehashedWHOISSearchRequest{
		Domain:     domain,
		SearchType: "whois",
	}

	if w.debug {
		debug.PrintInfo("building request body")
		debug.PrintJson(fmt.Sprintf("Request Body: %v\n", whoisSearchRequest))
		zap.L().Info("whois_search_debug",
			zap.String("message", "building request body"),
			zap.String("body", fmt.Sprintf("%v", whoisSearchRequest)),
		)
	}

	reqBody, _ := json.Marshal(whoisSearchRequest)
	req, err := http.NewRequest("POST", "https://api.dehashed.com/v2/whois/search", bytes.NewReader(reqBody))
	if err != nil {
		return whoisRecord, err
	}

	if w.debug {
		debug.PrintInfo("setting headers")
		zap.L().Info("whois_search_debug",
			zap.String("message", "setting headers"),
		)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", w.apiKey)

	if w.debug {
		debug.PrintInfo("headers set")
		debug.PrintJson(fmt.Sprintf("Headers: %v\n", req.Header.Clone()))
		debug.PrintInfo("performing request")
		zap.L().Info("whois_search_debug",
			zap.String("message", "performing request"),
		)
	}

	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		if w.debug {
			debug.PrintInfo("failed to perform request")
			debug.PrintError(err)
		}
		zap.L().Error("whois_search",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return whoisRecord, err
	}
	if res == nil {
		if w.debug {
			debug.PrintInfo("response was nil")
		}
		zap.L().Error("whois_search",
			zap.String("message", "response was nil"),
		)
		return whoisRecord, errors.New("response was nil")
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("whois_search",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return whoisRecord, err
	}

	// Check for HTTP status code errors
	if res.StatusCode != 200 {
		if w.debug {
			debug.PrintInfo("received error status code")
			debug.PrintJson(fmt.Sprintf("Status Code: %d\n", res.StatusCode))
			debug.PrintJson(fmt.Sprintf("Body: %s\n", string(b[:])))
		}

		dhErr := dehashed.GetDehashedError(res.StatusCode)
		fmt.Printf("[%d] API Error message: %s\n", res.StatusCode, dhErr.Error())
		zap.L().Error("whois_search",
			zap.String("message", "received error status code"),
			zap.Int("status_code", res.StatusCode),
			zap.String("error", dhErr.Error()),
		)
		return whoisRecord, &dhErr
	}

	err = json.Unmarshal(b, &whois)
	if err != nil {
		if w.debug {
			debug.PrintInfo("failed to unmarshal response body")
			debug.PrintError(err)
		}
		zap.L().Error("whois_search",
			zap.String("message", "failed to unmarshal response body"),
			zap.Error(err),
		)
		fmt.Println("Error unmarshalling response body:", err)
		fmt.Println("Response body:", string(b))
		return whoisRecord, err
	}

	if w.debug {
		debug.PrintInfo("unmarshalled response body")
		debug.PrintJson(fmt.Sprintf("Remaining Credits: %d\n", whois.RemainingCredits))
		debug.PrintJson(fmt.Sprintf("Data: %v\n", whois.Data))
	}
	w.balance = whois.RemainingCredits

	return whois.Data.WhoisRecord, nil
}

func (w *DehashedWhoIs) WhoisHistory(domain string) ([]sqlite.HistoryRecord, error) {
	var whois sqlite.WhoIsHistory
	var historyRecords []sqlite.HistoryRecord

	if w.debug {
		debug.PrintInfo("performing whois history search")
		zap.L().Info("whois_history_debug",
			zap.String("message", "performing whois history search"),
		)
	}

	whoisSearchRequest := DehashedWHOISSearchRequest{
		Domain:     domain,
		SearchType: "whois-history",
	}
	reqBody, _ := json.Marshal(whoisSearchRequest)

	if w.debug {
		debug.PrintInfo("building request body")
		debug.PrintJson(fmt.Sprintf("Request Body: %v\n", whoisSearchRequest))
		zap.L().Info("whois_history_debug",
			zap.String("message", "building request body"),
			zap.String("body", fmt.Sprintf("%v", whoisSearchRequest)),
		)
	}

	req, err := http.NewRequest("POST", "https://api.dehashed.com/v2/whois/search", bytes.NewReader(reqBody))
	if err != nil {
		if w.debug {
			debug.PrintInfo("failed to create request")
			debug.PrintError(err)
		}
		zap.L().Error("whois_history",
			zap.String("message", "failed to create request"),
			zap.Error(err),
		)
		return historyRecords, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", w.apiKey)

	if w.debug {
		debug.PrintInfo("performing request")
		debug.PrintJson(fmt.Sprintf("Headers: %v\n", req.Header.Clone()))
		zap.L().Info("whois_history_debug",
			zap.String("message", "performing request"),
		)
	}

	res, err := http.DefaultClient.Do(req)
	if res != nil {
		if w.debug {
			debug.PrintInfo("response was not nil")
		}
		zap.L().Info("whois_history",
			zap.String("message", "response was not nil"),
		)
		defer res.Body.Close()
	}
	if err != nil {
		if w.debug {
			debug.PrintInfo("failed to perform request")
			debug.PrintError(err)
		}
		zap.L().Error("whois_history",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return historyRecords, err
	}
	if res == nil {
		if w.debug {
			debug.PrintInfo("response was nil")
		}
		zap.L().Error("whois_history",
			zap.String("message", "response was nil"),
		)
		return historyRecords, errors.New("response was nil")
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		if w.debug {
			debug.PrintInfo("failed to read response body")
			debug.PrintError(err)
		}
		zap.L().Error("whois_history",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return historyRecords, err
	}

	// Check for HTTP status code errors
	if res.StatusCode != 200 {
		if w.debug {
			debug.PrintInfo("received error status code")
			debug.PrintJson(fmt.Sprintf("Status Code: %d\n", res.StatusCode))
			debug.PrintJson(fmt.Sprintf("Body: %s\n", string(b[:])))
		}
		dhErr := dehashed.GetDehashedError(res.StatusCode)
		fmt.Printf("[%d] API Error message: %s\n", res.StatusCode, dhErr.Error())
		zap.L().Error("whois_history",
			zap.String("message", "received error status code"),
			zap.Int("status_code", res.StatusCode),
			zap.String("error", dhErr.Error()),
			zap.String("body_error", string(b)),
		)
		return historyRecords, &dhErr
	}

	err = json.Unmarshal(b, &whois)
	if err != nil {
		if w.debug {
			debug.PrintInfo("failed to unmarshal response body")
			debug.PrintError(err)
		}
		zap.L().Error("whois_history",
			zap.String("message", "failed to unmarshal response body"),
			zap.Error(err),
		)
		fmt.Println("Error unmarshalling response body:", err)
		fmt.Println("Response body:", string(b))
		return historyRecords, err
	}

	if w.debug {
		debug.PrintInfo("unmarshalled response body")
		debug.PrintJson(fmt.Sprintf("Remaining Credits: %d\n", whois.RemainingCredits))
		debug.PrintJson(fmt.Sprintf("Data: %v\n", whois.Data))
	}
	w.balance = whois.RemainingCredits

	return whois.Data.Records, nil
}

func (w *DehashedWhoIs) ReverseWHOIS(include []string, exclude []string, reverseType string) (string, error) {
	if w.debug {
		debug.PrintInfo("performing reverse whois search")
		zap.L().Info("reverse_whois_debug",
			zap.String("message", "performing reverse whois search"),
		)
	}

	whoisSearchRequest := DehashedWHOISSearchRequest{
		Include:     include,
		Exclude:     exclude,
		ReverseType: reverseType,
		SearchType:  "reverse-whois",
	}
	reqBody, _ := json.Marshal(whoisSearchRequest)

	if w.debug {
		debug.PrintInfo("building request body")
		debug.PrintJson(fmt.Sprintf("Request Body: %v\n", whoisSearchRequest))
		zap.L().Info("reverse_whois_debug",
			zap.String("message", "building request body"),
			zap.String("body", fmt.Sprintf("%v", whoisSearchRequest)),
		)
	}

	req, err := http.NewRequest("POST", "https://api.dehashed.com/v2/whois/search", bytes.NewReader(reqBody))
	if err != nil {
		zap.L().Error("reverse_whois",
			zap.String("message", "failed to create request"),
			zap.Error(err),
		)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", w.apiKey)

	if w.debug {
		debug.PrintInfo("performing reverse whois search")
		debug.PrintJson(fmt.Sprintf("Headers: %v\n", req.Header.Clone()))
		zap.L().Info("reverse_whois_debug",
			zap.String("message", "performing reverse whois search"),
		)
	}

	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		if w.debug {
			debug.PrintInfo("failed to perform request")
			debug.PrintError(err)
		}
		zap.L().Error("reverse_whois",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return "", err
	}
	if res == nil {
		if w.debug {
			debug.PrintInfo("response was nil")
		}
		zap.L().Error("reverse_whois",
			zap.String("message", "response was nil"),
		)
		return "", errors.New("response was nil")
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		if w.debug {
			debug.PrintInfo("failed to read response body")
			debug.PrintError(err)
		}
		zap.L().Error("reverse_whois",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return "", err
	}

	// Check for HTTP status code errors
	if res.StatusCode != 200 {
		if w.debug {
			debug.PrintInfo("received error status code")
			debug.PrintJson(fmt.Sprintf("Status Code: %d\n", res.StatusCode))
			debug.PrintJson(fmt.Sprintf("Body: %s\n", string(b[:])))
		}
		dhErr := dehashed.GetDehashedError(res.StatusCode)
		fmt.Printf("[%d] API Error message: %s\n", res.StatusCode, dhErr.Error())
		zap.L().Error("reverse_whois",
			zap.String("message", "received error status code"),
			zap.Int("status_code", res.StatusCode),
			zap.String("error", dhErr.Error()),
			zap.String("body_error", string(b)),
		)
		return "", &dhErr
	}

	if w.debug {
		debug.PrintInfo("unmarshalled response body")
		debug.PrintJson(fmt.Sprintf("Body: %s\n", string(b[:])))
	}

	return string(b), nil
}

func (w *DehashedWhoIs) WhoisIP(ipAddress string) ([]sqlite.LookupResult, error) {
	if w.debug {
		debug.PrintInfo("performing whois ip search")
		zap.L().Info("whois_ip_debug",
			zap.String("message", "performing whois ip search"),
		)
	}

	type IPSearchRequest struct {
		IPAddress  string `json:"domain"`
		SearchType string `json:"search_type"`
	}

	whoisSearchRequest := IPSearchRequest{
		IPAddress:  ipAddress,
		SearchType: "reverse-ip",
	}
	reqBody, _ := json.Marshal(whoisSearchRequest)

	if w.debug {
		debug.PrintInfo("building request body")
		debug.PrintJson(fmt.Sprintf("Request Body: %v\n", whoisSearchRequest))
		zap.L().Info("whois_ip_debug",
			zap.String("message", "building request body"),
			zap.String("body", fmt.Sprintf("%v", whoisSearchRequest)),
		)
	}

	req, err := http.NewRequest("POST", "https://api.dehashed.com/v2/whois/search", bytes.NewReader(reqBody))
	if err != nil {
		zap.L().Error("whois_ip",
			zap.String("message", "failed to create request"),
			zap.Error(err),
		)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", w.apiKey)

	if w.debug {
		debug.PrintInfo("performing whois ip search")
		debug.PrintJson(fmt.Sprintf("Headers: %v\n", req.Header.Clone()))
		zap.L().Info("whois_ip_debug",
			zap.String("message", "performing whois ip search"),
		)
	}

	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		if w.debug {
			debug.PrintInfo("failed to perform request")
			debug.PrintError(err)
		}
		zap.L().Error("whois_ip",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return nil, err
	}
	if res == nil {
		if w.debug {
			debug.PrintInfo("response was nil")
		}
		zap.L().Error("whois_ip",
			zap.String("message", "response was nil"),
		)
		return nil, errors.New("response was nil")
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		if w.debug {
			debug.PrintInfo("failed to read response body")
			debug.PrintError(err)
		}
		zap.L().Error("whois_ip",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return nil, err
	}

	// Check for HTTP status code errors
	if res.StatusCode != 200 {
		if w.debug {
			debug.PrintInfo("received error status code")
			debug.PrintJson(fmt.Sprintf("Status Code: %d\n", res.StatusCode))
			debug.PrintJson(fmt.Sprintf("Body: %s\n", string(b[:])))
		}
		dhErr := dehashed.GetDehashedError(res.StatusCode)
		fmt.Printf("[%d] API Error message: %s\n", res.StatusCode, dhErr.Error())
		zap.L().Error("whois_ip",
			zap.String("message", "received error status code"),
			zap.Int("status_code", res.StatusCode),
			zap.String("error", dhErr.Error()),
			zap.String("body_error", string(b)),
		)
		return nil, &dhErr
	}

	if w.debug {
		debug.PrintInfo("read response body")
		debug.PrintJson(fmt.Sprintf("Response Body: %s\n", string(b)))
		zap.L().Info("whois_ip_debug",
			zap.String("message", "read response body"),
			zap.String("body", string(b)),
		)
	}

	var whois sqlite.WhoIsIPLookup
	err = json.Unmarshal(b, &whois)
	if err != nil {
		if w.debug {
			debug.PrintInfo("failed to unmarshal response body")
			debug.PrintError(err)
		}
		zap.L().Error("whois_ip",
			zap.String("message", "failed to unmarshal response body"),
			zap.Error(err),
		)
		return nil, err
	}

	if w.debug {
		debug.PrintInfo("unmarshalled response body")
		debug.PrintJson(fmt.Sprintf("Remaining Credits: %d\n", whois.RemainingCredits))
		debug.PrintJson(fmt.Sprintf("Data: %v\n", whois.Data))
	}
	w.balance = whois.RemainingCredits
	var lookups []sqlite.LookupResult

	for _, v := range whois.Data.Result {
		lookups = append(lookups, sqlite.LookupResult{
			FirstSeen:  v.FirstSeen,
			LastVisit:  v.LastVisit,
			Name:       v.Name,
			SearchTerm: ipAddress,
			Type:       "Reverse IP",
		})
	}

	sqlite.StoreIPLookup(lookups)

	return lookups, nil
}

func (w *DehashedWhoIs) WhoisMX(mxHostname string) ([]sqlite.LookupResult, error) {
	if w.debug {
		debug.PrintInfo("performing whois mx search")
		zap.L().Info("whois_mx_debug",
			zap.String("message", "performing whois mx search"),
		)
	}

	type ReverseMX struct {
		Domain     string `json:"domain"`
		SearchType string `json:"search_type"`
	}

	whoisSearchRequest := ReverseMX{
		Domain:     mxHostname,
		SearchType: "reverse-mx",
	}
	reqBody, _ := json.Marshal(whoisSearchRequest)

	if w.debug {
		debug.PrintInfo("building request body")
		debug.PrintJson(fmt.Sprintf("Request Body: %v\n", whoisSearchRequest))
		zap.L().Info("whois_mx_debug",
			zap.String("message", "building request body"),
			zap.String("body", fmt.Sprintf("%v", whoisSearchRequest)),
		)
	}

	req, err := http.NewRequest("POST", "https://api.dehashed.com/v2/whois/search", bytes.NewReader(reqBody))
	if err != nil {
		if w.debug {
			debug.PrintInfo("failed to create request")
			debug.PrintError(err)
		}
		zap.L().Error("whois_mx",
			zap.String("message", "failed to create request"),
			zap.Error(err),
		)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", w.apiKey)

	if w.debug {
		debug.PrintInfo("performing whois mx search")
		debug.PrintJson(fmt.Sprintf("Headers: %v\n", req.Header.Clone()))
		zap.L().Info("whois_mx_debug",
			zap.String("message", "performing whois mx search"),
		)
	}

	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		if w.debug {
			debug.PrintInfo("failed to perform request")
			debug.PrintError(err)
		}
		zap.L().Error("whois_mx",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return nil, err
	}
	if res == nil {
		if w.debug {
			debug.PrintInfo("response was nil")
		}
		zap.L().Error("whois_mx",
			zap.String("message", "response was nil"),
		)
		return nil, errors.New("response was nil")
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		if w.debug {
			debug.PrintInfo("failed to read response body")
			debug.PrintError(err)
		}
		zap.L().Error("whois_mx",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return nil, err
	}

	// Check for HTTP status code errors
	if res.StatusCode != 200 {
		if w.debug {
			debug.PrintInfo("received error status code")
			debug.PrintJson(fmt.Sprintf("Status Code: %d\n", res.StatusCode))
			debug.PrintJson(fmt.Sprintf("Body: %s\n", string(b[:])))
		}
		dhErr := dehashed.GetDehashedError(res.StatusCode)
		fmt.Printf("[%d] API Error message: %s\n", res.StatusCode, dhErr.Error())
		zap.L().Error("whois_mx",
			zap.String("message", "received error status code"),
			zap.Int("status_code", res.StatusCode),
			zap.String("error", dhErr.Error()),
			zap.String("body_error", string(b)),
		)
		return nil, &dhErr
	}

	if w.debug {
		debug.PrintInfo("read response body")
		debug.PrintJson(fmt.Sprintf("Body: %s\n", string(b[:])))
	}

	var whois sqlite.WhoIsMXLookup
	err = json.Unmarshal(b, &whois)
	if err != nil {
		if w.debug {
			debug.PrintInfo("failed to unmarshal response body")
			debug.PrintError(err)
		}
		zap.L().Error("whois_mx",
			zap.String("message", "failed to unmarshal response body"),
			zap.Error(err),
		)
		return nil, err
	}

	if w.debug {
		debug.PrintInfo("unmarshalled response body")
		debug.PrintJson(fmt.Sprintf("Remaining Credits: %d\n", whois.RemainingCredits))
		debug.PrintJson(fmt.Sprintf("Data: %v\n", whois.Data))
	}

	var mxLookups []sqlite.LookupResult

	for _, v := range whois.Data.Result {
		mxLookups = append(mxLookups, sqlite.LookupResult{
			FirstSeen:  v.FirstSeen,
			LastVisit:  v.LastVisit,
			Name:       v.Name,
			SearchTerm: mxHostname,
			Type:       "MX",
		})
	}

	sqlite.StoreIPLookup(mxLookups)

	return mxLookups, nil
}

func (w *DehashedWhoIs) WhoisNS(nsHostname string) ([]sqlite.LookupResult, error) {
	if w.debug {
		debug.PrintInfo("performing whois ns search")
		zap.L().Info("whois_ns_debug",
			zap.String("message", "performing whois ns search"),
		)
	}

	type NSLookup struct {
		Domain     string `json:"domain"`
		SearchType string `json:"search_type"`
	}

	whoisSearchRequest := NSLookup{
		Domain:     nsHostname,
		SearchType: "reverse-ns",
	}

	reqBody, _ := json.Marshal(whoisSearchRequest)

	if w.debug {
		debug.PrintInfo("building request body")
		debug.PrintJson(fmt.Sprintf("Request Body: %v\n", whoisSearchRequest))
		zap.L().Info("whois_ns_debug",
			zap.String("message", "building request body"),
			zap.String("body", fmt.Sprintf("%v", whoisSearchRequest)),
		)
	}

	req, err := http.NewRequest("POST", "https://api.dehashed.com/v2/whois/search", bytes.NewReader(reqBody))
	if err != nil {
		zap.L().Error("whois_ns",
			zap.String("message", "failed to create request"),
			zap.Error(err),
		)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", w.apiKey)

	if w.debug {
		debug.PrintInfo("performing request")
		debug.PrintJson(fmt.Sprintf("Headers: %v\n", req.Header.Clone()))
		zap.L().Info("whois_ns_debug",
			zap.String("message", "performing request"),
		)
	}

	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		if w.debug {
			debug.PrintInfo("failed to perform request")
			debug.PrintError(err)
		}
		zap.L().Error("whois_ns",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return nil, err
	}
	if res == nil {
		if w.debug {
			debug.PrintInfo("response was nil")
		}
		zap.L().Error("whois_ns",
			zap.String("message", "response was nil"),
		)
		return nil, errors.New("response was nil")
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("whois_ns",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return nil, err
	}

	// Check for HTTP status code errors
	if res.StatusCode != 200 {
		if w.debug {
			debug.PrintInfo("received error status code")
			debug.PrintJson(fmt.Sprintf("Status Code: %d\n", res.StatusCode))
			debug.PrintJson(fmt.Sprintf("Body: %s\n", string(b[:])))
		}
		dhErr := dehashed.GetDehashedError(res.StatusCode)
		fmt.Printf("[%d] API Error message: %s\n", res.StatusCode, dhErr.Error())
		zap.L().Error("whois_ns",
			zap.String("message", "received error status code"),
			zap.Int("status_code", res.StatusCode),
			zap.String("error", dhErr.Error()),
			zap.String("body_error", string(b)),
		)
		return nil, &dhErr
	}

	if w.debug {
		debug.PrintInfo("read response body")
		debug.PrintJson(fmt.Sprintf("Body: %s\n", string(b[:])))
		zap.L().Info("whois_ns_debug",
			zap.String("message", "read response body"),
			zap.String("body", string(b)),
		)
	}

	var whois sqlite.WhoIsNSLookup
	err = json.Unmarshal(b, &whois)
	if err != nil {
		if w.debug {
			debug.PrintInfo("failed to unmarshal response body")
			debug.PrintError(err)
		}
		zap.L().Error("whois_ns",
			zap.String("message", "failed to unmarshal response body"),
			zap.Error(err),
		)
		return nil, err
	}

	if w.debug {
		debug.PrintInfo("unmarshalled response body")
		debug.PrintJson(fmt.Sprintf("Remaining Credits: %d\n", whois.RemainingCredits))
		debug.PrintJson(fmt.Sprintf("Data: %v\n", whois.Data))
	}
	w.balance = whois.RemainingCredits
	var nsLookups []sqlite.LookupResult

	for _, v := range whois.Data.Result {
		nsLookups = append(nsLookups, sqlite.LookupResult{
			FirstSeen:  v.FirstSeen,
			LastVisit:  v.LastVisit,
			Name:       v.Name,
			SearchTerm: nsHostname,
			Type:       "NS",
		})
	}

	sqlite.StoreIPLookup(nsLookups)

	return nsLookups, nil
}

func (w *DehashedWhoIs) WhoisSubdomainScan(domain string) ([]sqlite.SubdomainRecord, error) {
	var whois sqlite.WhoIsSubdomainScan
	var subdomains []sqlite.SubdomainRecord

	if w.debug {
		debug.PrintInfo("performing whois subdomain scan")
		zap.L().Info("whois_subdomain_scan_debug",
			zap.String("message", "performing whois subdomain scan"),
		)
	}

	whoisSearchRequest := DehashedWHOISSearchRequest{
		Domain:     domain,
		SearchType: "subdomain-scan",
	}

	reqBody, _ := json.Marshal(whoisSearchRequest)

	if w.debug {
		debug.PrintInfo("building request body")
		debug.PrintJson(fmt.Sprintf("Request Body: %v\n", whoisSearchRequest))
		zap.L().Info("whois_subdomain_scan_debug",
			zap.String("message", "building request body"),
			zap.String("body", fmt.Sprintf("%v", whoisSearchRequest)),
		)
	}

	req, err := http.NewRequest("POST", "https://api.dehashed.com/v2/whois/search", bytes.NewReader(reqBody))
	if err != nil {
		if w.debug {
			debug.PrintInfo("failed to create request")
			debug.PrintError(err)
		}
		zap.L().Error("whois_subdomain_scan",
			zap.String("message", "failed to create request"),
			zap.Error(err),
		)
		return subdomains, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", w.apiKey)

	if w.debug {
		debug.PrintInfo("performing request")
		debug.PrintJson(fmt.Sprintf("Headers: %v\n", req.Header.Clone()))
		zap.L().Info("whois_subdomain_scan_debug",
			zap.String("message", "performing request"),
		)
	}

	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		if w.debug {
			debug.PrintInfo("failed to perform request")
			debug.PrintError(err)
		}
		zap.L().Error("whois_subdomain_scan",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return subdomains, err
	}
	if res == nil {
		if w.debug {
			debug.PrintInfo("response was nil")
		}
		zap.L().Error("whois_subdomain_scan",
			zap.String("message", "response was nil"),
		)
		return subdomains, errors.New("response was nil")
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		if w.debug {
			debug.PrintInfo("failed to read response body")
			debug.PrintError(err)
		}
		zap.L().Error("whois_subdomain_scan",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return subdomains, err
	}

	// Check for HTTP status code errors
	if res.StatusCode != 200 {
		if w.debug {
			debug.PrintInfo("received error status code")
			debug.PrintJson(fmt.Sprintf("Status Code: %d\n", res.StatusCode))
			debug.PrintJson(fmt.Sprintf("Body: %s\n", string(b[:])))
		}
		dhErr := dehashed.GetDehashedError(res.StatusCode)
		fmt.Printf("[%d] API Error message: %s\n", res.StatusCode, dhErr.Error())
		zap.L().Error("whois_subdomain_scan",
			zap.String("message", "received error status code"),
			zap.Int("status_code", res.StatusCode),
			zap.String("error", dhErr.Error()),
			zap.String("body_error", string(b)),
		)
		return subdomains, &dhErr
	}

	err = json.Unmarshal(b, &whois)
	if err != nil {
		zap.L().Error("whois_subdomain_scan",
			zap.String("message", "failed to unmarshal response body"),
			zap.Error(err),
		)
		fmt.Println("Error unmarshalling response body:", err)
		fmt.Println("Response body:", string(b))
		return subdomains, err
	}

	if w.debug {
		debug.PrintInfo("unmarshalled response body")
		debug.PrintJson(fmt.Sprintf("Remaining Credits: %d\n", whois.RemainingCredits))
		debug.PrintJson(fmt.Sprintf("Data: %v\n", whois.Data))
	}
	w.balance = whois.RemainingCredits

	return whois.Data.Result.Records, nil
}

func (w *DehashedWhoIs) Balance() (int, error) {
	if w.debug {
		debug.PrintInfo("getting whois credits")
		zap.L().Info("whois_debug",
			zap.String("message", "getting whois credits"),
		)
	}
	return w.getBalance()
}

func (w *DehashedWhoIs) getBalance() (int, error) {
	var whoisCredits sqlite.WhoIsCredits

	req, err := http.NewRequest("GET", "https://api.dehashed.com/v2/whois/credits", nil)
	if err != nil {
		if w.debug {
			debug.PrintInfo("failed to create request")
			debug.PrintError(err)
		}
		return whoisCredits.WhoisCredits, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Dehashed-Api-Key", w.apiKey)

	if w.debug {
		debug.PrintInfo("performing request")
		h := req.Header.Clone()
		debug.PrintJson(fmt.Sprintf("Headers: %v\n", h))
		zap.L().Info("whois_debug",
			zap.String("message", "performing request"),
		)
	}

	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		if w.debug {
			debug.PrintInfo("failed to perform request")
			debug.PrintError(err)
		}
		zap.L().Error("get_whois_credits",
			zap.String("message", "failed to perform request"),
			zap.Error(err),
		)
		return whoisCredits.WhoisCredits, err
	}
	if res == nil {
		if w.debug {
			debug.PrintInfo("response was nil")
		}
		zap.L().Error("get_whois_credits",
			zap.String("message", "response was nil"),
		)
		return whoisCredits.WhoisCredits, errors.New("response was nil")
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		if w.debug {
			debug.PrintInfo("failed to read response body")
			debug.PrintError(err)
		}
		zap.L().Error("get_whois_credits",
			zap.String("message", "failed to read response body"),
			zap.Error(err),
		)
		return whoisCredits.WhoisCredits, err
	}

	// Check for HTTP status code errors
	if res.StatusCode != 200 {
		if w.debug {
			debug.PrintInfo("received error status code")
			debug.PrintJson(fmt.Sprintf("Status Code: %d\n", res.StatusCode))
			debug.PrintJson(fmt.Sprintf("Body: %s\n", string(b[:])))
		}
		dhErr := dehashed.GetDehashedError(res.StatusCode)
		fmt.Printf("[%d] API Error message: %s\n", res.StatusCode, dhErr.Error())
		zap.L().Error("get_whois_credits",
			zap.String("message", "received error status code"),
			zap.Int("status_code", res.StatusCode),
			zap.String("error", dhErr.Error()),
			zap.String("body_error", string(b)),
		)
		return whoisCredits.WhoisCredits, &dhErr
	}

	err = json.Unmarshal(b, &whoisCredits)
	if err != nil {
		zap.L().Error("get_whois_credits",
			zap.String("message", "failed to unmarshal response body"),
			zap.Error(err),
		)
		fmt.Println("Error unmarshalling response body:", err)
		fmt.Println("Response body:", string(b))
		return whoisCredits.WhoisCredits, err
	}

	if w.debug {
		debug.PrintInfo("unmarshalled response body")
		debug.PrintJson(fmt.Sprintf("Remaining Credits: %d\n", whoisCredits.WhoisCredits))
	}

	return whoisCredits.WhoisCredits, nil
}
