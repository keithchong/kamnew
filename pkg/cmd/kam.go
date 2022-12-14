package cmd

import (
	"log"

	env "github.com/redhat-developer/kam/pkg/cmd/component/environment"

	bootstrapnew "github.com/redhat-developer/kam/pkg/cmd/component"
	"github.com/redhat-developer/kam/pkg/cmd/component/component"
	"github.com/redhat-developer/kam/pkg/cmd/environment"
	"github.com/redhat-developer/kam/pkg/cmd/service"
	"github.com/redhat-developer/kam/pkg/cmd/utility"
	"github.com/redhat-developer/kam/pkg/cmd/version"
	"github.com/redhat-developer/kam/pkg/cmd/webhook"
	"github.com/spf13/cobra"
)

var (
	kamLong  = "GitOps Application Manager (KAM) is a CLI tool to scaffold your GitOps repository"
	fullName = "kam"
)

// MakeRootCmd creates and returns the root command for the kam commands.
func MakeRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "kam",
		Short:             "kam",
		Long:              kamLong,
		DisableAutoGenTag: true,
	}

	// Add all subcommands to base command
	rootCmd.AddCommand(
		NewCmdBootstrap(BootstrapRecommendedCommandName, utility.GetFullName(fullName, BootstrapRecommendedCommandName)),
		environment.NewCmdEnv(environment.EnvRecommendedCommandName, utility.GetFullName(fullName, environment.EnvRecommendedCommandName)),
		service.NewCmd(service.RecommendedCommandName, utility.GetFullName(fullName, service.RecommendedCommandName)),
		version.NewCmd(version.RecommendedCommandName, utility.GetFullName(fullName, version.RecommendedCommandName)),
		webhook.NewCmdWebhook(webhook.RecommendedCommandName, utility.GetFullName(fullName, webhook.RecommendedCommandName)),
		NewCmdBuild(BuildRecommendedCommandName, utility.GetFullName(fullName, BuildRecommendedCommandName)),
		completionCmd,
		bootstrapnew.NewCmdBootstrapNew(bootstrapnew.BootstrapRecommendedCommandName, utility.GetFullName(fullName, bootstrapnew.BootstrapRecommendedCommandName)),
		component.NewCmdComp(component.CompRecommendedCommandName, utility.GetFullName(fullName, component.CompRecommendedCommandName)),
		env.NewCmdEnv(env.EnvRecommendedCommandName, utility.GetFullName(fullName, env.EnvRecommendedCommandName)),
		bootstrapnew.NewCmdDescribe(bootstrapnew.DescribeRecommendedCommandName, utility.GetFullName(fullName, bootstrapnew.DescribeRecommendedCommandName)),
	)
	return rootCmd
}

// Execute is the main entry point into this component.
func Execute() {
	if err := MakeRootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
