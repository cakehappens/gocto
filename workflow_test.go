package gocto

import (
	"maps"
	"os"
	"slices"
	"testing"

	"github.com/goforj/godump"
	"github.com/kaptinlin/jsonschema"
	"github.com/stretchr/testify/require"
)

func TestWorkflowSchemaValidation(t *testing.T) {
	compiler := jsonschema.NewCompiler()

	contents, err := os.ReadFile("./github-workflow.json")
	require.NoError(t, err)

	schema, err := compiler.Compile(contents)
	require.NoError(t, err)

	result := schema.ValidateStruct(Workflow{
		On: WorkflowOn{
			Push: OnPush{
				OnBranches: OnBranches{
					Branches: []string{"main"},
				},
			},
		},
		Jobs: map[string]Job{
			"foo": {
				RunsOn: []string{"ubuntu-latest"},
				Steps: []Step{
					{
						Run: `echo "foo"`,
					},
				},
			},
		},
	})
	resultErrKeys := slices.Sorted(maps.Keys(result.Errors))

	//assert.Len(t, result.Errors, 0)

	for _, k := range resultErrKeys {
		godump.Dump(result.Errors[k])
	}
}
