package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
)

func main() {

	ctx := context.Background()

	// Construct a Rego object that can be prepared or evaluated.

	r := rego.New(
		rego.Query(os.Args[2]),
		// rego.Query(`x = hello("bob")`),
		rego.Function1(
			&rego.Function{
				Name: "hello",
				Decl: types.NewFunction(types.Args(types.S), types.S),
			},
			func(_ rego.BuiltinContext, a *ast.Term) (*ast.Term, error) {
				if str, ok := a.Value.(ast.String); ok {
					return ast.StringTerm("hello, " + string(str)), nil
				}
				return nil, nil
			}),
		rego.Load([]string{os.Args[1]}, nil),
	)

	// r := rego.New(
	// 	rego.Query(os.Args[2]),
	// 	rego.Load([]string{os.Args[1]}, nil))

	// Create a prepared query that can be evaluated.
	query, err := r.PrepareForEval(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Load the input document from stdin.
	var input interface{}
	dec := json.NewDecoder(os.Stdin)
	dec.UseNumber()
	if err := dec.Decode(&input); err != nil {
		log.Fatal(err)
	}

	// Execute the prepared query.
	rs, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		log.Fatal(err)
	}

	// Do something with the result.
	fmt.Println(rs)
}
