package jobs

import (
	"errors"
	"fmt"

	"github.com/ghostsquad/alveus/internal/integrations/github"
)

type Options struct {
	Permissions github.Permissions
	Needs       []string
	If          string
	RunsOn      []string
	Outputs     map[string]string
	Environment github.Environment
	Concurrency github.Concurrency
	Env         map[string]string
	Defaults    github.Defaults
	Steps       []github.Step
}

type Option func(*Options) error

type CommitChangesInput struct {
	CommitMessage string
	Ref           string
}

func (input *CommitChangesInput) Validate() error {
	var errs []error

	if input.Ref == "" {
		errs = append(errs, fmt.Errorf("ref is required"))
	}

	if input.CommitMessage == "" {
		errs = append(errs, fmt.Errorf("commit message is required"))
	}

	return errors.Join(errs...)
}

type JobInput struct {
	JobName string
}

func (input *JobInput) Validate() error {
	var errs []error

	if input.JobName == "" {
		errs = append(errs, fmt.Errorf("job name is required"))
	}

	return errors.Join(errs...)
}

func New(input JobInput, options ...Option) (github.Job, error) {
	opts := &Options{
		Outputs: make(map[string]string),
		Env:     make(map[string]string),
	}

	var errs []error

	for _, o := range options {
		errs = append(errs, o(opts))
	}

	if err := input.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("invalid input: %w", err))
	}

	if err := errors.Join(errs...); err != nil {
		return github.Job{}, fmt.Errorf("input problem(s) found: %w", err)
	}

	job := github.Job{
		Name:        input.JobName,
		Permissions: opts.Permissions,
		Needs:       opts.Needs,
		If:          opts.If,
		RunsOn:      opts.RunsOn,
		Environment: opts.Environment,
		Concurrency: opts.Concurrency,
		Outputs:     opts.Outputs,
		Env:         opts.Env,
		Defaults:    opts.Defaults,
		Steps:       opts.Steps,
	}

	return job, nil
}
