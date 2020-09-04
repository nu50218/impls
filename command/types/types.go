package types

import (
	"errors"
	"flag"
	"fmt"
	"go/build"
	"go/types"
	"reflect"
	"strings"

	"github.com/nu50218/impls/command"
	"github.com/nu50218/impls/impls"
	"golang.org/x/tools/go/packages"
)

const name = "types"

var Command (command.Command) = &c{}

var flagSet = flag.NewFlagSet(name, flag.ExitOnError)

// オプション
var (
	exported        bool
	flagIncludeTest bool
)

func init() {
	flagSet.Usage = func() {
		fmt.Printf("Usage of %s:\n", name)
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  $ impls types go/ast.Expr")
		fmt.Println("  $ impls types io.Writer net/http")
		fmt.Println()
		fmt.Println("Options:")
		flagSet.PrintDefaults()
	}

	flagSet.BoolVar(&exported, "exported", false, "only exported")
	flagSet.BoolVar(&flagIncludeTest, "t", false, "include test package")
}

type c struct{}

func (*c) Name() string {
	return name
}

func (*c) Description() string {
	return "find all types by interface"
}

func (*c) FlagSet() *flag.FlagSet {
	return flagSet
}

// Run 下のような感じで動かしたい
// $ impls types [options] io.Writer
// $ impls types [options] io.Writer io
// $ impls types [options] io.Writer ./... io
func (*c) Run(args []string) error {
	if err := flagSet.Parse(args); err != nil {
		return err
	}

	flagArgs := flagSet.Args()
	if len(flagArgs) < 1 {
		return errors.New("invalid arguments")
	}

	target := flagArgs[0]
	var targetPkgPath string
	loadPkgs := flagArgs[1:]
	if strings.Contains(target, ".") {
		targetPkgPath = target[:strings.LastIndex(target, ".")]
		loadPkgs = append(loadPkgs, targetPkgPath)
	}

	type pkgPathsResponse struct {
		paths map[string]struct{}
		err   error
	}

	pkgPathsChan := make(chan *pkgPathsResponse, 1)
	go func() {
		paths, err := impls.PkgPaths(flagArgs[1:]...)
		pkgPathsChan <- &pkgPathsResponse{paths: paths, err: err}
	}()

	pkgs, err := impls.LoadPkgs(flagIncludeTest, loadPkgs...)
	if err != nil {
		return err
	}

	checkPkgIncludedResp := <-pkgPathsChan
	paths, err := checkPkgIncludedResp.paths, checkPkgIncludedResp.err
	if err != nil {
		return err
	}

	i, err := findInterface(target, pkgs)
	if err != nil {
		return err
	}

	covered := map[string]struct{}{}
	for _, pkg := range pkgs {
		if _, ok := paths[pkg.Types.Path()]; !ok && len(flagArgs) != 1 {
			continue
		}

		scope := pkg.Types.Scope()
		for _, n := range scope.Names() {
			obj := scope.Lookup(n)
			if exported && !obj.Exported() {
				continue
			}

			t, _ := obj.(*types.TypeName)
			if t == nil {
				continue
			}

			pos := pkg.Fset.Position(obj.Pos()).String()
			if _, ok := covered[pos]; ok {
				continue
			}

			if impls.Implements(obj.Type(), i) {
				covered[pos] = struct{}{}
				fmt.Println(pos, pkg.Types.Name()+"."+obj.Name())
			}
		}
	}

	return nil
}

func findInterface(s string, pkgs []*packages.Package) (*types.Interface, error) {
	if s == "error" {
		errType, _ := types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
		return errType, nil
	}

	if !strings.Contains(s, ".") {
		return nil, errors.New("invalid syntax")
	}

	lastComma := strings.LastIndex(s, ".")
	ifacePath := s[:lastComma]
	ifaceName := s[lastComma+1:]

	buildPkg, err := build.Default.Import(ifacePath, ".", build.ImportMode(0))
	if err != nil {
		return nil, err
	}

	for _, pkg := range pkgs {
		obj := pkg.Types.Scope().Lookup(ifaceName)
		if obj == nil || reflect.ValueOf(obj).IsNil() || obj.Pkg().Path() != buildPkg.ImportPath {
			continue
		}
		i, err := impls.UnderlyingInterface(obj.Type())
		if err != nil {
			return nil, err
		}
		return i, nil
	}

	return nil, errors.New("not found")
}
