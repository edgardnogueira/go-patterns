package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/edgardnogueira/go-patterns/behavioral/interpreter"
)

func main() {
	fmt.Println("==== Math Expression Interpreter ====")
	fmt.Println("Enter expressions to evaluate. Type 'exit' to quit.")
	fmt.Println("You can use variables (e.g., 'x + y') and functions (sin, cos, sqrt, log, abs).")
	fmt.Println("To set a variable, use: 'let x = 5'")
	fmt.Println()

	ctx := interpreter.NewContext()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "exit" {
			break
		}

		// Check if it's a variable assignment
		if strings.HasPrefix(input, "let ") {
			parts := strings.SplitN(strings.TrimPrefix(input, "let "), "=", 2)
			if len(parts) != 2 {
				fmt.Println("Invalid variable assignment syntax. Use: 'let variable = value'")
				continue
			}

			varName := strings.TrimSpace(parts[0])
			varExpr := strings.TrimSpace(parts[1])

			parser := interpreter.NewParser(varExpr)
			expr, err := parser.Parse()
			if err != nil {
				fmt.Printf("Error parsing expression: %v\n", err)
				continue
			}

			val, err := expr.Interpret(ctx)
			if err != nil {
				fmt.Printf("Error evaluating expression: %v\n", err)
				continue
			}

			ctx.SetVariable(varName, val)
			fmt.Printf("Set %s = %v\n", varName, val)
			continue
		}

		// Parse and evaluate the expression
		parser := interpreter.NewParser(input)
		expr, err := parser.Parse()
		if err != nil {
			fmt.Printf("Error parsing expression: %v\n", err)
			continue
		}

		result, err := expr.Interpret(ctx)
		if err != nil {
			fmt.Printf("Error evaluating expression: %v\n", err)
			continue
		}

		fmt.Printf("Result: %v\n", result)
		fmt.Printf("Expression tree: %s\n", expr.String())
	}
}
