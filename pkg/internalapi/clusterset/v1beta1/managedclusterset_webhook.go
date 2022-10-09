package v1beta1

import (
	"open-cluster-management.io/api/cluster/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
)

type ManagedClusterSet struct {
	v1beta1.ManagedClusterSet
}

func (r *ManagedClusterSet) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}
