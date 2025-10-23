package gocto

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/require"
)

func TestWorkflowSchemaValidation(t *testing.T) {
	wf := Workflow{
		On: WorkflowOn{
			Push: &OnPush{
				OnBranches: &OnBranches{
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
	}
	wfJson, err := json.Marshal(wf)
	require.NoError(t, err)

	c := jsonschema.NewCompiler()
	sch, err := c.Compile("./github-workflow.json")
	require.NoError(t, err)

	buf := bytes.NewBuffer(wfJson)
	inst, err := jsonschema.UnmarshalJSON(buf)
	require.NoError(t, err)

	require.NoError(t, sch.Validate(inst))
}
