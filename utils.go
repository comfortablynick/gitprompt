package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// PrettyPrint prints objects in a readable format for debugging
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Sprint(string(b))
	}
	return
}

// Detent removes leading tab from string
func detent(s string) string {
	return regexp.MustCompile("(?m)^[\t]*").ReplaceAllString(s, "")
}

// standardizeSpaces removes extra spaces and trims string
func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
