package main

import "regexp"

var BEXPRESSION_REGEX = regexp.MustCompile(`(?:%(Cc|Aa|w(?:\d+)?|d|c(?:\d+)?|C|a|A|@|\((.*?)\))(?:(\{\d+(?:,\d+)?\})|(\(\w(?:,\w)\))){0,3})`)

// Charsets
var CHARACTER_SMALL string = "abcdefghijklmnopqrstuvwxyz"
var CHARACTER_BIG string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var DIGIT = "0123456789"
var ALPHANUMBERIC_SMALL string = CHARACTER_SMALL + DIGIT
var ALPHANUMBERIC_BIG string = CHARACTER_BIG + DIGIT
var SYMBOL string = "`~!@#$%^&*()_-+={[}]:;\"'<,>.?/|\\"
