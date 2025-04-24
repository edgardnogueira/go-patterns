# Interpreter Pattern

## Intent
The Interpreter pattern defines a grammatical representation for a language and an interpreter to deal with this grammar. It is used to interpret context in a given language.

## Problem
Given a language, define a representation for its grammar along with an interpreter that uses the representation to interpret sentences in the language.

## Solution
The Interpreter pattern suggests modeling the grammar as a class hierarchy and implementing the interpret operation in these classes. Each rule in the grammar becomes a class, and expressions can be composed of terminal and non-terminal expressions.

## Structure
- **AbstractExpression**: Declares the interpret operation that all nodes in the abstract syntax tree must implement.
- **TerminalExpression**: Implements the interpret operation for terminal symbols in the grammar.
- **NonterminalExpression**: Implements the interpret operation for non-terminal symbols in the grammar, typically containing one or more child expressions.
- **Context**: Contains information that is global to the interpreter, such as variable mappings.
- **Client**: Builds (or is given) the abstract syntax tree representing a particular sentence in the language.

## Implementation
In this implementation, we create a simple mathematical expression evaluator that can handle:
- Basic numeric expressions
- Variables
- Arithmetic operations (+, -, *, /)
- Functions (e.g., sin, cos, sqrt)

## When to use
- The grammar is simple and can be represented as an abstract syntax tree.
- You need to interpret frequently occurring expressions in a well-defined domain.
- You want to create a domain-specific language (DSL).
- Performance is not a critical concern (interpreters often don't provide the best performance).

## Benefits
- Easy to add new expressions.
- Easy to change or extend the grammar.
- Implementing the grammar is straightforward.

## Drawbacks
- Complex grammars lead to complex class hierarchies.
- May be inefficient for large expressions.
- Maintenance can become difficult for complex grammars.

## Go-Specific Implementation Notes
In Go, we implement the pattern using interfaces and composition. Instead of creating a class hierarchy (as would be typical in languages like Java or C++), we define a common interface for all expressions and implement it in various struct types.
