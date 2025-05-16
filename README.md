# 🚀 Dehasher
### A CLI tool for seamless interaction with the Dehashed API

---

## 🌟 Features
- **Output Format Control**: JSON, YAML, XML, and TEXT support.
- **Regex & Wildcard Matching**: Flexible query options.
- **Local Database Storage**: Default or custom paths.
- **Database Querying**: Raw SQL and filtered queries.
- **Enhanced Logging**: Easy log parsing and rotation.
- **Error Handling**: Intelligent API error management.
- **WhoIs Lookups**: Domain, IP, MX, NS, and more.
- **Subdomain Scanning**: Identify subdomains.
- **Robust Logging**: Detailed logs for debugging.
- **API Key Management**: Securely store and manage API keys.
- **Formatted Output**: Easy to read and understand.
- **Intuitive Database Querying**: Query for specific information.

---

## 📦 Installation

Clone the repository and build the tool:
```bash
git clone https://github.com/Ar1ste1a/Dehasher.git
cd Dehasher
go build dehasher.go
```

<hr></hr>

## 🔰 Getting Started

To begin, clone the repository
``` bash-session
git clone https://github.com/Ar1ste1a/Dehasher.git
cd Dehasher
go build dehasher.go
```

<hr></hr>

## 🛠️ Initial Setup

Dehasher requires an API key from Dehashed. Set it up with:
```bash
ar1ste1a@kali:~$ dehasher set-key <redacted>
```

<hr></hr>

## 🗄️ Database Configuration

Dehasher supports two database storage options:

1. **Default Path** (default): Stores the database at `~/.local/share/Dehasher/db/dehashed.sqlite`
2. **Local Path**: Stores the database in the current directory as `./dehasher.sqlite`

The **Local Path** option allows for separate databases for different projects or engagements.

To configure the database location:

```bash
# Use local database in current directory
./dehasher set-local-db true

# Use default database path
./dehasher set-local-db false
```

<hr></hr>

## 🔍 Crafting Queries

### Simple Query
Dehasher can be used simply for example to query for credentials matching a given email domain.
``` go
# Provide credentials for emails matching @target.com
dehasher api -D @target.com -C
```

### Simple Credentials Query
Dehasher can also be used to return only credentials for a given query.
``` go
# Provide credentials for emails matching @target.com
dehasher api -E @target.com -C
```

### Multiple Match Query
Dehasher is capable of handling multiple queries for the same field.  
This is useful for when you want to search for multiple domains, or multiple usernames.
``` go
# Provide credentials for emails matching @target.com and @target2.com
dehasher api -E @target.com,@target2.com -C
```

### Wildcard Query
Dehasher is capable of handling wildcard queries.  
A wildcard query cannot begin with a wildcard.  
This is a limitation of the Dehashed API.
An asterisk can be used to denote multiple characters, and a question mark can be used to denote a single character.
![Alt text](.img/wildcard_sample.png "Wildcard Query")
``` go
# Provide credentials for emails matching @target.com and @target2.com
dehasher api -E @target?.com -C -W
```


### Regex Query
Dehasher is capable of handling regex queries.  
Simply denote regex queries with the `-R` flag.
Place all regex queries in quotes with the corresponding query flag in single quotes.
``` go
# Return matches for emails matching this given regex query
dehasher api -R -e '[a-zA-Z0-9]+(?:\.[a-zA-Z0-9]+)?@target.com'
```

### Output Text (default JSON)
Dehasher is capable of handling output formats.  
The default output format is JSON.  
To change the output format, use the `-f` flag.  
Dehasher currently supports JSON, YAML, XML, and TEXT output formats.
``` go
# Return matches for usernames exactly matching "admin" and write to text file 'admins_file.txt'
dehasher api -U admin -o admins_file -f txt
```

<hr></hr>

## 🌐 WhoIs Lookups
Dehasher supports WHOIS lookups, history searches, reverse WHOIS searches, IP lookups, MX lookups, NS lookups, and subdomain scans.
The WhoIs Lookups require a separate API Credit from the Dehashed API.

### Domain Lookup
Dehasher can perform a domain lookup for a given domain.
This provides a tree view of the domain's WHOIS information.
![Alt text](.img/tree_whois_lookup.png "WhoIs Tree View")
```bash
# Perform a WHOIS lookup for example.com
dehasher whois -d example.com
```

### History Lookup
History Lookups require 25 credits. 
This is a Dehashed API limitation.
The history lookup is immediately written to file and not displayed in the terminal or stored in the database.
```bash
# Perform a WHOIS history search for example.com
dehasher whois -d example.com -H
```

