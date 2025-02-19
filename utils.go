package main

import "strings"

func containsIgnoreCase(s, substr string) bool {
	return (s != "") && (substr != "") && (strings.Contains(strings.ToLower(s), strings.ToLower(substr)))
}

