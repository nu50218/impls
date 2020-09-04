package impls

import (
	"errors"
	"go/types"

	"golang.org/x/tools/go/packages"
)

func PkgPaths(patterns ...string) (map[string]struct{}, error) {
	mode := packages.NeedDeps | packages.NeedTypes
	cfg := &packages.Config{Mode: mode}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, err
	}
	var res = make(map[string]struct{}, len(pkgs))
	for _, pkg := range pkgs {
		res[pkg.Types.Path()] = struct{}{}
	}
	return res, nil
}

func LoadPkgs(incTest bool, patterns ...string) ([]*packages.Package, error) {
	mode := packages.NeedSyntax | packages.NeedTypes | packages.NeedDeps | packages.NeedTypesInfo | packages.NeedImports
	cfg := &packages.Config{Mode: mode, Tests: incTest}
	return packages.Load(cfg, patterns...)
}

func UnderlyingInterface(t types.Type) (*types.Interface, error) {
	switch t := t.(type) {
	case *types.Interface:
		return t, nil
	case *types.Named:
		return UnderlyingInterface(t.Underlying())
	default:
		return nil, errors.New("not interface")
	}
}

func Implements(V types.Type, T *types.Interface) bool {
	if types.Implements(V, T) {
		return true
	}

	if types.Implements(types.NewPointer(V), T) {
		return true
	}

	return false
}
