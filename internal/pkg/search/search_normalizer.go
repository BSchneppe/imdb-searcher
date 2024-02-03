package search

import (
	"html"
	"strings"

	"golang.org/x/text/unicode/norm"
)

func NormalizeString(instr string) string {
	instr = html.UnescapeString(instr)              // Unescape HTML entities
	instr = norm.NFC.String(instr)                  // Unicode normalization
	instr = substituteRuneF(strings.ToLower(instr)) // Substitute runes
	instr = strings.TrimSpace(instr)                // Trim spaces
	instr = removeStopWords(instr)                  // Remove stop words
	return instr
}

var stopWords = map[string]bool{
	"i": true, "me": true, "my": true, "myself": true,
	"we": true, "our": true, "ours": true, "ourselves": true,
	"you": true, "your": true, "yours": true, "yourself": true, "yourselves": true,
	"he": true, "him": true, "his": true, "himself": true,
	"she": true, "her": true, "hers": true, "herself": true,
	"it": true, "its": true, "itself": true,
	"they": true, "them": true, "their": true, "theirs": true, "themselves": true,
	"what": true, "which": true, "who": true, "whom": true,
	"this": true, "that": true, "these": true, "those": true,
	"am": true, "is": true, "are": true, "was": true, "were": true,
	"be": true, "been": true, "being": true,
	"have": true, "has": true, "had": true, "having": true,
	"do": true, "does": true, "did": true, "doing": true,
	"a": true, "an": true, "the": true,
	"and": true, "but": true, "if": true, "or": true, "because": true, "as": true,
	"until": true, "while": true, "of": true, "at": true, "by": true,
	"for": true, "with": true, "about": true, "against": true,
	"between": true, "into": true, "through": true, "during": true,
	"before": true, "after": true,
	"above": true, "below": true, "to": true, "from": true,
	"up": true, "down": true, "in": true, "out": true, "on": true, "off": true,
	"over": true, "under": true, "again": true, "further": true,
	"then": true, "once": true, "here": true, "there": true, "when": true,
	"where": true, "why": true, "how": true,
	"all": true, "any": true, "both": true, "each": true, "few": true,
	"more": true, "most": true, "other": true, "some": true, "such": true,
	"no": true, "nor": true, "not": true, "only": true, "own": true, "same": true,
	"so": true, "than": true, "too": true, "very": true, "s": true, "t": true,
	"can": true, "will": true, "just": true, "don": true, "should": true, "now": true,
}

func removeStopWords(instr string) string {

	// Tokenize the input string
	words := strings.Fields(instr)

	// Filter out the stopwords
	filteredWords := make([]string, 0, len(words))
	for _, word := range words {
		if _, found := stopWords[strings.ToLower(word)]; !found {
			filteredWords = append(filteredWords, word)
		}
	}

	// Reassemble the string
	return strings.Join(filteredWords, " ")
}

var subRune = map[rune]string{
	'&':  "and",
	'@':  "at",
	'"':  "",
	'\'': "",
	'’':  "",
	'_':  "",
	'‒':  "-", // figure dash
	'–':  "-", // en dash
	'—':  "-", // em dash
	'―':  "-", // horizontal bar
	'ä':  "ae",
	'Ä':  "Ae",
	'ö':  "oe",
	'Ö':  "Oe",
	'ü':  "ue",
	'Ü':  "Ue",
	'ß':  "ss",
}

func substituteRuneF(s string) string {
	var buf strings.Builder
	buf.Grow(len(s))

	for _, c := range s {
		if repl, ok := subRune[c]; ok {
			buf.WriteString(repl)
		} else {
			buf.WriteRune(c)
		}
	}
	return buf.String()
}
