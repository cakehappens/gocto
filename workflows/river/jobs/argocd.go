package jobs

import "github.com/ghostsquad/alveus/internal/integrations/github"

func WithArgoCDSyncAndWaitSteps() Option {
	return func(options *Options) error {
		options.Steps = append(options.Steps, github.Step{})

		return nil
	}
}
