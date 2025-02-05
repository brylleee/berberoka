package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// Flag options
var BEREXP string = ""
var OUTPUT_FILE string = ""
var FILE *os.File

// User supplied wordlists and charsets
var WORDLISTS map[int][]string
var CHARSETS map[int]string

func help(bannerOnly bool) {
	banner := `
▄▄▄▄   ▓█████  ██▀███   ▄▄▄▄   ▓█████  ██▀███   ▒█████   ██ ▄█▀▄▄▄      
▓█████▄ ▓█   ▀ ▓██ ▒ ██▒▓█████▄ ▓█   ▀ ▓██ ▒ ██▒▒██▒  ██▒ ██▄█▒▒████▄    
▒██▒ ▄██▒███   ▓██ ░▄█ ▒▒██▒ ▄██▒███   ▓██ ░▄█ ▒▒██░  ██▒▓███▄░▒██  ▀█▄  
▒██░█▀  ▒▓█  ▄ ▒██▀▀█▄  ▒██░█▀  ▒▓█  ▄ ▒██▀▀█▄  ▒██   ██░▓██ █▄░██▄▄▄▄██ 
░▓█  ▀█▓░▒████▒░██▓ ▒██▒░▓█  ▀█▓░▒████▒░██▓ ▒██▒░ ████▓▒░▒██▒ █▄▓█   ▓██▒
░▒▓███▀▒░░ ▒░ ░░ ▒▓ ░▒▓░░▒▓███▀▒░░ ▒░ ░░ ▒▓ ░▒▓░░ ▒░▒░▒░ ▒ ▒▒ ▓▒▒▒   ▓▒█░
▒░▒   ░  ░ ░  ░  ░▒ ░ ▒░▒░▒   ░  ░ ░  ░  ░▒ ░ ▒░  ░ ▒ ▒░ ░ ░▒ ▒░ ▒   ▒▒ ░
░    ░    ░     ░░   ░  ░    ░    ░     ░░   ░ ░ ░ ░ ▒  ░ ░░ ░  ░   ▒   
░         ░  ░   ░      ░         ░  ░   ░         ░ ░  ░  ░        ░  ░
      ░                       ░                                          
berberoka v1.0 - ka1ro | The Wordlist Crafter (because crunching isn't enough)`

	usage := "\n[=] Usage: %s [Berberoka Expression] [Output Wordlist] [...OPTIONS]"

	manual := `
Metacharacters:
	%w(digit)        word from wordlist
	%s(digit)        customized charset from flag
	%d               digit
	%c               character
	%C               - all caps
	%@               symbols

Metacharacters may be joined together:
	%dc              alphanumberic characters
	%@C              capital letters with symbols
	%dcC@            letters, digits and symbols in one charset

Metacharacter modifiers:
	%_{2,3}          2-3 duplicates
	%_(b,e)          from b to e only

Examples:
	Password Format: [4 digit pin]
	Berexp:          %d%d%d%d
	Commandline:     ./berberoka "%d%d%d%d" pins.lst

	Password Format: PC-[2 digit number]-USER_[Name]
	Berexp:          PC-%d%d-USER_%w1
	Commandline:     ./berberoka "PC-%d%d-USER_%w1" password.lst -w1 names.lst

	Password Format: 2020-[LastName]-[FirstName]-[3 distinct random numbers]
	Berexp:          2020-%w1-%w2-%d%d%d
	Commandline:     ./berberoka "2020-%w1-%w2-%d%d%d" password.lst -w1 lastnames.lst -w2 firstnames.lst

Custom lists:
    -w(digit) <wordlist>  Include a custom wordlist to use in a berexp
    -s(digit) <charset>   Include a custom charset to use in a berexp

    -h  Output this help menu
    `

	// Gradient banner color
	numberOfColors := 4
	g := 50
	b := 250
	bannerLines := strings.Split(banner, "\n")
	for ctr := 0; ctr < len(bannerLines); ctr++ {
		color.RGB(0, g, b).Println(bannerLines[ctr])

		g = g + 5*(numberOfColors%(len(bannerLines)/numberOfColors))
		b = b - 25*(numberOfColors%(len(bannerLines)/numberOfColors))
	}

	fmt.Println()

	if !bannerOnly {
		fmt.Println(fmt.Sprintf(usage, os.Args[0]))
		fmt.Println(manual)
	}
}

