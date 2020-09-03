package types

import (
	"errors"
	"flag"
	"fmt"
	"go/build"
	"go/types"
	"strings"

	"github.com/nu50218/impls/command"
	"github.com/nu50218/impls/impls"
	"golang.org/x/tools/go/packages"
)

const name = "types"

var Command (command.Command) = &c{}

var flagSet = flag.NewFlagSet(name, flag.ExitOnError)

// オプション
var ()

type c struct{}

func (*c) Name() string {
	return name
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
	loadPkgs := append(flagArgs[1:], target[:strings.LastIndex(target, ".")])
	pkgs, err := impls.LoadPkgs(loadPkgs...)
	if err != nil {
		return err
	}

	i, err := findInterface(target, pkgs)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		scope := pkg.Types.Scope()
		for _, n := range scope.Names() {
			obj := scope.Lookup(n)
			if implements(obj.Type(), i) {
				fmt.Println(pkg.Fset.Position(obj.Pos()), pkg.Types.Name()+"."+obj.Name())
			}
		}
	}

	return nil
}

func findInterface(s string, pkgs []*packages.Package) (*types.Interface, error) {
	lastComma := strings.LastIndex(s, ".")
	ifacePath := s[:lastComma]
	ifaceName := s[lastComma+1:]

	buildPkg, err := build.Default.Import(ifacePath, ".", build.ImportMode(0))
	if err != nil {
		return nil, err
	}

	for _, pkg := range pkgs {
		scope := pkg.Types.Scope()
		for _, n := range scope.Names() {
			obj := scope.Lookup(n)
			if obj.Name() != ifaceName || obj.Pkg().Path() != buildPkg.ImportPath {
				continue
			}
			i, err := impls.UnderlyingInterface(obj.Type())
			if err != nil {
				return nil, err
			}
			return i, nil
		}
	}

	return nil, errors.New("not found")
}

func implements(V types.Type, T *types.Interface) bool {
	if types.Implements(V, T) {
		return true
	}

	if types.Implements(types.NewPointer(V), T) {
		return true
	}

	return false
}
