package main

import "os"

// resolveExistingPath returns the first path that exists on disk.
func resolveExistingPath(candidates ...string) string {
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}
