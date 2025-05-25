package main

import (
	"go/types"
)

// getTupleRepresentation transforms a types.Tuple (representing function parameters or results)
// into three parallel slices: variable names, type names, and "variable type" strings.
//
//   - vars: Slice of variable names (e.g., ["a0", "b0", "c0"])
//   - types: Slice of type types (e.g., ["string", "int", "context.Context"])
//   - varsTypes: Slice of "variable type" strings (e.g., ["a0 string", "b0 int", "c0 context.Context"])
//
// This function is used to generate function signatures in the generated code,
// creating sequential variable names and proper type representations.
func getTupleRepresentation(params *types.Tuple, imports map[string]string, nbVariables int) ([]string, []string, []string) {
	var (
		vars      []string
		typs      []string
		varsTypes []string
	)

	for i := range params.Len() {
		paramType := params.At(i).Type()
		fillUsedImports(paramType, imports)

		typ := types.TypeString(paramType, packageQualifier(imports))

		variable := string([]byte{
			byte('a' + (i+(nbVariables%26))%26),
			byte('0' + (i+(nbVariables/26))/26),
		})

		vars = append(vars, variable)
		typs = append(typs, typ)
		varsTypes = append(varsTypes, variable+" "+typ)
	}

	return vars, typs, varsTypes
}

// generateCombinations creates all possible non-empty subsets of the provided set.
// Adapted from: https://github.com/mxschmitt/golang-combinations/blob/main/combinations.go
func generateCombinations[T any](set []T) [][]T {
	var combinations [][]T
	length := uint(len(set))

	// Go through all possible combinations of objects
	// from 1 (only first object in subset) to 2^length-1 (all objects in subset)
	for subsetBits := 1; subsetBits < (1 << length); subsetBits++ {
		var subset []T

		for object := range length {
			if (subsetBits>>object)&1 == 1 {
				subset = append(subset, set[object])
			}
		}
		combinations = append(combinations, subset)
	}

	return combinations
}

// fillUsedImports recursively analyzes a type to identify all packages that need to be imported.
// It modifies the provided imports map to include all required packages.
//
// This function recursively examines composite types (slices, maps, etc.) to find
// all named types that will require imports in the generated code.
func fillUsedImports(t types.Type, imports map[string]string) {
	switch t := t.(type) {
	case *types.Named: // for named types, add the package to imports if it exists
		obj := t.Obj()
		pkg := obj.Pkg()
		if pkg != nil && pkg.Path() != "" {
			imports[pkg.Path()] = pkg.Name()
		}

		if t.TypeArgs() != nil { // handle generic types by recursively processing type arguments
			args := t.TypeArgs()
			for i := range args.Len() {
				fillUsedImports(args.At(i), imports)
			}
		}

	case *types.Pointer: // for pointers, process the element type
		fillUsedImports(t.Elem(), imports)

	case *types.Slice: // for slices, process the element type
		fillUsedImports(t.Elem(), imports)

	case *types.Array: // for arrays, process the element type
		fillUsedImports(t.Elem(), imports)

	case *types.Map: // for maps, process both key and value types
		fillUsedImports(t.Key(), imports)
		fillUsedImports(t.Elem(), imports)

	case *types.Chan: // for channels, process the element type
		fillUsedImports(t.Elem(), imports)

	case *types.Signature: // for function signatures, process all parameter and result types
		params := t.Params()
		for i := range params.Len() {
			fillUsedImports(params.At(i).Type(), imports)
		}
		results := t.Results()
		for i := range results.Len() {
			fillUsedImports(results.At(i).Type(), imports)
		}
	}
}

// packageQualifier creates a types.Qualifier function that maps package paths to their names.
func packageQualifier(imports map[string]string) types.Qualifier {
	return func(p *types.Package) string {
		return imports[p.Path()]
	}
}
