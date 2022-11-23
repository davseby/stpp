package analyzer

import (
	"flag"
	"go/ast"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

var flagSet flag.FlagSet

var (
	prefix string
)

func init() {
	flagSet.StringVar(&prefix, "prefix", "_", "what prefix to expect")
}

func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:  "fixedglobals",
		Doc:   "checks that all private global variables contains predefined prefix",
		Run:   run,
		Flags: flagSet,
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, func(node ast.Node) bool {
			_, ok := node.(*ast.FuncDecl)
			if ok {
				return false
			}

			v, ok := node.(*ast.ValueSpec)
			if !ok {
				return true
			}

			for _, name := range v.Names {
				if unicode.IsUpper(rune(name.Name[0])) {
					continue
				}

				if !strings.HasPrefix(name.Name, prefix) {
					pass.Reportf(node.Pos(), "%q name is not prefixed with %q", name.Name, f.Name.Name)
				}
			}

			return true
		})
	}

	return nil, nil
}
