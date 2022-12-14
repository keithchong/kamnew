package webhook

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/openshift/odo/pkg/log"
	"github.com/spf13/cobra"

	"github.com/redhat-developer/kam/pkg/cmd/genericclioptions"
	backend "github.com/redhat-developer/kam/pkg/pipelines/webhook"
	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

const deleteRecommendedCommandName = "delete"

var (
	deleteExample = ktemplates.Examples(`	
	# Delete a Git repository webhook.
	# Example: kam webhook delete --git-host-access-token <git host access token> --env-name dev --service-name taxi 
	
	%[1]s`)
)

type deleteOptions struct {
	options
}

// Run contains the logic for the kam command
func (o *deleteOptions) Run() error {
	ids, err := backend.Delete(o.accessToken, o.pipelinesFolderPath, o.getAppServiceNames(), o.isCICD)

	if len(ids) > 0 {
		if log.IsJSON() {
			outputSuccess(ids)
		} else {
			w := tabwriter.NewWriter(os.Stdout, 5, 2, 3, ' ', tabwriter.TabIndent)
			fmt.Fprintln(w, "DELETED ID")
			fmt.Fprintln(w, "==========")
			for _, id := range ids {
				fmt.Fprintln(w, id)
			}
			w.Flush()
		}
	}

	return err
}

func newCmdDelete(name, fullName string) *cobra.Command {

	o := &deleteOptions{}
	command := &cobra.Command{
		Use:     name,
		Short:   "Delete webhooks.",
		Long:    "Delete all Git repository webhooks that trigger event to CI/CD Pipeline Event Listeners.",
		Example: fmt.Sprintf(deleteExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}

	o.setFlags(command)
	return command
}
