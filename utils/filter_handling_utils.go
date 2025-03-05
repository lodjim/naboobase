package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/expr-lang/expr"
	"go.mongodb.org/mongo-driver/bson"
)

var filterLexer = lexer.MustSimple([]lexer.SimpleRule{
	{"Whitespace", `\s+`},
	{"Operator", `&&|\|\||=|\!=|>|>=|<|<=|~|\!~`},
	{"Arithmetic", `\+|\-|\*|\/`},
	{"Punct", `[\(\):]`},
	{"String", `"[^"]*"|'[^']*'`},
	{"Number", `[0-9]+(\.[0-9]+)?([eE][+-]?[0-9]+)?`},
	{"Bool", `true|false`},
	{"Null", `null`},
	{"Identifier", `[a-zA-Z_][a-zA-Z0-9_]*`},
	{"Macro", `@[a-zA-Z_][a-zA-Z0-9_]*`},
})

type Expression struct {
	Logical *LogicalExpression `@@`
}

type LogicalExpression struct {
	Left     *ComparisonExpression `@@`
	Operator string                `( @("&&" | "||")`
	Right    *LogicalExpression    `@@ )?`
}

type ComparisonExpression struct {
	Pos      lexer.Position
	Field    *Identifier      `@@`
	Operator string           `@("=" | "!=" | ">" | ">=" | "<" | "<=" | "~" | "!~")`
	Value    *ValueExpression `@@`
	EndPos   lexer.Position
}

type Identifier struct {
	Name     string  `@Identifier`
	Modifier *string `(":" @Identifier)?`
}

type ValueExpression struct {
	Additive *AdditiveValueExpression `@@`
}

type AdditiveValueExpression struct {
	Pos    lexer.Position
	Left   *MultiplicativeValueExpression   `@@`
	Ops    []string                         `( @("+" | "-")`
	Rights []*MultiplicativeValueExpression `@@ )*`
}

type MultiplicativeValueExpression struct {
	Left   *PrimaryValueExpression   `@@`
	Ops    []string                  `( @("*" | "/")`
	Rights []*PrimaryValueExpression `@@ )*`
}

type PrimaryValueExpression struct {
	String  *string          `@String`
	Number  *float64         `| @Number`
	Bool    *bool            `| @("true" | "false")`
	Null    *string          `| @"null"`
	Macro   *string          `| @Macro`
	SubExpr *ValueExpression `| "(" @@ ")"`
}

// Parser instance
var parser = participle.MustBuild[Expression](
	participle.Lexer(filterLexer),
	participle.Elide("Whitespace"),
)

// TransformFilterToMongoQuery transforms a filter string to a MongoDB query
func TransformFilterToMongoQuery(filter string) (bson.M, error) {
	// Parse the filter string into an AST
	ast, err := parser.ParseString("", filter)
	if err != nil {
		return nil, fmt.Errorf("failed to parse filter: %v", err)
	}

	// Build the MongoDB query
	query, err := buildQuery(ast.Logical, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}
	return query, nil
}

// buildQuery constructs the MongoDB query from a LogicalExpression
func buildQuery(logical *LogicalExpression, filter string) (bson.M, error) {
	if logical.Operator == "" {
		return buildComparisonQuery(logical.Left, filter)
	}

	op := "$and"
	if logical.Operator == "||" {
		op = "$or"
	}

	leftQuery, err := buildComparisonQuery(logical.Left, filter)
	if err != nil {
		return nil, err
	}

	rightQuery, err := buildQuery(logical.Right, filter)
	if err != nil {
		return nil, err
	}

	return bson.M{op: []bson.M{leftQuery, rightQuery}}, nil
}

