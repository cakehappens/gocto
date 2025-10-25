package gocto

import (
	"encoding/json"
	"errors"
	"maps"
	"path"
	"slices"
	"strings"

	"github.com/cakehappens/gocto/internal/util"
)

const (
	DefaultPathToWorkflows = "./.github/workflows"
)

// Workflow
// https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#about-yaml-syntax-for-workflows
type Workflow struct {
	Name        string            `json:"name"`
	RunName     string            `json:"run-name,omitempty,omitzero"`
	On          WorkflowOn        `json:"on"`
	Concurrency Concurrency       `json:"concurrency,omitempty,omitzero"`
	Defaults    Defaults          `json:"defaults,omitempty,omitzero"`
	Env         map[string]string `json:"env,omitempty,omitzero"`
	Permissions Permissions       `json:"permissions,omitempty,omitzero"`
	Jobs        map[string]Job    `json:"jobs"`
	// Storing filename here is useful when you need to reference reusable workflows
	filename string
}

func (w *Workflow) SetFilename(value string) {
	if w == nil {
		return
	}

	w.filename = value
}

func (w *Workflow) GetFilename() string {
	if w == nil {
		return ""
	}

	if w.filename == "" {
		w.filename = FilenameFor(*w)
	}

	return w.filename
}

func (w *Workflow) GetRelativePathAndFilename() string {
	return path.Join(DefaultPathToWorkflows, w.GetFilename())
}

func FilenameFor(w Workflow) string {
	newName := strings.Map(func(r rune) rune {
		switch {
		case '0' <= r && r <= '9':
			fallthrough
		case 'A' <= r && r <= 'Z':
			fallthrough
		case 'a' <= r && r <= 'z':
			return r
		default:
			return '-'
		}
	}, w.Name)

	newName = strings.Trim(newName, "-")
	newName = util.RemoveDupOf(newName, '-')
	newName = strings.ToLower(newName + ".yml")

	return newName
}

// WorkflowOn
// https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#on
type WorkflowOn struct {
	Call              *OnCall        `json:"workflow_call,omitempty,omitzero"`
	Run               *OnWorkflowRun `json:"workflow_run,omitempty,omitzero"`
	Dispatch          *OnDispatch    `json:"workflow_dispatch,omitempty,omitzero"`
	Schedule          *OnSchedule    `json:"schedule,omitempty,omitzero"`
	PullRequest       *OnPullRequest `json:"pull_request,omitempty"`
	PullRequestTarget *OnPullRequest `json:"pull_request_target,omitempty,omitzero"`
	Push              *OnPush        `json:"push,omitempty,omitzero"`
}

type OnCall struct {
	Inputs  map[string]CallInput  `json:"inputs,omitempty,omitzero"`
	Outputs map[string]CallOutput `json:"outputs,omitempty,omitzero"`
}

type CallInput struct {
	Description string        `json:"description,omitempty,omitzero"`
	Default     string        `json:"default,omitempty,omitzero"`
	Required    bool          `json:"required,omitempty,omitzero"`
	Type        CallInputType `json:"type,omitempty,omitzero"`
}

type CallOutput struct {
	Description string `json:"description,omitempty,omitzero"`
	Value       string `json:"value,omitempty,omitzero"`
}

type CallSecrets struct {
	Description string `json:"description,omitempty,omitzero"`
	Required    bool   `json:"required,omitempty,omitzero"`
}

type CallInputType string

const (
	CallInputTypeString  CallInputType = "string"
	CallInputTypeBoolean CallInputType = "boolean"
	CallInputTypeNumber  CallInputType = "number"
)

type OnWorkflowRun struct {
	Workflows   []string `json:"workflows,omitempty,omitzero"`
	Types       []string `json:"types,omitempty,omitzero"`
	*OnBranches `json:",inline"`
}

type OnDispatch struct {
	Inputs map[string]OnDispatchInput `json:"inputs,omitempty,omitzero"`
}

type OnDispatchInput struct {
	Description string              `json:"description,omitempty,omitzero"`
	Required    bool                `json:"required"`
	Default     string              `json:"default,omitempty,omitzero"`
	Type        OnDispatchInputType `json:"type,omitempty,omitzero"`
	Options     []string            `json:"options,omitempty,omitzero"`
}

type OnDispatchInputType string

const (
	OnDispatchInputTypeString      OnDispatchInputType = "string"
	OnDispatchInputTypeBoolean     OnDispatchInputType = "boolean"
	OnDispatchInputTypeNumber      OnDispatchInputType = "number"
	OnDispatchInputTypeEnvironment OnDispatchInputType = "environment"
	OnDispatchInputTypeChoice      OnDispatchInputType = "choice"
)

type OnPaths struct {
	Paths       []string `json:"paths,omitempty,omitzero"`
	PathsIgnore []string `json:"paths-ignore,omitempty,omitzero"`
}

type OnBranches struct {
	Branches       []string `json:"branches,omitempty,omitzero"`
	BranchesIgnore []string `json:"branches-ignore,omitempty,omitzero"`
}