### Reverse WHOIS Lookup
Dehasher can perform a reverse WHOIS lookup for given criteria.  
This provides a list of all domains that match the given query.  
The reverse WHOIS lookup is immediately written to file and not displayed in the terminal or stored in the database.
```bash
# Perform a reverse WHOIS lookup for example.com
dehasher whois -I example.com
```

### IP Lookup
Dehasher can perform a reverse IP lookup for a given IP address.  
This provides a list of all domains that match the given query.
![Alt text](.img/reverse_ip_lookup.png "WhoIs Tree View")
```bash
# Perform a reverse IP lookup for 8.8.8.8
dehasher whois -i 8.8.8.8
```

### MX Lookup
Dehasher can perform an MX lookup for a given MX hostname.  
This provides a list of all domains that match the given query.
![Alt text](.img/mx_lookup.png "WhoIs Tree View")
```bash
# Perform a reverse MX lookup for google.com
dehasher whois -m google.com
```
### NS Lookup
Dehasher can perform an NS lookup for a given NS hostname.  
This provides a list of all domains that match the given query.
The picture below also includes the --debug global flag.
![Alt text](.img/debug_ns_search.png "WhoIs Tree View")
```bash
# Perform a reverse NS lookup for google.com
dehasher whois -n google.com
```
### Subdomain Scan
Dehasher can perform a subdomain scan for a given domain.  
This provides a list of all subdomains that match the given query.
![Alt text](.img/subdomains_lookup.png "WhoIs Tree View")
```bash
# Perform a WHOIS subdomain scan for google.com
dehasher whois -d google.com -s
```

---

## 📊 Database Querying
Dehasher stores query results in a local database.  
This database can be queried for previous results.
This database also includes WhoIs Information and Subdomain Scan results, but does **not** include historical lookups.

## Simple Query
![Alt text](.img/simple_query_db.png "Simple Query")

Dehasher supports querying the database for previous results.  
This is useful for when you want to query for specific information.
```bash
# Query the database for all results containing the word 'admin' in the username
dehasher query -t results -q "username LIKE '%admin%'"
```


## Raw SQL Queries
![Alt text](.img/raw_query_db.png "Raw Query")

Dehasher also supports raw SQL queries.  This is useful for when you want to query for specific information.
```bash
# Query the database for all results containing the word 'admin' in the username
dehasher query -r "SELECT * FROM results WHERE username LIKE '%admin%'"
```

## Query Options
Dehasher supports a number of query options.  These options can be used to filter the results of a query.
```bash
# Query the database for all results containing the word 'admin' in the username
dehasher query -t results -q "username LIKE '%admin%'" -n username,email,password
```

## Listing Tables and Columns
Dehasher supports listing all available tables and columns.  
This is useful for when you want to query for specific information.
```bash
# List all available tables and columns
dehasher query -a
```

The current tables available for query are:
- results
  - Results from a dehashed query
- creds
  - Credentials parsed from dehashed results
- whois
  - Results from a whois record lookup
- subdomains
  - Subdomains discovered in a whois subdomain scan 
- runs
  - Previous query runs to the dehashed API
- lookup 
  - Results of any Whois NS, MX, or IP lookup

---

# Exporting Results
Dehasher supports exporting results to a file.  
This is useful for when you want to requery for specific information without touching the Dehashed API.
The export subcommand supports all the same options as the query subcommand.
The export subcommand also supports file naming and output format control.
```bash
# Export all results containing the word 'admin' in the username to a text file
dehasher export -t results -q "username LIKE '%admin%'" -o admins_file -f txt
```

## 🐛 Debugging

Dehasher uses the `zap` logging library for logging.  The logs are stored in `~/.local/share/Dehasher/logs`.
The logs can be easily queried from the Dehasher CLI.
```bash
# Show the last 10 logs
dehasher logs -l 10

# Show logs from the last 24 hours
dehasher logs -s "24 hours ago"

# Show logs from the last 24 hours with a severity of error or fatal
dehasher logs -s "24 hours ago" -v error,fatal
```
## 🎉 Sample Run
```bash
ar1ste1a@kali:~$ dehasher api -D <redacted>.com -o <redacted> -f json
Making 3 Requests for 10000 Records (30000 Total)
[*] Querying Dehashed API...
	[*] Performing Request...
		[+] Retrieved 2740 Records
        [-] Not Enough Entries, ending queries
	[+] Discovered 10 Credentials
	[*] Writing entries to file: <redacted>.json
		[*] Success
[*] Completing Process
```

## 🤝 Contributing
Contributions are welcome! Submit a pull request to help improve Dehasher.



<div align="center">
    <img src="https://img.wanman.io/fUSu0/jUtovIFE52.png/raw" style="width: 350px; height: auto" alt="Ar1ste1a" title="Ar1ste1a Offensive Security">
</div>

## **Release The Kraken**
