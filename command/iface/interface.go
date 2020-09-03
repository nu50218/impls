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

var (
	flagSet          = flag.NewFlagSet(name, flag.ExitOnError)
	flagIncludeError bool
)

var errorIface, _ = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

func init() {
	flagSet.BoolVar(&flagIncludeError, "e", true, "include error interface (default = true)")
}

type c struct{}

func (*c) Name() string {
	return name
}

func typeObjFromName(s string, pkgs []*packages.Package) (types.Object, error) {
	if s == "error" {
		return types.Universe.Lookup("error"), nil
	}

	comma := strings.LastIndex(s, ".")
	path := s[:comma]
	name := s[comma+1:]

	buildPkg, err := build.Default.Import(path, ".", build.ImportMode(0))
	if err != nil {
		return nil, err
	}

	for _, pkg := range pkgs {
		obj, _ := pkg.Types.Scope().Lookup(name).(*types.TypeName)
		if obj == nil {
			continue
		}
		if obj.Pkg().Path() == buildPkg.ImportPath {
			return obj, nil
		}
	}

	return nil, errors.New("not found")
}

func interfacesCmd(args []string) error {
	if len(args) < 1 {
		return errors.New("invalid arguments")
	}

	target := args[0]
	loadPkgs := args[1:]
	if strings.Contains(target, ".") {
		loadPkgs = append(loadPkgs, target[:strings.LastIndex(target, ".")])
	}
	pkgs, err := impls.LoadPkgs(loadPkgs...)
	if err != nil {
		return err
	}

	obj, err := typeObjFromName(target, pkgs)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		scoop := pkg.Types.Scope()
		for _, name := range scoop.Names() {
			iface, _ := scoop.Lookup(name).(*types.TypeName)
			if iface == nil {
				continue
			}

			i, err := impls.UnderlyingInterface(iface.Type())
			if err != nil {
				continue
			}

			if impls.Implements(obj.Type(), i) {
				fmt.Printf("%s %s.%s\n", pkg.Fset.Position(iface.Pos()), pkg.Types.Name(), iface.Name())
			}
		}
	}

	if flagIncludeError && impls.Implements(obj.Type(), errorIface) {
		fmt.Println("error")
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
