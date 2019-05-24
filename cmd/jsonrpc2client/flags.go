package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
)

// FatalUsage report usage error in same way as flag.Parse().
func FatalUsage(format string, a ...interface{}) {
	if format != "" {
		fmt.Fprintf(os.Stderr, format, a...)
	}
	flag.Usage()
	os.Exit(2)
}

// FatalFlagValue report invalid flag values in same way as flag.Parse().
func FatalFlagValue(msg, name string, val interface{}) {
	FatalUsage("invalid value %#v for flag -%s: %s\n", val, name, msg)
}

// Endpoint validates url and ensure it doesn't have trailing slashes.
func Endpoint(s *string) bool {
	clean := regexp.MustCompile(`/+$`).ReplaceAllString(*s, "")

	p, err := url.Parse(clean)
	if err != nil {
		return false
	}
	if p.Host == "" {
		return false
	}

	*s = clean
	return true
}
