package interpreter

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Parser converts a string expression into an abstract syntax tree
// of Expression objects that can be interpreted
type Parser struct {
	tokens []string
	pos    int
}

// NewParser creates a new parser with the tokenized input
func NewParser(input string) *Parser {
	tokens := tokenize(input)
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

// Parse parses the input and returns the root expression
func (p *Parser) Parse() (Expression, error) {
	if len(p.tokens) == 0 {
		return nil, fmt.Errorf("empty expression")
	}
	return p.parseExpression()
}

// parseExpression parses an expression which can be a sum or difference
func (p *Parser) parseExpression() (Expression, error) {
	// Parse the first term
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	// Continue parsing if there are more tokens
	for p.pos < len(p.tokens) {
		if p.tokens[p.pos] == "+" {
			p.pos++ // Skip the + operator
			right, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			left = NewAddExpression(left, right)
		} else if p.tokens[p.pos] == "-" {
			p.pos++ // Skip the - operator
			right, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			left = NewSubtractExpression(left, right)
		} else {
			// Not a + or - operator, so we're done parsing the expression
			break
		}
	}

	return left, nil
}

// parseTerm parses a term which can be a product or quotient
func (p *Parser) parseTerm() (Expression, error) {
	// Parse the first factor
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	// Continue parsing if there are more tokens
	for p.pos < len(p.tokens) {
		if p.tokens[p.pos] == "*" {
			p.pos++ // Skip the * operator
			right, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			left = NewMultiplyExpression(left, right)
		} else if p.tokens[p.pos] == "/" {
			p.pos++ // Skip the / operator
			right, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			left = NewDivideExpression(left, right)
		} else {
			// Not a * or / operator, so we're done parsing the term
			break
		}
	}

	return left, nil
}

// parseFactor parses a factor which can be a number, variable, function call,
// or a parenthesized expression
func (p *Parser) parseFactor() (Expression, error) {
	if p.pos >= len(p.tokens) {
		return nil, fmt.Errorf("unexpected end of expression")
	}

	token := p.tokens[p.pos]
	p.pos++

	// Check if it's a number
	if value, err := strconv.ParseFloat(token, 64); err == nil {
		return NewNumberExpression(value), nil
	}

	// Check if it's a function
	if p.pos < len(p.tokens) && p.tokens[p.pos] == "(" {
		p.pos++ // Skip the opening parenthesis
		arg, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		if p.pos >= len(p.tokens) || p.tokens[p.pos] != ")" {
			return nil, fmt.Errorf("expected closing parenthesis")
		}
		p.pos++ // Skip the closing parenthesis

		return NewFunctionExpression(token, arg)
	}

	// Check if it's an opening parenthesis
	if token == "(" {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		if p.pos >= len(p.tokens) || p.tokens[p.pos] != ")" {
			return nil, fmt.Errorf("expected closing parenthesis")
		}
		p.pos++ // Skip the closing parenthesis

		return expr, nil
	}

	// Otherwise it's a variable
	return NewVariableExpression(token), nil
}

// tokenize breaks the input string into tokens
func tokenize(input string) []string {
	input = strings.TrimSpace(input)
	var tokens []string
	var buffer strings.Builder

	i := 0
	for i < len(input) {
		char := rune(input[i])

		// Skip whitespace
		if unicode.IsSpace(char) {
			i++
			continue
		}

		// Handle operators and parentheses
		if char == '+' || char == '-' || char == '*' || char == '/' || char == '(' || char == ')' {
			tokens = append(tokens, string(char))
			i++
			continue
		}

		// Handle numbers and identifiers
		buffer.Reset()
		if unicode.IsDigit(char) {
			// Parse number
			for i < len(input) && (unicode.IsDigit(rune(input[i])) || input[i] == '.') {
				buffer.WriteByte(input[i])
				i++
			}
		} else if unicode.IsLetter(char) {
			// Parse identifier
			for i < len(input) && (unicode.IsLetter(rune(input[i])) || unicode.IsDigit(rune(input[i]))) {
				buffer.WriteByte(input[i])
				i++
			}
		} else {
			// Unexpected character
			i++
			continue
		}

		tokens = append(tokens, buffer.String())
	}

	return tokens
}
