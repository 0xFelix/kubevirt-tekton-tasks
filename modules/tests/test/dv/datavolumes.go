package dv

import (
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1beta12 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1beta1"
)

type TestDataVolumeAttachmentType string

const (
	OwnedPVC TestDataVolumeAttachmentType = "owned-pvc"
	PVC      TestDataVolumeAttachmentType = "pvc"
	DV       TestDataVolumeAttachmentType = "dv"
	OwnedDV  TestDataVolumeAttachmentType = "owned-dv"
)

type TestDataVolume struct {
	Data           *v1beta12.DataVolume
	AttachmentType TestDataVolumeAttachmentType
	DiskName       string
}

func NewBlankDataVolume(name string) *TestDataVolume {
	volumeMode := v1.PersistentVolumeFilesystem
	datavolume := &v1beta12.DataVolume{
		TypeMeta: metav1.TypeMeta{
			APIVersion: constants.DataVolumeApiVersion,
			Kind:       constants.DataVolumeKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1beta12.DataVolumeSpec{
			PVC: &v1.PersistentVolumeClaimSpec{
				AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceStorage: *resource.NewScaledQuantity(100, resource.Mega),
					},
				},
				VolumeMode: &volumeMode,
			},
			Source: v1beta12.DataVolumeSource{
				Blank: &v1beta12.DataVolumeBlankImage{},
			},
		},
	}

	return &TestDataVolume{
		datavolume,
		"",
		"",
	}
}

func (d *TestDataVolume) WithoutTypeMeta() *TestDataVolume {
	d.Data.Kind = ""
	d.Data.APIVersion = ""
	return d
}

func (d *TestDataVolume) AttachAs(attachmentType TestDataVolumeAttachmentType) *TestDataVolume {
	d.AttachmentType = attachmentType
	return d
}

func (d *TestDataVolume) AttachWithDiskName(diskName string) *TestDataVolume {
	d.DiskName = diskName
	return d
}

func (d *TestDataVolume) WithURLSource(url string) *TestDataVolume {
	d.Data.Spec.Source.Blank = nil
	d.Data.Spec.Source.HTTP = &v1beta12.DataVolumeSourceHTTP{
		URL: url,
	}
	return d
}

func (d *TestDataVolume) WithRegistrySource(registryURL string) *TestDataVolume {
	d.Data.Spec.Source.Blank = nil
	d.Data.Spec.Source.Registry = &v1beta12.DataVolumeSourceRegistry{
		URL: registryURL,
	}
	return d
}

func (d *TestDataVolume) WithSize(size int64, scale resource.Scale) *TestDataVolume {
	d.Data.Spec.PVC.Resources.Requests[v1.ResourceStorage] = *resource.NewScaledQuantity(size, scale)
	return d
}

func (d *TestDataVolume) Build() *v1beta12.DataVolume {
	return d.Data
}
