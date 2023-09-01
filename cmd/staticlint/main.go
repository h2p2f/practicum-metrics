// Custom linter based on staticcheck and an extended standard library,
// which checks the code for errors and suboptimal places in the code.
//
// Build and run: go build -o staticlint cmd/staticlint/main.go
// ./staticlint --help
//
// More about staticcheck: https://staticcheck.io/docs/
// Added analyzers cheknoglobals and nakedret
// connected custom analyzer OsExitInMain
//used Analyzers:
// assign - detects useless assignments
// atomic - detects common mistaken usages of the sync/atomic package
// atomicalign - detects structs that would be smaller if their fields were sorted
// bools - detects common mistakes involving boolean operators
// buildssa - builds an SSA-form program representation
// buildtag - checks that build tags are well-formed
// composite - detects composite literals that can be simplified
// copylock - detects locks erroneously passed by value
// deepequalerrors - detects errors that are compared with reflect.DeepEqual
// defers - detects defers that will never be reached
// directive - checks that directives are followed by a blank line
// errorsas - detects errors that are compared with errors.As
// fieldalignment - detects structs that would be smaller if their fields were sorted
// findcall - finds calls to a function with a given name
// httpresponse - detects common mistakes involving HTTP responses
// ifaceassert - detects redundant type assertions from/to interfaces
// loopclosure - detects references to loop variables from within nested functions
// lostcancel - detects contexts that are canceled too late
// nilfunc - detects useless comparisons between functions and nil
// nilness - detects redundant nil comparisons
// pkgfact - detects packages with large numbers of facts
// printf - detects errors in Printf-style format strings
// reflectvaluecompare - detects values that may be compared with reflect.Value
// shadow - detects shadowed variables
// sigchanyzer - detects common mistaken usages of the sync/atomic package
// slog - detects common mistakes involving the slog package
// sortslice - detects common mistakes involving sorting slices
// stdmethods - detects methods that can be declared on the standard library's types
// stringintconv - detects redundant conversions between strings and integers
// structtag - checks that struct field tags conform to a format string
// testinggoroutine - detects common mistaken usages of the sync/atomic package
// tests - detects common mistakes involving tests
// timeformat - detects invalid time format strings
// unmarshal - detects invalid Unmarshal usage
// unreachable - detects unreachable code
// unsafeptr - detects invalid unsafe.Pointer conversions
// unusedresult - detects unused results of calls to some functions
// unusedwrite - detects unused results of calls to Write
// usesgenerics - detects uses of generics
// osexitinmain - detects os.Exit in main function
// rowserr -  checks whether sql.Rows.Err is correctly checked
// nakedret - detects naked returns

package main

import (
	"strings"

	"github.com/alexkohler/nakedret/v2"
	"github.com/jingyugao/rowserrcheck/passes/rowserr"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/staticcheck"
)

func main() {

	// Slice of analyzers that will be used
	var analyzers []*analysis.Analyzer

	// Add all SA analyzers and one ST from staticcheck
	for _, v := range staticcheck.Analyzers {
		if strings.Contains(v.Analyzer.Name, "SA") || v.Analyzer.Name == "ST1013" {
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	// Add analyzers from the extended standard library
	analyzers = append(analyzers,
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		buildssa.Analyzer,
		buildtag.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		deepequalerrors.Analyzer,
		defers.Analyzer,
		directive.Analyzer,
		errorsas.Analyzer,
		fieldalignment.Analyzer,
		findcall.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		pkgfact.Analyzer,
		printf.Analyzer,
		reflectvaluecompare.Analyzer,
		shadow.Analyzer,
		sigchanyzer.Analyzer,
		slog.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		timeformat.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer,
		usesgenerics.Analyzer,
	)

	// Add our analyzer OsExitInMain
	analyzers = append(analyzers, OsExitInMainAnalyser)

	// Add nakedret analyzer
	analyzers = append(analyzers, nakedret.NakedReturnAnalyzer(5))

	// Add rowserr analyzer
	analyzers = append(analyzers, rowserr.NewAnalyzer("github.com/jmoiron/sqlx"))
	// Run all analyzers
	multichecker.Main(
		analyzers...,
	)
}
