package iface

import (
	"errors"
	"flag"
	"fmt"
	"go/types"
	"path/filepath"
	"strings"

	"github.com/nu50218/impls/command"
	"github.com/nu50218/impls/impls"
	"golang.org/x/tools/go/packages"
)

const name = "interfaces"

var Command (command.Command) = &c{}

var errorIface = types.Universe.Lookup("error").(*types.TypeName)

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

func typeObjFromName(pkg string, name string) (types.Object, error) {
	mode := packages.NeedSyntax | packages.NeedTypes | packages.NeedDeps | packages.NeedTypesInfo | packages.NeedImports
	cfg := &packages.Config{Mode: mode}
	pkgs, err := packages.Load(cfg, pkg)
	if err != nil {
		return nil, err
	}

	obj := pkgs[0].Types.Scope().Lookup(name)
	if obj == nil {
		return nil, fmt.Errorf("lookup: not found type %s.%s", pkg, name)
	}

	return obj, nil
}

func interfacesCmd(args []string) error {
	if len(args) < 1 || (len(args) < 2 && !flagIncludeError) {
		return errors.New("invalid arguments")
	}

	targetType := args[0]
	searchPkgs := args[1:]

	typ := strings.TrimLeft(filepath.Ext(targetType), ".")
	pkg := strings.TrimRight(strings.TrimSuffix(targetType, typ), ".")
	if pkg == "" || typ == "" {
		return errors.New("invalid type name")
	}

	ifs, err := impls.InterfacesFromPkgs(searchPkgs...)
	if err != nil {
		return err
	}

	if flagIncludeError {
		ifs = append(ifs, errorIface)
	}

	obj, err := typeObjFromName(pkg, typ)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", obj.Type())
	for _, iface := range ifs {
		i, _ := iface.Type().Underlying().(*types.Interface)
		if i == nil {
			return errors.New("invalid interface")
		}
		if types.Implements(obj.Type(), i) || types.Implements(types.NewPointer(obj.Type()), i) {
			fmt.Printf("\t%s\n", iface.Type())
		}
	}

	return nil
}

func (*c) Run(args []string) error {
	if err := flagSet.Parse(args); err != nil {
		return err
	}

	if err := interfacesCmd(args); err != nil {
		return err
	}

	return nil
}
