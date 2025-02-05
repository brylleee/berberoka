# Berberoka
The Wordlist Crafter (because crunching isn't enough)
---

![image](https://github.com/user-attachments/assets/c4dd64c5-6730-470d-a631-68189d4fe384)

### What is Berberoka
**Berberoka** is similar to and is inspired by Hashcat's **maskprocessor**. Both tools allows you to generate wordlists that can be configured per-charset
It is extremely powerful for generating a wordlist from a known format and can be used in situations wherein default passwords follow a certain format.
It  can also be used to generate pins, and modify, integrate, and create even more specific wordlists from existing wordlists.

### Comparison with existing tools

| Feature | Berberoka | Crunch |
| :---: | --- | --- |
| Generate wordlist from criteria | ✅ | ✅ |
| Specify custom character sets | ✅ | ⭕ Not as rigid |
| **Use existing wordlists** | ✅ | ❌ |
| Generate permutations | ❌ | ✅ |

| Feature | Berberoka | maskprocessor |
| :---: | --- | --- |
| Generate wordlist from criteria | ✅ | ✅ |
| Specify custom character sets | ✅ | ✅ |
| **Use existing wordlists** | ✅ | ❌ |
| Generate permutations | ❌ | ❌ |

### Usage
Usage: `./berberoka [Berberoka Expression] [Output Wordlist] [...OPTIONS]`

Flags:
```
    -w(digit) <wordlist>  Include a custom wordlist to use in a berexp
    -s(digit) <charset>   Include a custom charset to use in a berexp

    -h  Output this help menu
```

### Berexpressions
Berberoka utilizes expressions similar to Regular Expressions to create its wordlists.
It works like an inverse regex; A regex matches strings, a berexp creates strings that match.

Metacharacters:
```
	%w(digit)        word from wordlist
	%s(digit)        customized charset from flag
	%d               digit
	%c               character
	%C               - all caps
	%@               symbols
```

Metacharacters may be joined together:
```
	%dc              alphanumberic characters
	%@C              capital letters with symbols
	%dcC@            letters, digits and symbols in one charset
```

Examples:
```
	Password Format: [4 digit pin]
	Berexp:          %d%d%d%d
	Commandline:     ./berberoka "%d%d%d%d" pins.lst

	Password Format: PC-[2 digit number]-USER_[Name]
	Berexp:          PC-%d%d-USER_%w1
	Commandline:     ./berberoka "PC-%d%d-USER_%w1" password.lst -w1 names.lst

	Password Format: 2020-[LastName]-[FirstName]-[3 distinct random numbers]
	Berexp:          2020-%w1-%w2-%d%d%d
	Commandline:     ./berberoka "2020-%w1-%w2-%d%d%d" password.lst -w1 lastnames.lst -w2 firstnames.lst
```

### Installation
1) Clone the project
   
   `git clone https://github.com/brylleee/berberoka/`
   
2) Build using Go (install [here](https://go.dev/doc/install))

   `go build`

3) Run

   `./berberoka`

---

This tool is in active development, so there will be more features to be integrated soon.
