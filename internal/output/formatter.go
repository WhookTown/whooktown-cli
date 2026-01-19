package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// Formatter defines the output interface
type Formatter interface {
	// Format outputs data to the writer
	Format(w io.Writer, data interface{}) error

	// FormatTable outputs tabular data with headers
	FormatTable(w io.Writer, headers []string, rows [][]string) error

	// Success prints a success message
	Success(msg string)

	// Error prints an error message
	Error(msg string)

	// Info prints an info message
	Info(msg string)
}

// OutputFormat represents output type
type OutputFormat string

const (
	FormatTable OutputFormat = "table"
	FormatJSON  OutputFormat = "json"
)

// New creates a formatter based on format type
func New(format OutputFormat) Formatter {
	switch format {
	case FormatJSON:
		return &JSONFormatter{out: os.Stdout, errOut: os.Stderr}
	default:
		return &TableFormatter{out: os.Stdout, errOut: os.Stderr}
	}
}

// TableFormatter outputs human-readable tables
type TableFormatter struct {
	out    io.Writer
	errOut io.Writer
}

func (f *TableFormatter) Format(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(data) // Fallback to JSON for complex data
}

func (f *TableFormatter) FormatTable(w io.Writer, headers []string, rows [][]string) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	// Print headers
	for i, h := range headers {
		if i > 0 {
			fmt.Fprint(tw, "\t")
		}
		fmt.Fprint(tw, h)
	}
	fmt.Fprintln(tw)

	// Print rows
	for _, row := range rows {
		for i, cell := range row {
			if i > 0 {
				fmt.Fprint(tw, "\t")
			}
			fmt.Fprint(tw, cell)
		}
		fmt.Fprintln(tw)
	}

	return tw.Flush()
}

func (f *TableFormatter) Success(msg string) {
	fmt.Fprintf(f.out, "OK %s\n", msg)
}

func (f *TableFormatter) Error(msg string) {
	fmt.Fprintf(f.errOut, "ERR %s\n", msg)
}

func (f *TableFormatter) Info(msg string) {
	fmt.Fprintf(f.out, "%s\n", msg)
}

// JSONFormatter outputs JSON
type JSONFormatter struct {
	out    io.Writer
	errOut io.Writer
}

func (f *JSONFormatter) Format(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func (f *JSONFormatter) FormatTable(w io.Writer, headers []string, rows [][]string) error {
	// Convert to array of maps for JSON
	result := make([]map[string]string, len(rows))
	for i, row := range rows {
		result[i] = make(map[string]string)
		for j, cell := range row {
			if j < len(headers) {
				result[i][headers[j]] = cell
			}
		}
	}
	return f.Format(w, result)
}

func (f *JSONFormatter) Success(msg string) {
	f.Format(f.out, map[string]string{"status": "success", "message": msg})
}

func (f *JSONFormatter) Error(msg string) {
	f.Format(f.errOut, map[string]string{"status": "error", "message": msg})
}

func (f *JSONFormatter) Info(msg string) {
	f.Format(f.out, map[string]string{"message": msg})
}
