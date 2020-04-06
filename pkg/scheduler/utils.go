package scheduler

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/kubernetes-local-volume/kubernetes-local-volume/pkg/common/types"
)

func (lvs *LocalVolumeScheduler) getPodLocalVolumeRequestSize(pod *corev1.Pod) uint64 {
	var result uint64

	for _, volume := range pod.Spec.Volumes {
		if volume.PersistentVolumeClaim != nil {
			pvcName := volume.PersistentVolumeClaim.ClaimName

			// get pvc
			pvc, err := lvs.pvcLister.PersistentVolumeClaims(corev1.NamespaceDefault).Get(pvcName)
			if err != nil {
				continue
			}

			// get storageclass
			sc, err := lvs.storageclassLister.Get(*pvc.Spec.StorageClassName)
			if err != nil {
				continue
			}

			if types.DriverName == sc.Provisioner {
				requestSize, ok := pvc.Spec.Resources.Requests.StorageEphemeral().AsInt64()
				if !ok {
					continue
				}
				result = result + uint64(requestSize)
			}
		}
	}
	return result
}

func (lvs *LocalVolumeScheduler) getPodLocalVolumePVCNames(pod *corev1.Pod) map[string]string {
	result := make(map[string]string)

	for _, volume := range pod.Spec.Volumes {
		if volume.PersistentVolumeClaim != nil {
			pvcName := volume.PersistentVolumeClaim.ClaimName

			// get pvc
			pvc, err := lvs.pvcLister.PersistentVolumeClaims(corev1.NamespaceDefault).Get(pvcName)
			if err != nil {
				continue
			}

			// get storageclass
			sc, err := lvs.storageclassLister.Get(*pvc.Spec.StorageClassName)
			if err != nil {
				continue
			}

			if sc.Provisioner == types.DriverName {
				result[pvc.Name] = ""
			}
		}
	}
	return result
}

func (lvs *LocalVolumeScheduler) getNodeFreeSize(nodeName string) uint64 {
	lv, err := lvs.localvolumeLister.LocalVolumes(corev1.NamespaceDefault).Get(nodeName)
	if err != nil {
		return 0
	}

	var preallocateSize uint64
	for pvcName, _ := range lv.Status.PreAllocated {
		pvc, err := lvs.pvcLister.PersistentVolumeClaims(corev1.NamespaceDefault).Get(pvcName)
		if err != nil {
			continue
		}
		requestSize, ok := pvc.Spec.Resources.Requests.StorageEphemeral().AsInt64()
		if !ok {
			continue
		}
		preallocateSize = preallocateSize + uint64(requestSize)
	}
	return lv.Status.FreeSize - preallocateSize
}
