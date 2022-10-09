package v1beta2

import (
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"open-cluster-management.io/api/cluster/v1beta1"
	"open-cluster-management.io/api/cluster/v1beta2"

	internalv1beta1 "open-cluster-management.io/registration/pkg/internalapi/clusterset/v1beta1"
)

func TestConvertTo(t *testing.T) {
	cases := []struct {
		name           string
		oriSet         *ManagedClusterSet
		expectedDstSet *internalv1beta1.ManagedClusterSet
	}{
		{
			name: "test empty spec set",
			oriSet: &ManagedClusterSet{
				v1beta2.ManagedClusterSet{
					ObjectMeta: metav1.ObjectMeta{
						Name: "mcs1",
					},
				},
			},
			expectedDstSet: &internalv1beta1.ManagedClusterSet{
				ManagedClusterSet: v1beta1.ManagedClusterSet{
					ObjectMeta: metav1.ObjectMeta{
						Name: "mcs1",
					},
					Spec: v1beta1.ManagedClusterSetSpec{
						ClusterSelector: v1beta1.ManagedClusterSelector{
							SelectorType: v1beta1.LegacyClusterSetLabel,
						},
					},
				},
			},
		},
		{
			name: "test empty spec set",
			oriSet: &ManagedClusterSet{
				v1beta2.ManagedClusterSet{
					ObjectMeta: metav1.ObjectMeta{
						Name: "mcs1",
					},
					Spec: v1beta1.ManagedClusterSetSpec{}
				},
			},
			expectedDstSet: &internalv1beta1.ManagedClusterSet{
				ManagedClusterSet: v1beta1.ManagedClusterSet{
					ObjectMeta: metav1.ObjectMeta{
						Name: "mcs1",
					},
					Spec: v1beta1.ManagedClusterSetSpec{
						ClusterSelector: v1beta1.ManagedClusterSelector{
							SelectorType: v1beta1.LegacyClusterSetLabel,
						},
					},
				},
			},
		},
		{
			name: "test exclusive set",
			oriSet: &ManagedClusterSet{
				v1beta2.ManagedClusterSet{
					ObjectMeta: metav1.ObjectMeta{
						Name: "mcs1",
					},
					Spec: v1beta1.ManagedClusterSetSpec{
						ClusterSelector: v1beta1.ManagedClusterSelector{
							SelectorType: v1beta1.ExclusiveClusterSetLabel,
						},
					},
				},
			},
			expectedDstSet: &internalv1beta1.ManagedClusterSet{
				ManagedClusterSet: v1beta1.ManagedClusterSet{
					ObjectMeta: metav1.ObjectMeta{
						Name: "mcs1",
					},
					Spec: v1beta1.ManagedClusterSetSpec{
						ClusterSelector: v1beta1.ManagedClusterSelector{
							SelectorType: v1beta1.LegacyClusterSetLabel,
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dstSet := &internalv1beta1.ManagedClusterSet{}
			c.oriSet.ConvertTo(dstSet)
			if !reflect.DeepEqual(dstSet, c.expectedDstSet) {
				t.Errorf("Faild to convert clusterset. expectDstSet:%v , dstSet:%v", c.expectedDstSet, dstSet)
			}
		})
	}
}
