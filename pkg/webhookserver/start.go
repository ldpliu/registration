package webhookserver

import (
	"k8s.io/klog/v2"

	"k8s.io/apimachinery/pkg/runtime"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	kbatchv1 "k8s.io/api/batch/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	clusterv1beta1 "open-cluster-management.io/api/cluster/v1beta1"
	clusterv1beta2 "open-cluster-management.io/api/cluster/v1beta2"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(kbatchv1.AddToScheme(scheme)) // we've added this ourselves
	utilruntime.Must(clusterv1beta1.AddToScheme(scheme))
	utilruntime.Must(clusterv1beta2.AddToScheme(scheme))
}

func (c *Options) RunWebhookServer() error {

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:           scheme,
		Port:             c.Port,
		LeaderElection:   c.EnableLeaderElection,
		LeaderElectionID: "webhook-server",
		CertDir:          c.CertDir,
	})

	if err != nil {
		klog.Error(err, "unable to start manager")
		return err
	}

	if err = (&clusterv1beta1.ManagedClusterSet{}).SetupWebhookWithManager(mgr); err != nil {
		klog.Error(err, "unable to create webhook", "webhook", "CronJob")
		return err
	}
	if err = (&clusterv1beta2.ManagedClusterSet{}).SetupWebhookWithManager(mgr); err != nil {
		klog.Error(err, "unable to create webhook", "webhook", "CronJob")
		return err
	}

	klog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		klog.Error(err, "problem running manager")
		return err
	}
	return nil
}
