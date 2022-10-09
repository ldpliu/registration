package v1beta2

import (
	"open-cluster-management.io/api/cluster/v1beta2"
	ctrl "sigs.k8s.io/controller-runtime"
)

type ManagedClusterSet struct {
	v1beta2.ManagedClusterSet
}

func (r *ManagedClusterSet) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}
