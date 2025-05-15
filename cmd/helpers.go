package cmd

import (
	"dehasher/internal/sqlite"
	"strings"
)

type Table int64

const (
	ResultsTable Table = iota
	RunsTable
	CredsTable
	WhoIsTable
	SubdomainsTable
	HistoryTable
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
	default:
		return UnknownTable
	}
}

func (t Table) Object() interface{} {
	switch t {
	case ResultsTable:
		return sqlite.Result{}
	case RunsTable:
		return sqlite.QueryOptions{}
	case CredsTable:
		return sqlite.Creds{}
	case WhoIsTable:
		return sqlite.WhoisRecord{}
	case SubdomainsTable:
		return sqlite.SubdomainRecord{}
	case HistoryTable:
		return sqlite.HistoryRecord{}
	default:
		return nil
	}
}
