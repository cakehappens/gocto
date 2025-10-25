package gocto

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkflowSchemaValidation(t *testing.T) {
	type testCase struct {
		wf         Workflow
		assertions []func(t *testing.T, marshalled string)
	}

	cases := []testCase{
		{
			wf: Workflow{
				Name: "minimal",
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
			},
		},
		{
			wf: Workflow{
				Name: "inherit secrets",
				On: WorkflowOn{
					Push: &OnPush{},
				},
				Jobs: map[string]Job{
					"foo": {
						Uses: "./foo.yaml",
						Secrets: &Secrets{
							Inherit: true,
						},
					},
				},
			},
			assertions: []func(t *testing.T, marshalled string){
				func(t *testing.T, marshalled string) {
					assert.Regexp(t, `"secrets":"inherit"`, marshalled)
				},
			},
		},
		{
			wf: Workflow{
				Name: "secrets map",
				On: WorkflowOn{
					Push: &OnPush{},
				},
				Jobs: map[string]Job{
					"foo": {
						Uses: "./foo.yaml",
						Secrets: &Secrets{
							Map: map[string]string{
								"foo": "bar",
							},
						},
					},
				},
			},
			assertions: []func(t *testing.T, marshalled string){
				func(t *testing.T, marshalled string) {
					assert.Regexp(t, `"secrets":{`, marshalled)
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.wf.Name, func(t *testing.T) {
			wfJson, err := json.Marshal(tc.wf)
			require.NoError(t, err)

			c := jsonschema.NewCompiler()
			sch, err := c.Compile("./github-workflow.json")
			require.NoError(t, err)

			buf := bytes.NewBuffer(wfJson)
			inst, err := jsonschema.UnmarshalJSON(buf)
			require.NoError(t, err)

			require.NoError(t, sch.Validate(inst))

			for _, assertion := range tc.assertions {
				assertion(t, string(wfJson))
			}
		})
	}
}