type OnTags struct {
	Tags       []string `json:"tags,omitempty,omitzero"`
	TagsIgnore []string `json:"tags-ignore,omitempty,omitzero"`
}

type OnSchedule struct{}

type OnPullRequest struct {
	*OnPaths    `json:",inline"`
	*OnBranches `json:",inline"`
}

type OnPush struct {
	*OnPaths    `json:",inline"`
	*OnBranches `json:",inline"`
	*OnTags     `json:",inline"`
}

// Concurrency
// https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#concurrency
type Concurrency struct {
	Group            string `json:"group,omitempty,omitzero"`
	CancelInProgress bool   `json:"cancel-in-progress,omitempty,omitzero"`
}

type Job struct {
	Name            string            `json:"name,omitempty,omitzero"`
	Permissions     Permissions       `json:"permissions,omitempty,omitzero"`
	Needs           []string          `json:"needs,omitempty,omitzero"`
	If              string            `json:"if,omitempty,omitzero"`
	RunsOn          []string          `json:"runs-on,omitempty,omitzero"`
	Environment     Environment       `json:"environment,omitempty,omitzero"`
	Concurrency     Concurrency       `json:"concurrency,omitempty,omitzero"`
	Outputs         map[string]string `json:"outputs,omitempty,omitzero"`
	Env             map[string]string `json:"env,omitempty,omitzero"`
	Defaults        Defaults          `json:"defaults,omitempty,omitzero"`
	Steps           []Step            `json:"steps,omitempty,omitzero"`
	TimeoutMinutes  int               `json:"timeout-minutes,omitempty,omitzero"`
	ContinueOnError bool              `json:"continue-on-error,omitempty,omitzero"`
	Uses            string            `json:"uses,omitempty,omitzero"`
	With            map[string]any    `json:"with,omitempty,omitzero"`
	Secrets         *Secrets          `json:"secrets,omitempty,omitzero"`
	Container       Container         `json:"container,omitempty,omitzero"`
}

type Secrets struct {
	Inherit bool              `json:"inherit,omitempty,omitzero"`
	Map     map[string]string `json:",inline,omitempty,omitzero"`
}

func (s *Secrets) MarshalJSON() ([]byte, error) {
	if s == nil {
		return nil, nil
	}

	if s.Inherit {
		type inheritOnlyAlias struct {
			Inherit bool `json:"inherit,omitempty,omitzero"`
		}
		alias := inheritOnlyAlias{
			Inherit: s.Inherit,
		}
		return json.Marshal(&alias)
	}

	return json.Marshal(s.Map)
}

func (s *Secrets) UnmarshalJSON(data []byte) error {
	if util.IsJSONNull(data) {
		return nil
	}

	type inheritOnlyAlias struct {
		Inherit bool `json:"inherit,omitempty,omitzero"`
	}
	var alias inheritOnlyAlias
	if err := json.Unmarshal(data, &alias); err == nil {
		*s = Secrets{
			Inherit: alias.Inherit,
		}
		return nil
	}

	mapVal := make(map[string]string)
	if err := json.Unmarshal(data, &mapVal); err == nil {
		*s = Secrets{
			Map: mapVal,
		}
	}

	return errors.New("unable to unmarshal secrets field, expected bool or map[string]string")
}

type Permissions struct {
	Actions        AccessLevel `json:"actions,omitempty,omitzero"`
	Attestations   AccessLevel `json:"attestations,omitempty,omitzero"`
	Checks         AccessLevel `json:"checks,omitempty,omitzero"`
	Contents       AccessLevel `json:"contents,omitempty,omitzero"`
	Deployments    AccessLevel `json:"deployments,omitempty,omitzero"`
	Discussions    AccessLevel `json:"discussions,omitempty,omitzero"`
	IDToken        AccessLevel `json:"id-token,omitempty,omitzero"`
	Issues         AccessLevel `json:"issues,omitempty,omitzero"`
	Models         AccessLevel `json:"models,omitempty,omitzero"`
	Packages       AccessLevel `json:"packages,omitempty,omitzero"`
	Pages          AccessLevel `json:"pages,omitempty,omitzero"`
	PullRequests   AccessLevel `json:"pull-requests,omitempty,omitzero"`
	SecurityEvents AccessLevel `json:"security-events,omitempty,omitzero"`
	Statuses       AccessLevel `json:"statuses,omitempty,omitzero"`
}

type AccessLevel string

const (
	AccessLevelWrite AccessLevel = "write"
	AccessLevelRead  AccessLevel = "read"
	AccessLevelNone  AccessLevel = "none"
)

type Environment struct {
	Name string `json:"name,omitempty,omitzero"`
	URL  string `json:"url,omitempty,omitzero"`
}

type Defaults struct {
	Run DefaultsRun `json:"run,omitempty,omitzero"`
}

type DefaultsRun struct {
	Shell            Shell  `json:"shell,omitempty,omitzero"`
	WorkingDirectory string `json:"working-directory,omitempty,omitzero"`
}

