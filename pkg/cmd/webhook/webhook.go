package webhook

import (
	"crypto/tls"
	"net/http"
	"os"
	"sync"
	"time"

	admissionserver "github.com/openshift/generic-admission-server/pkg/cmd/server"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	genericapiserver "k8s.io/apiserver/pkg/server"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	clusterwebhook "open-cluster-management.io/registration/pkg/webhook/cluster"
	clustersetwebhook "open-cluster-management.io/registration/pkg/webhook/clusterset"
	clustersetbindingwebhook "open-cluster-management.io/registration/pkg/webhook/clustersetbinding"
)

func NewAdmissionHook() *cobra.Command {
	ops := NewOptions()
	cmd := &cobra.Command{
		Use:   "webhook",
		Short: "Start Managed Cluster Admission Server",
		RunE: func(c *cobra.Command, args []string) error {
			stopCh := genericapiserver.SetupSignalHandler()

			if err := ops.ServerOptions.Complete(); err != nil {
				return err
			}
			if err := ops.ServerOptions.Validate(args); err != nil {
				return err
			}
			if err := ops.RunAdmissionServer(ops.ServerOptions, stopCh); err != nil {
				return err
			}
			return nil
		},
	}

	flags := cmd.Flags()
	ops.AddFlags(flags)
	return cmd
}

// Config contains the server (the webhook) cert and key.
type Options struct {
	CertFile      string
	KeyFile       string
	QPS           float32
	Burst         int
	ServerOptions *admissionserver.AdmissionServerOptions
}

// NewOptions constructs a new set of default options for webhook.
func NewOptions() *Options {
	return &Options{
		CertFile: "",
		KeyFile:  "",
		QPS:      100.0,
		Burst:    200,
		ServerOptions: admissionserver.NewAdmissionServerOptions(
			os.Stdout,
			os.Stderr,
			&clusterwebhook.ManagedClusterValidatingAdmissionHook{},
			&clusterwebhook.ManagedClusterMutatingAdmissionHook{},
			&clustersetbindingwebhook.ManagedClusterSetBindingValidatingAdmissionHook{}),
	}
}

func (c *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.CertFile, "tls-cert-file", c.CertFile, ""+
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated "+
		"after server cert).")
	fs.StringVar(&c.KeyFile, "tls-private-key-file", c.KeyFile, ""+
		"File containing the default x509 private key matching --tls-cert-file.")
	fs.Float32Var(&c.QPS, "max-qps", c.QPS,
		"Maximum QPS to the hub server from this webhook.")
	fs.IntVar(&c.Burst, "max-burst", c.Burst,
		"Maximum burst for throttle.")

	featureGate := utilfeature.DefaultMutableFeatureGate
	featureGate.AddFlag(fs)

	c.ServerOptions.RecommendedOptions.FeatureGate = featureGate
	c.ServerOptions.RecommendedOptions.AddFlags(fs)
}

// change the default QPS and Butst, so rewrite this func
func (c *Options) RunAdmissionServer(o *admissionserver.AdmissionServerOptions, stopCh <-chan struct{}) error {
	config, err := o.Config()
	if err != nil {
		return err
	}
	config.RestConfig.QPS = c.QPS
	config.RestConfig.Burst = c.Burst
	server, err := config.Complete().New()
	if err != nil {
		return err
	}

	http.HandleFunc("/clustersetconvert", clustersetwebhook.ServeExampleConvert)

	httpserver := &http.Server{
		Addr:      ":443",
		TLSConfig: ConfigTLS(c),
	}
	httpserver.ListenAndServeTLS("", "")
	return server.GenericAPIServer.PrepareRun().Run(stopCh)
}

type certificateCacheEntry struct {
	cert  *tls.Certificate
	err   error
	birth time.Time
}

// isStale returns true when this cache entry is too old to be usable
func (c *certificateCacheEntry) isStale() bool {
	return time.Since(c.birth) > time.Second
}

func newCertificateCacheEntry(certFile, keyFile string) certificateCacheEntry {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	return certificateCacheEntry{cert: &cert, err: err, birth: time.Now()}
}

// cachingCertificateLoader ensures that we don't hammer the filesystem when opening many connections
// the underlying cert files are read at most once every second
func cachingCertificateLoader(certFile, keyFile string) func() (*tls.Certificate, error) {
	current := newCertificateCacheEntry(certFile, keyFile)
	var currentMtx sync.RWMutex

	return func() (*tls.Certificate, error) {
		currentMtx.RLock()
		if current.isStale() {
			currentMtx.RUnlock()

			currentMtx.Lock()
			defer currentMtx.Unlock()

			if current.isStale() {
				current = newCertificateCacheEntry(certFile, keyFile)
			}
		} else {
			defer currentMtx.RUnlock()
		}

		return current.cert, current.err
	}
}

func ConfigTLS(o *Options) *tls.Config {
	dynamicCertLoader := cachingCertificateLoader(o.CertFile, o.KeyFile)
	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
			return dynamicCertLoader()
		},
	}
}
