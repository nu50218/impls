package impls

import (
	"errors"
	"go/types"

	"golang.org/x/tools/go/packages"
)

func LoadPkgs(patterns ...string) ([]*packages.Package, error) {
	mode := packages.NeedSyntax | packages.NeedTypes | packages.NeedDeps | packages.NeedTypesInfo | packages.NeedImports
	cfg := &packages.Config{Mode: mode}
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
