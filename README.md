<div align="center">
    <img src=.img/crowsnest.png  style="width: 500px; height: auto" alt="Ar1ste1a" title="CrowsNest Logo">
</div>

### A CLI tool for seamless interaction with the Dehashed and Hunter.io APIs.

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
- **Person and Company Enrichment**: Retrieve detailed information about people and companies.
- **Email Verification**: Verify the existence and quality of email addresses.

---

## 🔰 Getting Started

To begin, clone the repository
``` bash-session
git clone https://github.com/Ar1ste1a/CrowsNest.git
cd crowsnest
go build crowsnest.go
```

---

## 🗄️ Database Configuration

CrowsNest supports two database storage options:

1. **Default Path** (default): Stores the database at `~/.local/share/crowsnest/db/dehashed.sqlite`
2. **Local Path**: Stores the database in the current directory as `./crowsnest.sqlite`

The **Local Path** option allows for separate databases for different projects or engagements.

To configure the database location:

```bash
# Use local database in current directory
./crowsnest set-local-db true

# Use default database path
./crowsnest set-local-db false
```

---

## 🌐 Dehashed

###️ Initial Setup

CrowsNest requires an API key from Dehashed. Set it up with:
```bash
ar1ste1a@kali:~$ crowsnest set-dehashed <redacted>
```

### Simple Query
CrowsNest can be used simply for example to query for credentials matching a given email domain.
``` go
# Provide credentials for domains matching target.com
crowsnest api -D target.com -C
```

### Simple Credentials Query
CrowsNest can also be used to return only credentials for a given query.
``` go
# Provide credentials for emails matching @target.com
crowsnest api -E @target.com -C
```

### Multiple Match Query
CrowsNest is capable of handling multiple queries for the same field.  
This is useful for when you want to search for multiple domains, or multiple usernames.
``` go
# Provide credentials for domains matching target.com and target2.com, retrieving only credentials
crowsnest api -D target.com,target2.com -C
```

### Wildcard Query
CrowsNest is capable of handling wildcard queries.  
A wildcard query cannot begin with a wildcard.  
This is a limitation of the Dehashed API.
An asterisk can be used to denote multiple characters, and a question mark can be used to denote a single character.
![Alt text](.img/wildcard_sample.png "Wildcard Query")
``` go
# Provide credentials for emails matching @target.com and @target2.com
crowsnest api -E @target?.com -C -W
```

### Email Query
Dehashed has dictated that emails should be searched in the following format:
`email:target.name&domain:target.com`.
As such, to query an email, please use the following format (note, wildcard is not required but can be useful):
``` go
# Provide credentials for emails matching target.*@target.com
crowsnest api -W -E 'target*' -D target.com
```
You may also query the domain and find emails as well
``` go
# Provide credentials for emails matching target.com
crowsnest api -D target.com -C
```


### Regex Query
CrowsNest is capable of handling regex queries.  
Simply denote regex queries with the `-R` flag.
Place all regex queries in quotes with the corresponding query flag in single quotes.
``` go
# Return matches for emails matching this given regex query
crowsnest api -R -E '[a-zA-Z0-9]+(?:\.[a-zA-Z0-9]+)?@target.com'
```

### Output Text (default JSON)
CrowsNest is capable of handling output formats.  
The default output format is JSON.  
To change the output format, use the `-f` flag.  
CrowsNest currently supports JSON, YAML, XML, and TEXT output formats.
``` go
# Return matches for usernames exactly matching "admin" and write to text file 'admins_file.txt'
crowsnest api -U admin -o admins_file -f txt
```

---

## 🌐 WhoIs Lookups
CrowsNest supports WHOIS lookups, history searches, reverse WHOIS searches, IP lookups, MX lookups, NS lookups, and subdomain scans.
The WhoIs Lookups require a separate API Credit from the Dehashed API.

### Domain Lookup
CrowsNest can perform a domain lookup for a given domain.
This provides a tree view of the domain's WHOIS information.
![Alt text](.img/tree_whois_lookup.png "WhoIs Tree View")
```bash
# Perform a WHOIS lookup for example.com
crowsnest whois -d example.com
```

