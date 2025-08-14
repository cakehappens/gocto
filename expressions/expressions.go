package expressions

import (
	"fmt"
	"strconv"
)

func StringLiteral(val string) Expression {
	return Expression(`"` + val + `"`)
}

func BoolLiteral(val bool) Expression {
	return Expression(strconv.FormatBool(val))
}

func IntLiteral(val int64) Expression {
	return Expression(strconv.FormatInt(val, 10))
}

func FloatLiteral(val float64) Expression {
	return Expression(strconv.FormatFloat(val, 'f', -1, 64))
}

func From(val string) Expression {
	return Expression(val)
}

type Expression string

func (e Expression) String() string {
	return "${{" + string(e) + "}}"
}

func (e Expression) Or(val Expression) Expression {
	return Expression(string(e) + " || " + string(val)).WithParentheses()
}

func (e Expression) And(val Expression) Expression {
	return Expression(string(e) + " && " + string(val)).WithParentheses()
}

// https://docs.github.com/en/actions/reference/evaluate-expressions-in-workflows-and-actions#literals
// Note that in conditionals, falsy values (false, 0, -0, "", ”, null) are coerced to false
// and truthy (true and other non-falsy values) are coerced to true.
var alwaysFalseVals = []string{
	"",
	`""`,
	`''`,
	// note, this represents boolean false, not a string literal
	// a string literal would be surrounded in quotes
	"false",
	"0",
	"-0",
	"null",
}

func IsAlwaysFalse(val string) bool {
	for _, v := range alwaysFalseVals {
		if val == v {
			return true
		}
	}

	return false
}

func CheckAlwaysFalse(val string) error {
	if IsAlwaysFalse(val) {
		return fmt.Errorf("value is always false: %q", val)
	}

	return nil
}

func IsAlwaysTrue(val string) bool {
	// note, this represents boolean true, not a string literal
	// a string literal would be surrounded in quotes
	return val == "true"
}

func CheckAlwaysTrue(val string) error {
	if IsAlwaysTrue(val) {
		return fmt.Errorf("value is always true: %q", val)
	}

	return nil
}

func IsAlwaysFunc(val string, fn func(string) bool) func() bool {
	return func() bool {
		return fn(val)
	}
}

func CheckAlwaysFunc(val string, fn func(string) error) func() error {
	return func() error {
		return fn(val)
	}
}

func (e Expression) WithParentheses() Expression {
	return Expression("( " + string(e) + " )")
}

func Inputs(val string) Expression {
	return Expression("inputs." + val)
}

func Secrets(val string) Expression {
	return Expression("secrets." + val)
}

func StepOutput(stepID, outputKey string) Expression {
	return Expression("steps." + stepID + ".outputs." + outputKey)
}
