package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	pal "github.com/abusomani/go-palette/palette"
)

// Flag options
var BEREXP string = ""
var OUTPUT_FILE string = ""
var FILE *os.File

// User supplied wordlists and charsets
var WORDLISTS map[int]string
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
%w(digit)        word from wordlist [exhausts every line in wordlist]
%d               digit
%c               character
    %C           all caps
    %Cc          both caps
%a               alphanumberic
    %A           all caps
    %Aa          both caps
%@               symbols

Metacharacter modifiers:
%_{2,3}          2-3 duplicates
%_(b,e)          from b to e only

Misc metacharacters:
%(...|...)       or grouping
%[2001,2005]     continuous number

Example:
A company gives out accounts with default password with format
(Birthmonth first 3 letters)-(Birthyear)-(Random 4 letter word)

Berexp: %w1-%(19|20)%d%d-%Cc%Cc%Cc%Cc   (Supply -w1 months.list, where months.list contains Jan, Feb ... Dec)
Output: Jan-1900-aaaa ... Dec-2099-ZZZZ

Custom lists:
    -w(digit) <wordlist>  Include a custom wordlist to use in a berexp
    -c(digit) <charset>   Include a custom charset to use in a berexp

    -h  Output this help menu
    `

	fmt.Println(banner)

	if !bannerOnly {
		fmt.Println(fmt.Sprintf(usage, os.Args[0]))
		fmt.Println(manual)
	}
}

func main() {
	p := pal.New()
	p.Print()

	WORDLISTS = make(map[int]string)
	CHARSETS = make(map[int]string)

	if len(os.Args) == 1 {
		help(true)
		fmt.Println("[=] Use flag -h for help")
	}

	// Flag processing
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i][:2] {
		case "-w":
			index, err := strconv.Atoi(os.Args[i][2:])
			if err != nil {
				log.Fatal("[!] ERROR: -w flag must be followed by an int to mark uniqueness (ex.: -w1 <wordlist>)")
			}

			WORDLISTS[index] = os.Args[i+1]
			i++
		case "-c":
			index, err := strconv.Atoi(os.Args[i][2:])
			if err != nil {
				log.Fatal("[!] ERROR: -c flag must be followed by an int to mark uniqueness (ex.: -c1 <charset>)")
			}

			CHARSETS[index] = os.Args[i+1]
			i++
		case "-h":
			help(false)
		default:
			if BEREXP != "" && OUTPUT_FILE != "" {
				log.Fatal("[!] ERROR: Unrecognized flag " + os.Args[i])
			}

			if BEREXP == "" {
				BEREXP = os.Args[i]
			} else if OUTPUT_FILE == "" {
				OUTPUT_FILE = os.Args[i]
			}
		}
	}

	if OUTPUT_FILE == "" {
		log.Fatal("[!] ERROR: No output wordlist file specified")
	}

	if _, err := os.Stat(OUTPUT_FILE); errors.Is(err, os.ErrNotExist) {
		fmt.Println("[=] Creating wordlist file `" + OUTPUT_FILE + "`...")

		FILE, err = os.OpenFile(OUTPUT_FILE, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("[!] ERROR: Unable to create file `" + OUTPUT_FILE + "`!")
		}
	} else if err == nil {
		fmt.Println("[=] Wordlist file `" + OUTPUT_FILE + "` already exists. Will append to wordlist...")
		FILE, err = os.OpenFile(OUTPUT_FILE, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("[!] ERROR: Unable to open file `" + OUTPUT_FILE + "`!")
		}
	} else {
		log.Fatal("[!] ERROR: Something bad happened while trying to create your wordlist file!")
	}

	Process(Parse(BEREXP))

	fmt.Println("[=] Wordlist locked and loaded.")
}
