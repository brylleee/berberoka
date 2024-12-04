package main

type Expression struct {
	Data     string
	IsBerexp bool
}

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

func Process(parsed []Expression) {
	var entry string

	var Cyclers []*Cycler
	var previousCycler *Cycler = nil

	// Create a cycler
	for _, expression := range parsed {
		if !expression.IsBerexp {
			continue
		}

		// Capture groups
		// 0 - Entire
		// 1 - Metachar
		// 2 - Or group expression
		// 3 - Repeat value
		// 4 - Limit value
		berexpParts := BEXPRESSION_REGEX.FindStringSubmatch(expression.Data)

		switch berexpParts[1] {
		case "d":
			Cyclers = append(Cyclers, &Cycler{
				Charsets:       DIGIT,
				Current:        string(DIGIT[0]),
				Steps:          0,
				PreviousCycler: previousCycler,
			})
		case "c":
			Cyclers = append(Cyclers, &Cycler{
				Charsets:       CHARACTER_SMALL,
				Current:        string(CHARACTER_SMALL[0]),
				Steps:          0,
				PreviousCycler: previousCycler,
			})
		case "C":
			Cyclers = append(Cyclers, &Cycler{
				Charsets:       CHARACTER_BIG,
				Current:        string(CHARACTER_BIG[0]),
				Steps:          0,
				PreviousCycler: previousCycler,
			})
		case "Cc":
			Cyclers = append(Cyclers, &Cycler{
				Charsets:       CHARACTER_SMALL + CHARACTER_BIG,
				Current:        string(CHARACTER_SMALL[0]),
				Steps:          0,
				PreviousCycler: previousCycler,
			})
		case "a":
			Cyclers = append(Cyclers, &Cycler{
				Charsets:       ALPHANUMBERIC_SMALL,
				Current:        string(ALPHANUMBERIC_SMALL[0]),
				Steps:          0,
				PreviousCycler: previousCycler,
			})
		case "A":
			Cyclers = append(Cyclers, &Cycler{
				Charsets:       ALPHANUMBERIC_BIG,
				Current:        string(ALPHANUMBERIC_BIG[0]),
				Steps:          0,
				PreviousCycler: previousCycler,
			})
		case "Aa":
			Cyclers = append(Cyclers, &Cycler{
				Charsets:       ALPHANUMBERIC_SMALL + ALPHANUMBERIC_BIG,
				Current:        string(ALPHANUMBERIC_SMALL[0]),
				Steps:          0,
				PreviousCycler: previousCycler,
			})
		case "@":
			Cyclers = append(Cyclers, &Cycler{
				Charsets:       SYMBOL,
				Current:        string(SYMBOL[0]),
				Steps:          0,
				PreviousCycler: previousCycler,
			})
		default:
			// CUSTOM WORDLISTS
		}

		previousCycler = Cyclers[len(Cyclers)-1]
	}

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

		FILE.WriteString(entry + "\n")
		// fmt.Println(entry)

		if Cyclers[0].LastFinish {
			break
		}
	}

	FILE.Close()
}
