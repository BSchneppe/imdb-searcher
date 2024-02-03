package search

import "testing"

func TestNormalizeString(t *testing.T) {
	normalized := NormalizeString("The quick brown fox jumps over the lazy dog")
	if normalized != "quick brown fox jumps lazy dog" {
		t.Errorf("NormalizeString() returned %s, expected the quick brown fox jumps over the lazy dog", normalized)
	}
}

func TestNormalizeString_empty_string(t *testing.T) {
	normalized := NormalizeString("")
	if normalized != "" {
		t.Errorf("NormalizeString() returned %s, expected an empty string", normalized)
	}
}

func TestNormalizeString_single_word(t *testing.T) {
	normalized := NormalizeString("word")
	if normalized != "word" {
		t.Errorf("NormalizeString() returned %s, expected word", normalized)
	}
}

func TestNormalizeString_single_word_with_spaces(t *testing.T) {
	normalized := NormalizeString(" word ")
	if normalized != "word" {
		t.Errorf("NormalizeString() returned %s, expected word", normalized)
	}
}

func TestNormalizeStringWithAmbiguousChar(t *testing.T) {
	normalized := NormalizeString("SpongeBob SquarePants & Patrick Star")
	if normalized != "spongebob squarepants patrick star" {
		t.Errorf("NormalizeString() returned %s, expected spongebob squarepants patrick star", normalized)
	}

}
