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
	flagIncludeTest  bool
)

var errorIface, _ = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

func init() {
	flagSet.Usage = func() {
		fmt.Printf("Usage of %s:\n", name)
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  $ impls interfaces io.PipeWriter")
		fmt.Println("  $ impls interfaces bytes.Buffer io")
		fmt.Println()
		fmt.Println("Options:")
		flagSet.PrintDefaults()
	}

	flagSet.BoolVar(&flagIncludeError, "e", true, "include error interface (default = true)")
	flagSet.BoolVar(&flagIncludeTest, "t", true, "include test package (default = true)")
}

type c struct{}

func (*c) Name() string {
	return name
}

func (*c) Description() string {
	return "find all interfaces by type"
}

func (*c) FlagSet() *flag.FlagSet {
	return flagSet
}

func typeObjFromName(s string, pkgs []*packages.Package) (types.Object, error) {
	if s == "error" {
		return types.Universe.Lookup("error"), nil
	}

	if !strings.Contains(s, ".") {
		return nil, errors.New("invalid syntax")
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

	obj, err := typeObjFromName(target, pkgs)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		if _, ok := paths[pkg.Types.Path()]; !ok && len(flagArgs) != 1 {
			continue
		}

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
