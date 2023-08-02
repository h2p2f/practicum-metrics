package main

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func TestOsExitAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), OsExitInMainAnalyser, "./...")

}
