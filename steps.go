package gocto

import (
	"fmt"
	"strings"
)

// TernaryExpressionStep is a better way to do if/else than with pure expressions,
// because of the pitfalls of truthy and falsy values
//
// https://github.com/orgs/community/discussions/26738
// Ternary Expressions suck in GitHub Actions, don't use them, instead use a bash expression
// more about all the pitfalls of trying to use || && here:
// https://7tonshark.com/posts/github-actions-ternary-operator/
//
// the output name is "value", e.g.
// steps.my-step-id.outputs.value
func TernaryExpressionStep(stepID, bashCond, thenVal, elseVal string) Step {
	run := strings.TrimSpace(fmt.Sprintf(`
if [[ %s ]]; then
	echo "::set-output name=value::%s"
else
	echo "::set-output name=value::%s"
fi
`, bashCond, thenVal, elseVal))

	return Step{
		ID:    stepID,
		Run:   run,
		Shell: ShellBash,
	}
}