// buildComparisonQuery constructs a MongoDB query from a ComparisonExpression
func buildComparisonQuery(comp *ComparisonExpression, filter string) (bson.M, error) {
	field := comp.Field.Name
	modifier := ""
	if comp.Field.Modifier != nil {
		modifier = *comp.Field.Modifier
	}

	// Find the end position for the value expression
	// If EndPos is not set, we need to find the end of this expression
	// This could be until the next operator (&&, ||) or end of string
	valueExprStr := ""
	if comp.Value != nil {
		startPos := comp.Value.Additive.Pos.Offset
		endPos := comp.EndPos.Offset

		// If the end position is invalid, find the next logical operator or use end of filter
		if endPos == 0 || endPos <= startPos {
			// Look for the next && or || after the start position
			andPos := strings.Index(filter[startPos:], "&&")
			orPos := strings.Index(filter[startPos:], "||")

			endPos = len(filter)
			if andPos >= 0 && (orPos < 0 || andPos < orPos) {
				endPos = startPos + andPos
			} else if orPos >= 0 {
				endPos = startPos + orPos
			}
		}

		valueExprStr = strings.TrimSpace(filter[startPos:endPos])
	}

	// Evaluate the value expression into an interface{} value
	value, err := evaluateValue(comp.Value, valueExprStr)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate value '%s': %v", valueExprStr, err)
	}

	// Map operators to MongoDB operators
	switch comp.Operator {
	case "=":
		if modifier == "lower" {
			if str, ok := value.(string); ok {
				return bson.M{field: bson.M{"$regex": "^" + regexp.QuoteMeta(str) + "$", "$options": "i"}}, nil
			}
			return bson.M{field: value}, nil
		} else if modifier == "length" {
			if num, ok := value.(float64); ok {
				return bson.M{field: bson.M{"$size": int(num)}}, nil
			}
			return nil, errors.New("length modifier requires a number")
		}
		return bson.M{field: value}, nil
	case "!=":
		return bson.M{field: bson.M{"$ne": value}}, nil
	case ">":
		return bson.M{field: bson.M{"$gt": value}}, nil
	case ">=":
		return bson.M{field: bson.M{"$gte": value}}, nil
	case "<":
		return bson.M{field: bson.M{"$lt": value}}, nil
	case "<=":
		return bson.M{field: bson.M{"$lte": value}}, nil
	case "~":
		if str, ok := value.(string); ok {
			regex := convertPatternToRegex(str)
			if modifier == "lower" {
				return bson.M{field: bson.M{"$regex": regex, "$options": "i"}}, nil
			}
			return bson.M{field: bson.M{"$regex": regex}}, nil
		}
		return nil, errors.New("~ operator requires a string pattern")
	case "!~":
		if str, ok := value.(string); ok {
			regex := convertPatternToRegex(str)
			if modifier == "lower" {
				return bson.M{field: bson.M{"$not": bson.M{"$regex": regex, "$options": "i"}}}, nil
			}
			return bson.M{field: bson.M{"$not": bson.M{"$regex": regex}}}, nil
		}
		return nil, errors.New("!~ operator requires a string pattern")
	default:
		return nil, fmt.Errorf("unsupported operator: %s", comp.Operator)
	}
}

func evaluateValue(valExpr *ValueExpression, exprStr string) (interface{}, error) {
	// Handle simple literal values directly from the AST
	primary := valExpr.Additive.Left.Left
	if primary.String != nil {
		return strings.Trim(*primary.String, `"'`), nil
	}
	if primary.Number != nil {
		return *primary.Number, nil
	}
	if primary.Bool != nil {
		return *primary.Bool, nil
	}
	if primary.Null != nil {
		return nil, nil
	}
	if primary.Macro != nil {
		return evaluateMacro(*primary.Macro)
	}
	if primary.SubExpr != nil {
		return evaluateExpr(exprStr)
	}

	// Check if we have a complex expression (with operators)
	hasOps := len(valExpr.Additive.Ops) > 0
	if hasOps || exprStr != "" {
		return evaluateExpr(exprStr)
	}

	return nil, fmt.Errorf("unable to evaluate value expression")
}

// evaluateMacro handles macro values like @now
func evaluateMacro(macro string) (interface{}, error) {
	switch macro {
	case "@now":
		return time.Now(), nil
	default:
		return nil, fmt.Errorf("unsupported macro: %s", macro)
	}
}

// evaluateExpr evaluates arithmetic expressions using the expr library
func evaluateExpr(exprStr string) (interface{}, error) {
	// Clean up the expression string
	exprStr = strings.TrimSpace(exprStr)

	// If it's a simple macro reference, handle directly
	if strings.HasPrefix(exprStr, "@") {
		return evaluateMacro(exprStr)
	}

	// Define environment with macros
	env := map[string]interface{}{
		"now": func() float64 {
			return float64(time.Now().UnixNano()) / 1e6 // milliseconds
		},
	}

	// Preprocess expression to handle @macros
	if strings.Contains(exprStr, "@now") {
		exprStr = strings.ReplaceAll(exprStr, "@now", "now()")
	}

	// Parse and evaluate the expression
	program, err := expr.Compile(exprStr, expr.Env(env))
	if err != nil {
		return nil, fmt.Errorf("failed to compile expression '%s': %v", exprStr, err)
	}

	result, err := expr.Run(program, env)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate expression '%s': %v", exprStr, err)
	}

	switch v := result.(type) {
	case float64:
		if v > 1e12 { // Assuming milliseconds since epoch
			return time.Unix(0, int64(v*1e6)), nil
		}
		return v, nil
	case int:
		return float64(v), nil
	case string:
		return strings.Trim(v, `"'`), nil
	case bool:
		return v, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported value type: %T", v)
	}
}

// convertPatternToRegex converts a SQL-like pattern to a regex pattern
func convertPatternToRegex(pattern string) string {
	pattern = strings.Trim(pattern, `"'`)
	if pattern == "%" {
		return ".*"
	}
	pattern = regexp.QuoteMeta(pattern)
	pattern = strings.ReplaceAll(pattern, "%", ".*")
	pattern = strings.ReplaceAll(pattern, "_", ".")
	return "^" + pattern + "$"
}

/*
func main() {
	// Example filters
	filters := []string{
		`status = "active" && created > @now - 604800000`,
		`title:lower ~ "test%"`,
		`items:length = 5`,
		`age > 42`,
	}

	for _, filter := range filters {
		query, err := TransformFilterToMongoQuery(filter)
		if err != nil {
			fmt.Printf("Error for filter '%s': %v\n", filter, err)
			continue
		}
		queryStr, _ := bson.MarshalExtJSON(query, false, false)
		fmt.Printf("Filter: %s\nMongoDB Query: %s\n\n", filter, string(queryStr))
	}
	}*/
