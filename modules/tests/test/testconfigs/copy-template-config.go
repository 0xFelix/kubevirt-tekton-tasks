package testconfigs

import (
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework/testoptions"
	v1 "github.com/openshift/api/template/v1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CopyTemplateTaskData struct {
	Template *v1.Template

	SourceTemplateName      string
	SourceTemplateNamespace string
	TargetTemplateName      string
	TargetTemplateNamespace string
}

type CopyTemplateTestConfig struct {
	TaskRunTestConfig
	TaskData CopyTemplateTaskData

	deploymentNamespace string
}

func (c *CopyTemplateTestConfig) Init(options *testoptions.TestOptions) {
	c.deploymentNamespace = options.DeployNamespace
}

func (c *CopyTemplateTestConfig) GetTaskRun() *v1beta1.TaskRun {
	params := []v1beta1.Param{
		{
			Name: SourceTemplateNameOptionName,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: c.TaskData.SourceTemplateName,
			},
		},
		{
			Name: SourceTemplateNamespaceOptionName,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: c.TaskData.SourceTemplateNamespace,
			},
		},
		{
			Name: TargetTemplateNameOptionName,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: c.TaskData.TargetTemplateName,
			},
		},
		{
			Name: TargetTemplateNamespaceOptionName,
			Value: v1beta1.ArrayOrString{
				Type:      v1beta1.ParamTypeString,
				StringVal: c.TaskData.TargetTemplateNamespace,
			},
		},
	}

	return &v1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestsRandomName(CopyTemplateTaskRunName),
			Namespace: c.deploymentNamespace,
		},
		Spec: v1beta1.TaskRunSpec{
			TaskRef: &v1beta1.TaskRef{
				Name: CopyTemplateClusterTaskName,
				Kind: v1beta1.ClusterTaskKind,
			},
			Timeout:            &metav1.Duration{Duration: c.GetTaskRunTimeout()},
			ServiceAccountName: c.ServiceAccount,
			Params:             params,
		},
	}
}
