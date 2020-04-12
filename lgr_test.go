package lgr

import (
	"testing"
)

func Test_ExtractLevel(t *testing.T) {
	lvl, msg := extractLevel("[INFO] run file server for %s, path %s")
	if lvl != "INFO" {
		t.Errorf("failed to extract log level")
	}
	if msg != "run file server for %s, path %s" {
		t.Errorf("failed to extract message")
	}
}

func Test_Printf(t *testing.T) {
	root := "fake-root"
	path := "fake-path"

	InitZapLogger(true)
	Printf("[INFO] run file server for %s, path %s", root, path)
}