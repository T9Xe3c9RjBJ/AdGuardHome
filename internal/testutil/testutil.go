// Package testutil contains the function to be wrapped by package's TestMain
// to discard logger output.
package testutil

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/AdguardTeam/golibs/log"
)

// TestMain discards logger output and runs tests.
func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)

	os.Exit(m.Run())
}
