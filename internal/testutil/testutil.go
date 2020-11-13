// Package testutil contains assistive utilities for testing.
package testutil

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/AdguardTeam/golibs/log"
)

// DiscardLogOutput runs tests with discarded logger output.
func DiscardLogOutput(m *testing.M) {
	// TODO(e.burkov): Using of global mutable logger is temporary solution.
	log.SetOutput(ioutil.Discard)

	os.Exit(m.Run())
}
