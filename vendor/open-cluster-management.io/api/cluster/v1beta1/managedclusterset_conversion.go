package v1beta1

import "k8s.io/klog/v2"

/*
Implementing the hub method is pretty easy -- we just have to add an empty
method called `Hub()` to serve as a
[marker](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/conversion?tab=doc#Hub).
*/

// Hub marks this type as a conversion hub.
func (*ManagedClusterSet) Hub() {
	klog.Errorf("######iN HUB")
}
