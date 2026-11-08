package main

import (
	"flag"
	"fmt"
	urlpkg "net/url"
	"os"
	"regexp"
)

const exitCodeUsage = 2 // Same as used by flag.Parse.

// FatalUsagef report usage error in same way as [flag.Parse].
func FatalUsagef(format string, a ...any) {
	if format != "" {
		fmt.Fprintf(os.Stderr, format, a...)
	}
	flag.Usage()
	os.Exit(exitCodeUsage) //nolint:revive // Exit in same way as flag.Parse.
}

// FatalFlagValue report invalid flag values in same way as [flag.Parse].
func FatalFlagValue(msg, name string, val any) {
	FatalUsagef("invalid value %#v for flag -%s: %s\n", val, name, msg)
}

// Endpoint validates url and ensure it doesn't have trailing slashes.
func Endpoint(s *string) bool {
	clean := regexp.MustCompile(`/+$`).ReplaceAllString(*s, "")

	p, err := urlpkg.Parse(clean)
	if err != nil {
		return false
	}
	if p.Host == "" {
		return false
	}

	*s = clean
	return true
}
