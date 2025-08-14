package jobs

import (
	"fmt"

	"github.com/ghostsquad/alveus/internal/integrations/github"
	"github.com/ghostsquad/alveus/internal/integrations/github/expressions"
	"github.com/ghostsquad/alveus/internal/util"
)

func WithCommitChangesSteps(input CommitChangesInput, intermediateSteps ...github.Step) Option {
	return func(options *Options) error {
		if err := input.Validate(); err != nil {
			return fmt.Errorf("invalid commit changes input: %w", err)
		}

		steps := []github.Step{
			{
				Uses: "actions/checkout@v4",
				With: map[string]any{
					"fetch-depth": 1,
					"ssh-key":     expressions.Secrets("WRITE_DEPLOY_KEY").String(),
					"ref":         input.Ref,
				},
			},
			{
				Name: "configure-git",
				Run: util.Join("\n",
					"git config --local user.email 'github-actions@github.com'",
					"git config --local user.name 'GitHub Actions'",
				),
			},
		}

		steps = append(steps, intermediateSteps...)

		steps = append(steps, github.Step{
			Run: util.Join("\n",
				"git add .",
				fmt.Sprintf("git commit -m %q", input.CommitMessage),
				"git push",
			),
		})
		options.Steps = append(options.Steps, steps...)

		return nil
	}
}
