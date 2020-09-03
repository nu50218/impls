package impls

import (
	"go/types"

	"golang.org/x/tools/go/packages"
)

func objsFromPkgs(isInterface bool, patterns ...string) ([]types.Object, error) {
	mode := packages.NeedSyntax | packages.NeedTypes | packages.NeedDeps | packages.NeedTypesInfo | packages.NeedImports
	cfg := &packages.Config{Mode: mode}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, err
	}

	objs := []types.Object{}
	for _, pkg := range pkgs {
		scoop := pkg.Types.Scope()
		for _, name := range scoop.Names() {
			obj, _ := scoop.Lookup(name).(*types.TypeName)
			if obj == nil {
				continue
			}

			if isInterface && !types.IsInterface(obj.Type()) {
				continue
			}

			objs = append(objs, obj)
		}
	}

	return objs, nil
}

// TypeObjsFromPkgs はパッケージ名の可変長引数で与えられたパッケージで定義されたユーザ定義型の types.Object を取得する
func TypeObjsFromPkgs(patterns ...string) ([]types.Object, error) {
	objs, err := objsFromPkgs(false, patterns...)
	if err != nil {
		return nil, err
	}

	return objs, nil
}

// InterfacesFromPkgs はパッケージ名の可変長引数で与えられたパッケージで定義されたユーザ定義型かつ interface の types.Object を取得する
func InterfacesFromPkgs(patterns ...string) ([]types.Object, error) {
	ifaces, err := objsFromPkgs(true, patterns...)
	if err != nil {
		return nil, err
	}

	return ifaces, nil
}
