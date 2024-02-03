package tsv_reader

import (
	"bufio"
	"io"
	"strings"
)

// TabNewlineReader is a custom reader that wraps another io.Reader and reads data line by line.
type TabNewlineReader struct {
	reader *bufio.Reader
}

// NewTabNewlineReader creates a new TabNewlineReader wrapping the provided io.Reader.
func NewTabNewlineReader(r io.Reader) *TabNewlineReader {
	return &TabNewlineReader{
		reader: bufio.NewReader(r),
	}
}

// Read reads one line from the underlying reader and returns a slice of strings,
// splitting the line by tabs.
func (r *TabNewlineReader) Read() ([]string, error) {
	line, err := r.reader.ReadString('\n')
	if err != nil {
		if err == io.EOF && len(line) > 0 {
			// Return the last line even if it doesn't end with a newline
			return strings.Split(line, "\t"), nil
		}
		return nil, err
	}

	// Split the line by tabs and return the fields
	return strings.Split(strings.TrimSuffix(line, "\n"), "\t"), nil
}
