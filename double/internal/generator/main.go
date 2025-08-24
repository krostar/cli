package main

import (
	"context"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		panic(fmt.Sprintf("usage: %s (path to cli sources)", os.Args[0]))
	}

	cliPath := os.Args[1]

	ctx := context.Background()

	command, others, err := findCommandInterfaces(cliPath)
	if err != nil {
		panic(fmt.Sprintf("unable to find all Command* interfaces: %v", err))
	}

	{
		if err := generateInterfaceWrappedImplementation(ctx, command, others, filepath.Join(cliPath, "double", "wrapper_generated.go")); err != nil {
			panic(fmt.Errorf("failed to generate wrapper: %v", err))
		}

		if err := generateInterfaceWrappedTestImplementation(ctx, command, others, filepath.Join(cliPath, "double", "wrapper_generated_test.go")); err != nil {
			panic(fmt.Errorf("failed to generate wrapper tests: %v", err))
		}
	}

	if err := generateFakeImplementation(ctx, command, others, filepath.Join(cliPath, "double", "fake_generated.go")); err != nil {
		panic(fmt.Errorf("failed to generate fake: %v", err))
	}

	if err := generateSpyImplementation(ctx, command, others, filepath.Join(cliPath, "double", "spy_generated.go")); err != nil {
		panic(fmt.Errorf("failed to generate spy: %v", err))
	}
}

func parseDirFiles(fset *token.FileSet, dir string) ([]*ast.File, error) {
	list, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("unable to read directory %s: %w", dir, err)
	}

	var files []*ast.File

	for _, f := range list {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".go") && !strings.HasSuffix(f.Name(), "_test.go") {
			file, err := parser.ParseFile(fset, filepath.Join(dir, f.Name()), nil, parser.SkipObjectResolution)
			if err != nil {
				return nil, fmt.Errorf("unable to parse file: %w", err)
			}

			files = append(files, file)
		}
	}

	return files, nil
}

func findCommandInterfaces(dir string) (*types.Interface, map[string]*types.Interface, error) {
	fset := token.NewFileSet()

	files, err := parseDirFiles(fset, dir)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse directory: %w", err)
	}

	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	pkg := types.NewPackage("github.com/krostar/cli", "cli")
	if err := types.NewChecker(&types.Config{Importer: importer.Default()}, fset, pkg, info).Files(files); err != nil {
		return nil, nil, fmt.Errorf("unable to parse package: %s", err)
	}

	var commandInterface *types.Interface

	otherInterfaces := make(map[string]*types.Interface)

	for _, name := range pkg.Scope().Names() {
		if !strings.HasPrefix(name, "Command") {
			continue
		}

		obj := pkg.Scope().Lookup(name)
		if i, ok := obj.Type().Underlying().(*types.Interface); ok {
			if name == "Command" {
				commandInterface = i
			} else {
				otherInterfaces[name] = i
			}
		}
	}

	if commandInterface == nil {
		return nil, nil, fmt.Errorf("unable to find command interface in %s", dir)
	}

	return commandInterface, otherInterfaces, nil
}
