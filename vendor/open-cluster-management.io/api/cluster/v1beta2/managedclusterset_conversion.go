package v1beta2

/*
For imports, we'll need the controller-runtime
[`conversion`](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/conversion?tab=doc)
package, plus the API version for our hub type (v1), and finally some of the
standard packages.
*/
import (
	"k8s.io/klog/v2"
	"open-cluster-management.io/api/cluster/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// +kubebuilder:docs-gen:collapse=Imports

/*
Our "spoke" versions need to implement the
[`Convertible`](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/conversion?tab=doc#Convertible)
interface.  Namely, they'll need `ConvertTo` and `ConvertFrom` methods to convert to/from
the hub version.
*/

/*
ConvertTo is expected to modify its argument to contain the converted object.
Most of the conversion is straightforward copying, except for converting our changed field.
*/
// ConvertTo converts this ManagedClusterSet to the Hub(v1beta1) version.
func (src *ManagedClusterSet) ConvertTo(dstRaw conversion.Hub) error {
	klog.Errorf("######iN CONVERT TO")

	dst := dstRaw.(*v1beta1.ManagedClusterSet)

	dst.ObjectMeta = src.ObjectMeta
	if len(src.Spec.ClusterSelector.SelectorType) == 0 || src.Spec.ClusterSelector.SelectorType == ExclusiveClusterSetLabel {
		dst.Spec.ClusterSelector.SelectorType = v1beta1.SelectorType(v1beta1.LegacyClusterSetLabel)
	} else {
		dst.Spec.ClusterSelector.SelectorType = v1beta1.SelectorType(src.Spec.ClusterSelector.SelectorType)
		dst.Spec.ClusterSelector.LabelSelector = src.Spec.ClusterSelector.LabelSelector
	}
	dst.Status = v1beta1.ManagedClusterSetStatus(src.Status)
	return nil
}

/*
ConvertFrom is expected to modify its receiver to contain the converted object.
Most of the conversion is straightforward copying, except for converting our changed field.
*/

// ConvertFrom converts from the Hub version (v1beta1) to this version.
func (dst *ManagedClusterSet) ConvertFrom(srcRaw conversion.Hub) error {
	klog.Errorf("######iN CONVERT from")

	src := srcRaw.(*v1beta1.ManagedClusterSet)

	dst.ObjectMeta = src.ObjectMeta
	if len(src.Spec.ClusterSelector.SelectorType) == 0 || src.Spec.ClusterSelector.SelectorType == v1beta1.LegacyClusterSetLabel {
		dst.Spec.ClusterSelector.SelectorType = ExclusiveClusterSetLabel
	} else {
		dst.Spec.ClusterSelector.SelectorType = SelectorType(src.Spec.ClusterSelector.SelectorType)
		dst.Spec.ClusterSelector.LabelSelector = src.Spec.ClusterSelector.LabelSelector
	}
	dst.Status = ManagedClusterSetStatus(src.Status)
	return nil
}
