package sqlite

import "strings"

type Table int64

const (
	ResultsTable Table = iota
	RunsTable
	CredsTable
	WhoIsTable
	SubdomainsTable
	HistoryTable
	LookupTable
	UnknownTable
)

func GetTable(userInput string) Table {
	switch strings.ToLower(userInput) {
	case "results":
		return ResultsTable
	case "runs":
		return RunsTable
	case "creds":
		return CredsTable
	case "whois":
		return WhoIsTable
	case "subdomains":
		return SubdomainsTable
	case "history":
		return HistoryTable
	case "lookup":
		return LookupTable
	default:
		return UnknownTable
	}
}

func (t Table) Object() interface{} {
	switch t {
	case ResultsTable:
		return Result{}
	case RunsTable:
		return QueryOptions{}
	case CredsTable:
		return Creds{}
	case WhoIsTable:
		return WhoisRecord{}
	case SubdomainsTable:
		return SubdomainRecord{}
	case HistoryTable:
		return HistoryRecord{}
	case LookupTable:
		return LookupResult{}
	default:
		return nil
	}
}