type Step struct {
	ID               string            `json:"id,omitempty,omitzero"`
	Name             string            `json:"name,omitempty,omitzero"`
	If               string            `json:"if,omitempty,omitzero"`
	Uses             string            `json:"uses,omitempty,omitzero"`
	Run              string            `json:"run,omitempty,omitzero"`
	WorkingDirectory string            `json:"working-directory,omitempty,omitzero"`
	Shell            Shell             `json:"shell,omitempty,omitzero"`
	With             map[string]any    `json:"with,omitempty,omitzero"`
	Env              map[string]string `json:"env,omitempty,omitzero"`
	ContinueOnError  bool              `json:"continue-on-error,omitempty,omitzero"`
	TimeoutMinutes   int               `json:"timeout-minutes,omitempty,omitzero"`
	Strategy         Strategy          `json:"strategy,omitempty,omitzero"`
}

func (s Step) WithName(name string) Step {
	s.Name = name
	return s
}

func (s Step) WithID(id string) Step {
	s.ID = id
	return s
}

func (s Step) WithEnv(key, value string) Step {
	if s.Env == nil {
		s.Env = make(map[string]string)
	}

	s.Env[key] = value
	return s
}

type Shell string

const (
	ShellBash Shell = "bash"
)

type Container struct {
	Image       string               `json:"image,omitempty,omitzero"`
	Env         map[string]string    `json:"env,omitempty,omitzero"`
	Ports       map[string]int       `json:"ports,omitempty,omitzero"`
	Volumes     []string             `json:"volumes,omitempty,omitzero"`
	Credentials ContainerCredentials `json:"credentials,omitempty,omitzero"`
	Options     string               `json:"options,omitempty,omitzero"`
}

type ContainerCredentials struct {
	Username string `json:"username,omitempty,omitzero"`
	Password string `json:"password,omitempty,omitzero"`
}

type Strategy struct {
	Matrix      *Matrix `json:"matrix,omitempty,omitzero"`
	FailFast    bool    `json:"fail-fast,omitempty,omitzero"`
	MaxParallel int     `json:"max-parallel,omitempty,omitzero"`
}

type Matrix struct {
	Map     map[string][]StringOrInt
	Include []map[string]StringOrInt
	Exclude []map[string]StringOrInt
}

func (m *Matrix) UnmarshalJSON(data []byte) error {
	if util.IsJSONNull(data) {
		return nil
	}

	rawMap := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	internalMap := make(map[string][]StringOrInt)
	var internalInclude []map[string]StringOrInt
	var internalExclude []map[string]StringOrInt

	keys := slices.Sorted(maps.Keys(rawMap))
	for _, k := range keys {
		switch k {
		case "include":
			err := json.Unmarshal(rawMap[k], &internalInclude)
			if err != nil {
				return err
			}
		case "exclude":
			err := json.Unmarshal(rawMap[k], &internalExclude)
			if err != nil {
				return err
			}
		default:
			var listStringOrInt []StringOrInt
			err := json.Unmarshal(rawMap[k], &listStringOrInt)
			if err != nil {
				return err
			}

			internalMap[k] = listStringOrInt
		}
	}

	*m = Matrix{
		Map:     internalMap,
		Include: internalInclude,
		Exclude: internalExclude,
	}

	return nil
}

func (m *Matrix) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte(util.JSONNull), nil
	}

	newMap := make(map[string]any)
	for k, v := range m.Map {
		newMap[k] = v
	}

	if len(m.Include) > 0 {
		newMap["include"] = m.Include
	}

	if len(m.Exclude) > 0 {
		newMap["exclude"] = m.Exclude
	}

	return json.Marshal(newMap)
}

type StringOrInt struct {
	StringValue *string
	IntValue    *int
}

func (x *StringOrInt) GetStringValue() string {
	if x.StringValue != nil {
		return *x.StringValue
	}

	return ""
}

func (x *StringOrInt) GetIntValue() int {
	if x.IntValue != nil {
		return *x.IntValue
	}

	return 0
}

func (x *StringOrInt) TryGetStringValue() (string, error) {
	if x.StringValue != nil {
		return *x.StringValue, nil
	} else {
		return "", errors.New("not a string value")
	}
}

func (x *StringOrInt) TryGetIntValue() (int, error) {
	if x.IntValue != nil {
		return *x.IntValue, nil
	} else {
		return 0, errors.New("not an int value")
	}
}

func (x *StringOrInt) UnmarshalJSON(data []byte) error {
	if util.IsJSONNull(data) {
		return nil
	}

	var strVal string
	if err := json.Unmarshal(data, &strVal); err == nil {
		*x = StringOrInt{StringValue: &strVal}
		return nil
	}

	var intVal int
	if err := json.Unmarshal(data, &intVal); err == nil {
		*x = StringOrInt{IntValue: &intVal}
	}

	return errors.New("invalid value in string or int")
}

func (x *StringOrInt) MarshalJSON() ([]byte, error) {
	if x == nil {
		return []byte(util.JSONNull), nil
	}

	if x.IntValue != nil {
		return json.Marshal(x.IntValue)
	}

	if x.StringValue != nil {
		return json.Marshal(x.StringValue)
	}

	return []byte(util.JSONNull), nil
}