func main() {
	errorC := color.New(color.FgRed).PrintfFunc()
	warnC := color.New(color.FgYellow).PrintfFunc()
	infoC := color.New(color.FgBlue).PrintfFunc()

	WORDLISTS = make(map[int][]string)
	CHARSETS = make(map[int]string)

	if len(os.Args) == 1 {
		help(true)
		warnC("[=] ")
		fmt.Println("Use flag -h for help")
	}

	// Flag processing
	for i := 1; i < len(os.Args); i++ {
		// Flags are strictly checked for first two characters only
		switch os.Args[i][:2] {
		// Specify wordlist file
		case "-w":
			index, err := strconv.Atoi(os.Args[i][2:])
			if err != nil {
				errorC("[!] ")
				log.Fatal("ERROR: -w flag must be followed by an int to mark uniqueness (ex.: -w1 <wordlist file>)")
			}

			wordlistFileName := os.Args[i+1]

			wordlistFile, err := os.Open(wordlistFileName)
			if err != nil {
				errorC("[!] ")
				log.Fatal("ERROR: Cannot find or access file `" + string(wordlistFileName) + "`")
			}

			defer wordlistFile.Close()

			infoC("[=] ")
			fmt.Println("Loading wordlist `" + wordlistFileName + "` in memory...")

			buffer := bufio.NewScanner(wordlistFile)
			for buffer.Scan() {
				WORDLISTS[index] = append(WORDLISTS[index], buffer.Text())
			}

			infoC("[=] ")
			fmt.Println("Wordlist loaded in memory...")

			i++
		// Specify character set string
		case "-s":
			index, err := strconv.Atoi(os.Args[i][2:])
			if err != nil {
				errorC("[!] ")
				log.Fatal("ERROR: -s flag must be followed by an int to mark uniqueness (ex.: -c1 \"<charset>\")")
			}

			infoC("[=] ")
			fmt.Println("Loading custom charset in memory...")

			CHARSETS[index] = os.Args[i+1]
			i++
		// Output help
		case "-h":
			help(false)
			os.Exit(0)
		default:
			// If both berexpression and output file has been read as flag
			if BEREXP != "" && OUTPUT_FILE != "" {
				errorC("[!] ")
				log.Fatal("ERROR: Unrecognized flag " + os.Args[i])
			}

			// If berexpression is not yet initialized, treat unknown flag as berexpression
			// else as the output wordlist file
			if BEREXP == "" {
				BEREXP = os.Args[i]
			} else if OUTPUT_FILE == "" {
				OUTPUT_FILE = os.Args[i]
			}
		}
	}

	if OUTPUT_FILE == "" {
		errorC("[!] ")
		log.Fatal("ERROR: No output wordlist file specified")
	}

	if _, err := os.Stat(OUTPUT_FILE); errors.Is(err, os.ErrNotExist) {
		infoC("[=] ")
		fmt.Println("Creating wordlist file `" + OUTPUT_FILE + "`...")

		FILE, err = os.OpenFile(OUTPUT_FILE, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			errorC("[!] ")
			log.Fatal("ERROR: Unable to create file `" + OUTPUT_FILE + "`!")
		}
	} else if err == nil {
		infoC("[=] ")
		fmt.Println("Wordlist file `" + OUTPUT_FILE + "` already exists. Will append to wordlist...")
		FILE, err = os.OpenFile(OUTPUT_FILE, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			errorC("[!] ")
			log.Fatal("ERROR: Unable to open file `" + OUTPUT_FILE + "`!")
		}
	} else {
		errorC("[!] ")
		log.Fatal("ERROR: Something bad happened while trying to create your wordlist file!")
	}

	Process(Parse(BEREXP))

	infoC("[=] ")
	fmt.Println("Wordlist locked and loaded.")
}
