package iface

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

const name = "interfaces"

var Command (command.Command) = &c{}

// var errorIface = types.Universe.Lookup("error").(*types.TypeName)

var (
	flagSet          = flag.NewFlagSet(name, flag.ExitOnError)
	flagIncludeError bool
)

func init() {
	flagSet.BoolVar(&flagIncludeError, "e", true, "include error interface (default = true)")
}

type c struct{}

func (*c) Name() string {
	return name
}

func typeObjFromName(s string, pkgs []*packages.Package) (types.Object, error) {
	comma := strings.LastIndex(s, ".")
	path := s[:comma]
	name := s[comma+1:]

	buildPkg, err := build.Default.Import(path, ".", build.ImportMode(0))
	if err != nil {
		return nil, err
	}

	for _, pkg := range pkgs {
		obj := pkg.Types.Scope().Lookup(name)
		if obj.Pkg().Path() != buildPkg.ImportPath {
			continue
		}

		return obj, nil
	}

	return nil, errors.New("not found")
}

func interfacesCmd(args []string) error {
	if len(args) < 1 || (len(args) < 2 && !flagIncludeError) {
		return errors.New("invalid arguments")
	}

	target := args[0]
	loadPkgs := append(args[1:], target[:strings.LastIndex(target, ".")])
	pkgs, err := impls.LoadPkgs(loadPkgs...)
	if err != nil {
		return err
	}

	obj, err := typeObjFromName(target, pkgs)
	if err != nil {
		return err
	}

	// if flagIncludeError {
	// 	ifs = append(ifs, errorIface)
	// }

	for _, pkg := range pkgs {
		scoop := pkg.Types.Scope()
		for _, name := range scoop.Names() {
			iface, _ := scoop.Lookup(name).(*types.TypeName)
			if iface == nil {
				continue
			}

			if !types.IsInterface(iface.Type()) {
				continue
			}

			i, err := impls.UnderlyingInterface(iface.Type())
			if err != nil {
				return err
			}

			if impls.Implements(obj.Type(), i) {
				fmt.Printf("%s %s.%s\n", pkg.Fset.Position(iface.Pos()), pkg.Types.Name(), iface.Name())
			}
		}
	}

	return nil
}

func (*c) Run(args []string) error {
	if err := flagSet.Parse(args); err != nil {
		return err
	}

	if err := interfacesCmd(flagSet.Args()); err != nil {
		return err
	}

	return nil
}
