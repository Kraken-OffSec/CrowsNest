package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	// Add Subcommand to db command
	rootCmd.AddCommand(exportCmd)

}

// DB export command
var (
	exportLimitRows int
	exportListAll   bool
	exportTableName string
	exportNotNull   string
	exportColumns   string
	exportUserQuery string
	exportRawQuery  string
	exportFormat    string
	exportFile      string

	exportCmd = &cobra.Command{
		Use:   "export",
		Short: "Export database to file",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Exporting database...")
		},
	}
)

func init() {
	// Add flags specific to export command
	exportCmd.Flags().IntVarP(&exportLimitRows, "limit", "l", 100, "Limit number of results")
	exportCmd.Flags().BoolVarP(&exportListAll, "list-all", "a", false, "List all columns")
	exportCmd.Flags().StringVarP(&exportTableName, "table", "t", "", "Table to export")
	exportCmd.Flags().StringVarP(&exportNotNull, "not-null", "n", "", "Filter for non-null values (comma-separated list, e.g., 'password,email')")
	exportCmd.Flags().StringVarP(&exportColumns, "columns", "c", "", "Columns to display in output (comma-separated list, e.g., 'username,email,password')")
	exportCmd.Flags().StringVarP(&exportUserQuery, "query", "q", "", "User query to execute")
	exportCmd.Flags().StringVarP(&exportRawQuery, "raw-query", "r", "", "Raw SQL query to execute")
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "json", "Output format (json, yaml, xml, txt)")
	exportCmd.Flags().StringVarP(&exportFile, "file", "o", "export", "File to output results to including extension")
}
