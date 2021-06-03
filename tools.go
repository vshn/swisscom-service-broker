// +build tools

// Package tools contains runtime dependencies as imports.
// Add such runtime dependencies top this file.
// Go modules will be forced to download and install them.
// It also keeps track of the versions.
package tools

import (
	// Required for the e2e-test
	_ "sigs.k8s.io/kustomize/kustomize/v3"
)
