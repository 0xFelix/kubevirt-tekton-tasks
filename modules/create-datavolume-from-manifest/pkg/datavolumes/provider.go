package datavolumes

import (
	"context"
	"time"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-datavolume-from-manifest/pkg/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-datavolume-from-manifest/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	cdiclientv1beta1 "kubevirt.io/containerized-data-importer/pkg/client/clientset/versioned/typed/core/v1beta1"
	"sigs.k8s.io/yaml"
)

type dataVolumeProvider struct {
	client cdiclientv1beta1.CdiV1beta1Interface
}

type DataVolumeProvider interface {
	Get(string, string) (*cdiv1beta1.DataVolume, error)
	Create(*cdiv1beta1.DataVolume) (*cdiv1beta1.DataVolume, error)
	RESTClient() rest.Interface
}

func NewDataVolumeProvider(client cdiclientv1beta1.CdiV1beta1Interface) DataVolumeProvider {
	return &dataVolumeProvider{
		client: client,
	}
}

func (d *dataVolumeProvider) Get(namespace string, name string) (*cdiv1beta1.DataVolume, error) {
	return d.client.DataVolumes(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (d *dataVolumeProvider) Create(dataVolume *cdiv1beta1.DataVolume) (*cdiv1beta1.DataVolume, error) {
	return d.client.DataVolumes(dataVolume.Namespace).Create(context.TODO(), dataVolume, metav1.CreateOptions{})
}

func (d *dataVolumeProvider) RESTClient() rest.Interface {
	return d.client.RESTClient()
}

type DataVolumeCreator struct {
	cliOptions         *parse.CLIOptions
	dataVolumeProvider DataVolumeProvider
}

func NewDataVolumeCreator(cliOptions *parse.CLIOptions) (*DataVolumeCreator, error) {
	log.Logger().Debug("initialized clients and providers")

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	return &DataVolumeCreator{
		cliOptions:         cliOptions,
		dataVolumeProvider: NewDataVolumeProvider(cdiclientv1beta1.NewForConfigOrDie(config)),
	}, nil
}

func (d *DataVolumeCreator) CreateDataVolume() (*cdiv1beta1.DataVolume, error) {
	dv := &cdiv1beta1.DataVolume{}
	if err := yaml.Unmarshal([]byte(d.cliOptions.DataVolumeManifest), dv); err != nil {
		return nil, zerrors.NewSoftError("could not read DV manifest: %v", err.Error())
	}
	if dv.Kind == "" || dv.APIVersion == "" {
		return nil, zerrors.NewSoftError("could not read DV manifest: kind or apiVersion missing")
	}
	dv.Namespace = d.cliOptions.GetDataVolumeNamespace()

	log.Logger().Debug("creating DV", zap.Reflect("dv", dv))
	dv, err := d.dataVolumeProvider.Create(dv)
	if err != nil {
		return nil, zerrors.NewSoftError("could not create DV: %v", err.Error())
	}

	if d.cliOptions.GetWaitForSuccess() {
		log.Logger().Debug("waiting for success of DV", zap.Reflect("dv", dv))
		if err := d.waitForSuccess(dv); err != nil {
			return nil, zerrors.NewSoftError("Failed to wait for success of DV: %v", err.Error())
		}
	}

	return dv, nil
}

func (d *DataVolumeCreator) waitForSuccess(dv *cdiv1beta1.DataVolume) error {
	return wait.PollImmediate(constants.PollInterval, time.Second*600, func() (bool, error) {
		dv, err := d.dataVolumeProvider.Get(dv.Namespace, dv.Name)
		if err != nil {
			return false, err
		}

		if isDataVolumeImportStatusSuccessful(dv) {
			return true, nil
		}

		if hasDataVolumeFailedToImport(dv) {
			return false, zerrors.NewSoftError("Import of DV failed: %v", dv)
		}

		if dv.Status.Phase == cdiv1beta1.Failed {
			return false, zerrors.NewSoftError("DV is in phase failed: %v", dv)
		}

		return false, nil
	})
}
