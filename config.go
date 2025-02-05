package main

import "regexp"

var BEXPRESSION_REGEX = regexp.MustCompile(`(?:%((?:w(?:\d+)?|s(?:\d+)?|c|C|d|@)+))(?:(\{\d+(?:,\d+)?\})|(\(\w(?:,\w)\))){0,2}`)

// Charsets
var CHARACTER_SMALL string = "abcdefghijklmnopqrstuvwxyz"
var CHARACTER_BIG string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var DIGIT = "0123456789"
var SYMBOL string = "`~!@#$%^&*()_-+={[}]:;\"'<,>.?/|\\"