### History Lookup
History Lookups require 25 credits. 
This is a Dehashed API limitation.
The history lookup is immediately written to file and not displayed in the terminal or stored in the database.
```bash
# Perform a WHOIS history search for example.com
crowsnest whois -d example.com -H
```

### Reverse WHOIS Lookup
CrowsNest can perform a reverse WHOIS lookup for given criteria.  
This provides a list of all domains that match the given query.  
The reverse WHOIS lookup is immediately written to file and not displayed in the terminal or stored in the database.
```bash
# Perform a reverse WHOIS lookup for example.com
crowsnest whois -I example.com
```

### IP Lookup
CrowsNest can perform a reverse IP lookup for a given IP address.  
This provides a list of all domains that match the given query.
![Alt text](.img/reverse_ip_lookup.png "WhoIs Tree View")
```bash
# Perform a reverse IP lookup for 8.8.8.8
crowsnest whois -i 8.8.8.8
```

### MX Lookup
CrowsNest can perform an MX lookup for a given MX hostname.  
This provides a list of all domains that match the given query.
![Alt text](.img/mx_lookup.png "WhoIs Tree View")
```bash
# Perform a reverse MX lookup for google.com
crowsnest whois -m google.com
```
### NS Lookup
CrowsNest can perform an NS lookup for a given NS hostname.  
This provides a list of all domains that match the given query.
The picture below also includes the --debug global flag.
![Alt text](.img/debug_ns_search.png "WhoIs Tree View")
```bash
# Perform a reverse NS lookup for google.com
crowsnest whois -n google.com
```
### Subdomain Scan
CrowsNest can perform a subdomain scan for a given domain.  
This provides a list of all subdomains that match the given query.
![Alt text](.img/subdomains_lookup.png "WhoIs Tree View")
```bash
# Perform a WHOIS subdomain scan for google.com
crowsnest whois -d google.com -s
```

---

## 🌐 Hunter.io
CrowsNest supports Hunter.io lookups.  
Hunter.io lookups require a separate API Key from the Dehashed API.
This can be set using the `set-hunter` command.
```bash
# Set the Hunter.io API key
crowsnest set-hunter <redacted>
```

### Domain Search
CrowsNest can perform a domain search for a given domain.  
This provides information about  company including a description, social media information and any technologies in use.
![Alt text](.img/hunter_domain_search.png "Hunter.io Domain Search")
```bash
# Perform a Hunter.io domain search for example.com
crowsnest hunter -d example.com -D
```

### Email Finder
CrowsNest can perform an email finder search for a given domain, first name, and last name.  
This provides information about a user including a confidence score, and any social media accounts linked to a first name, last name and email.
![Alt text](.img/hunter_email_finder.png "Hunter.io Email Finder")
```bash
# Perform a Hunter.io email finder search for example.com
crowsnest hunter -d example.com -F John -L Doe -E
```

### Email Verification
CrowsNest can perform an email verification search for a given email.  
This provides a verification and score of a given email address.
![Alt text](.img/email_verification.png "Hunter.io Email Verification")
```bash
# Perform a Hunter.io email verification search for example@target.com
crowsnest hunter -e example@target.com -V
```

### Company Enrichment
CrowsNest can perform a company enrichment search for a given domain.  
This provides information about a company given its domain.
![Alt text](.img/company_enrichment.png "Hunter.io Company Enrichment")
```bash
# Perform a Hunter.io company enrichment search for example.com
crowsnest hunter -d example.com -C
```

### Person Enrichment
CrowsNest can perform a person enrichment search for a given email.  
This provides information about a user given an email address..
![Alt text](.img/person_enrichment.png "Hunter.io Person Enrichment")
```bash
# Perform a Hunter.io person enrichment search for example@target.com
crowsnest hunter -e example@target.com -P
```

