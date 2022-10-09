package conversion

import "github.com/spf13/pflag"

// Config contains the server (the webhook) cert and key.
type Options struct {
	Port                 int
	EnableLeaderElection bool
	CertDir              string
}

// NewOptions constructs a new set of default options for webhook.
func NewOptions() *Options {
	return &Options{
		EnableLeaderElection: true,
		Port:                 9443,
	}
}

func (c *Options) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&c.EnableLeaderElection, "enable-leader-election", c.EnableLeaderElection,
		"Enable leader election for webhook server. "+
			"Enabling this will ensure there is only one active webhook server.")
	fs.IntVar(&c.Port, "port", c.Port,
		"Port is the port that the webhook server serves at.")
	fs.StringVar(&c.CertDir, "certdir", c.CertDir,
		"CertDir is the directory that contains the server key and certificate. If not set, webhook server would look up the server key and certificate in {TempDir}/k8s-webhook-server/serving-certs")
}
