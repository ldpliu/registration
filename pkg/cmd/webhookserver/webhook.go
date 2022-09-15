package hub

import (
	"github.com/spf13/cobra"

	"open-cluster-management.io/registration/pkg/webhookserver"
)

func NewWebhook() *cobra.Command {
	ops := webhookserver.NewOptions()

	cmd := &cobra.Command{
		Use:   "webhook-server",
		Short: "Start the webhook server",
		RunE: func(c *cobra.Command, args []string) error {
			err := ops.RunWebhookServer()
			return err
		},
	}

	flags := cmd.Flags()
	ops.AddFlags(flags)
	return cmd
}
