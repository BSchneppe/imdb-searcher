package tsv_reader

import (
	"strings"
	"testing"
)

func TestTabNewlineReader_Read(t *testing.T) {
	tab_reader := NewTabNewlineReader(strings.NewReader("a\tb\tc\n1\t2\t3\n"))

	first_line, err := tab_reader.Read()
	if err != nil {
		t.Errorf("Read() returned an error: %v", err)
	}
	if len(first_line) != 3 {
		t.Errorf("Read() returned a slice of length %d, expected 3", len(first_line))
	}
	if first_line[0] != "a" || first_line[1] != "b" || first_line[2] != "c" {
		t.Errorf("Read() returned %v, expected [a b c]", first_line)
	}
	second_line, err := tab_reader.Read()
	if err != nil {
		t.Errorf("Read() returned an error: %v", err)
	}
	if len(second_line) != 3 {
		t.Errorf("Read() returned a slice of length %d, expected 3", len(second_line))
	}
	if second_line[0] != "1" || second_line[1] != "2" || second_line[2] != "3" {
		t.Errorf("Read() returned %v, expected [1 2 3]", second_line)
	}
	_, err = tab_reader.Read()
	if err == nil {
		t.Errorf("Read() didn't return an error on EOF")
	}

}

func TestTabNewlineReader_Read_last_line_no_newline(t *testing.T) {
	tab_reader := NewTabNewlineReader(strings.NewReader("a\tb\tc"))

	first_line, err := tab_reader.Read()
	if err != nil {
		t.Errorf("Read() returned an error: %v", err)
	}
	if len(first_line) != 3 {
		t.Errorf("Read() returned a slice of length %d, expected 3", len(first_line))
	}
	if first_line[0] != "a" || first_line[1] != "b" || first_line[2] != "c" {
		t.Errorf("Read() returned %v, expected [a b c]", first_line)
	}
}
