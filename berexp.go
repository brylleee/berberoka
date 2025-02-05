package main

import (
	"strconv"
)

type Expression struct {
	Data     string
	IsBerexp bool
}

func backslashEscape(s string) string {
	var result []rune
	n := len(s)

	for i := 0; i < n; i++ {
		if s[i] == '\\' {
			// If next character is also '\', keep both and skip next
			if i+1 < n && s[i+1] == '\\' {
				result = append(result, '\\', '\\')
				i++ // Skip next '\\'
			}
			// Else, remove single '\'
		} else {
			result = append(result, rune(s[i]))
		}
	}

	return string(result)
}

// Find all metacharacters from given string ang separate from static ones
// RETURNS: an array of Expression struct, containing tokens of metacharacters and static string
func Parse(bexpression string) []Expression {
	bexpression_match_indexes := BEXPRESSION_REGEX.FindAllStringIndex(bexpression, -1)

	// Return if string has no berexp
	if len(bexpression_match_indexes) == 0 {
		return []Expression{{Data: bexpression, IsBerexp: false}}
	}

	index := 0
	var parsed []Expression

	for i, value := range bexpression_match_indexes {
		// Add non-berexpression part before the match
		if index < value[0] {
			parsed = append(parsed, Expression{Data: bexpression[index:value[0]], IsBerexp: false})
		}

		// Add the matched berexpression
		parsed = append(parsed, Expression{Data: bexpression[value[0]:value[1]], IsBerexp: true})

		index = value[1]

		// Add the remaining expression
		if i == len(bexpression_match_indexes)-1 && index < len(bexpression) {
			parsed = append(parsed, Expression{Data: bexpression[index:], IsBerexp: false})
		}
	}

	return parsed
}

// Process tokens by replacing them with proper characters or words from
// specified wordlists and charsets. The result is directly written into a file
func Process(parsed []Expression) {
	var entry string

	var Cyclers []*Cycler
	var previousCycler *Cycler = nil

	// Create a cycler
	for _, expression := range parsed {
		if !expression.IsBerexp {
			continue
		}

		berexpParts := BEXPRESSION_REGEX.FindStringSubmatch(expression.Data)
		var charsets []string

		// Process WORDLISTS
		if berexpParts[1][0] == 'w' {
			index, _ := strconv.Atoi(berexpParts[1][1:])
			for _, words := range WORDLISTS[index] {
				charsets = append(charsets, string(words))
			}
		} else if berexpParts[1][0] == 's' { // Process CUSTOMCHARSETS
			index, _ := strconv.Atoi(berexpParts[1][1:])
			for _, chr := range CHARSETS[index] {
				charsets = append(charsets, string(chr))
			}
		} else {
			for _, charset := range berexpParts[1] { // Process REGULAR METACHARACTERS
				switch charset {
				case 'd':
					for _, chr := range DIGIT {
						charsets = append(charsets, string(chr))
					}
				case 'c':
					for _, chr := range CHARACTER_SMALL {
						charsets = append(charsets, string(chr))
					}
				case 'C':
					for _, chr := range CHARACTER_BIG {
						charsets = append(charsets, string(chr))
					}
				case '@':
					for _, chr := range SYMBOL {
						charsets = append(charsets, string(chr))
					}
				}
			}
		}

		Cyclers = append(Cyclers, &Cycler{
			Charsets:       charsets,
			Current:        string(charsets[0]),
			Steps:          0,
			PreviousCycler: previousCycler,
		})

		previousCycler = Cyclers[len(Cyclers)-1]
	}

	// Cycle through every tokens, and increment them
	// to exhaust all possible combinations
	for {
		entry = ""
		cyclerCount := 0
		Cyclers[len(Cyclers)-1].cycle()

		for _, expression := range parsed {
			if !expression.IsBerexp {
				entry += expression.Data
				continue
			}

			entry += Cyclers[cyclerCount].Current
			cyclerCount++
		}

		FILE.WriteString(backslashEscape(entry) + "\n")
		// fmt.Println(entry)

		if Cyclers[0].LastFinish {
			break
		}
	}

	FILE.Close()
}