### Combined Enrichment
CrowsNest can perform a combined enrichment search for a given email.  
This is a combination of the company and person enrichments given an email address.
![Alt text](.img/combined_enrichment_1.png "Hunter.io Combined Enrichment")
![Alt text](.img/combined_enrichment_2.png "Hunter.io Combined Enrichment")
```bash
# Perform a Hunter.io combined enrichment search for example@target.com
crowsnest hunter -e example@target.com -B
```

---
## 📊 Database Querying
CrowsNest stores query results in a local database.  
This database can be queried for previous results.
This is useful for when you want to query for specific information.
This database also includes WhoIs Information and Subdomain Scan results, but does **not** include historical lookups.

## Simple Query
#### It's possible to query the database using shorthand and without knowing any SQL at all.
#### The following queries the results table where username is not null, only showing the username, email and password columns.
![Alt text](.img/simple_query_db.png "Simple Query")
#### You may also add in a simple query using the `-q` flag. The following displays a 'LIKE' clause on the email column.
#### Note the '%\<clause\>%' is still required.
![Alt text](.img/simple_where.png "Simple Query")
  
```bash
# Query the database for all results containing the word 'admin' in the username
crowsnest query -t results -q "username LIKE '%admin%'"
```


## Raw SQL Queries
![Alt text](.img/raw_query_db.png "Raw Query")

CrowsNest also supports raw SQL queries.  This is useful for when you want to query for specific information.
```bash
# Query the database for all results containing the word 'admin' in the username
crowsnest query -r "SELECT * FROM results WHERE username LIKE '%admin%'"
```

## Query Options
CrowsNest supports a number of query options.  These options can be used to filter the results of a query.
```bash
# Query the database for all results containing the word 'admin' in the username
crowsnest query -t results -q "username LIKE '%admin%'" -n username,email,password
```

## Listing Tables and Columns
CrowsNest supports listing all available tables and columns.  
This is useful for when you want to query for specific information.
```bash
# List all available tables and columns
crowsnest query -a
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
- hunter_domain
  - Results from a hunter.io domain search
- hunter_email
  - Results extracted from a domain saerch and email finder.

---

# Exporting Results
CrowsNest supports exporting results to a file.  
This is useful for when you want to requery for specific information without touching the Dehashed API.
The export subcommand supports all the same options as the query subcommand.
The export subcommand also supports file naming and output format control.
```bash
# Export all results containing the word 'admin' in the username to a text file
crowsnest export -t results -q "username LIKE '%admin%'" -o admins_file -f txt
```

## 🐛 Debugging

CrowsNest uses the `zap` logging library for logging.  The logs are stored in `~/.local/share/crowsnest/logs`.
The logs can be easily queried from the crowsnest CLI.

### Logs Dates
#### crowsnest utilized 'easy time' to determine the appropriate time for a given query.
![Alt text](.img/easy_time_parsing.png "Easy Time")
#### You may also used dates mixed with easy time to perform queries.
![Alt text](.img/mixed_time_query.png "Mixed Time")
#### The following formats are supported:
- `last 24 hours`
- `last 2 days`
- `30 minutes ago`
- `45 seconds ago`
- `1 week ago`
- `05-01-2025`
- `05/01/2025`
- `05/01/25`
- `05-01-25`
- `May 01, 2025`

```bash
# Show the last 10 logs
crowsnest logs -l 10

# Show logs from the last 24 hours
crowsnest logs -s "last 24 hours"

# Show logs from the last 24 hours with a severity of error or fatal
crowsnest logs -s "05-01-2025" -v error,fatal
```

## 🎉 Sample Run
```bash
ar1ste1a@kali:~$ crowsnest api -D <redacted>.com -o <redacted> -f json
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
Contributions are welcome! Submit a pull request to help improve CrowsNest.



<div align="center">
    <img src="https://img.wanman.io/fUSu0/jUtovIFE52.png/raw" style="width: 350px; height: auto" alt="Ar1ste1a" title="Ar1ste1a Offensive Security">
</div>

## **Release The Kraken**
