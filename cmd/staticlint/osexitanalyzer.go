// Package staticlint содержит в себе анализатор, который проверяет наличие вызова os.Exit() в main модуле.
//
// package staticlint contains an analyzer that checks for the presence of an os.Exit() call in the main module.
package main

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/go/types/typeutil"
)

// osExitInMainAnalyser - анализатор, который проверяет наличие вызова os.Exit() в main модуле.
// osExitInMainAnalyser - analyzer that checks for the presence of an os.Exit() call in the main module.
var OsExitInMainAnalyser = &analysis.Analyzer{
	Name: "osExitInMain",
	Doc:  "os.Exit() in main module check",
	Run:  run,
}

// run - функция, которая проверяет наличие вызова os.Exit() в main модуле.
// run - function that checks for the presence of an os.Exit() call in the main module.
func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}
	for _, fun := range pass.Files {
		ast.Inspect(fun, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.CallExpr:
				sel, ok := n.Fun.(*ast.SelectorExpr)
				if !ok {
					return false
				}
				fn, ok := pass.TypesInfo.ObjectOf(sel.Sel).(*types.Func)
				if !ok {
					return false
				}
				if typeutil.FuncName(fn) == "os.Exit" {

					pass.Reportf(n.Pos(), "os.Exit() in main module")
				}
			}
			return true
		})
	}
	return nil, nil
}
