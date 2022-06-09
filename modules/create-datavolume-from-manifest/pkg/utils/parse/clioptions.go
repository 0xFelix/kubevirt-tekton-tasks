package parse

import (
	"strings"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/env"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/output"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"go.uber.org/zap/zapcore"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	"sigs.k8s.io/yaml"
)

const (
	dvManifestOptionName = "dv-manifest"
)

type CLIOptions struct {
	DataVolumeManifest  string            `arg:"--dv-manifest,env:DV_MANIFEST" placeholder:"MANIFEST" help:"YAML manifest of a DataVolume resource to be created (can be set by DV_MANIFEST env variable)."`
	DataVolumeNamespace string            `arg:"--dv-namespace,env:DV_NAMESPACE" placeholder:"NAMESPACE" help:"Namespace where to create the DV (can be set by DV_NAMESPACE env variable)."`
	WaitForSuccess      string            `arg:"--wait-for-success,env:WAIT_FOR_SUCCESS" help:"Set to \"true\" or \"false\" if container should wait for Ready condition of a DataVolume (can be set by WAIT_FOR_SUCCESS env variable)."`
	Output              output.OutputType `arg:"-o" placeholder:"FORMAT" help:"Output format. One of: yaml|json"`
	Debug               bool              `arg:"--debug" help:"Sets DEBUG log level"`
}

func (c *CLIOptions) GetDebugLevel() zapcore.Level {
	if c.Debug {
		return zapcore.DebugLevel
	}
	return zapcore.InfoLevel
}

func (c *CLIOptions) GetDataVolumeManifest() string {
	return c.DataVolumeManifest
}

func (c *CLIOptions) GetDataVolumeNamespace() string {
	return c.DataVolumeNamespace
}

func (c *CLIOptions) GetWaitForSuccess() bool {
	return c.WaitForSuccess == "true"
}

func (c *CLIOptions) Init() error {
	c.trimSpaces()

	if err := c.assertValidParams(); err != nil {
		return err
	}

	if err := c.assertValidTypes(); err != nil {
		return err
	}

	if err := c.setValues(); err != nil {
		return err
	}

	return nil
}

func (c *CLIOptions) setValues() error {
	if c.GetDataVolumeNamespace() == "" {
		dv := cdiv1beta1.DataVolume{}
		if err := yaml.Unmarshal([]byte(c.DataVolumeManifest), &dv); err != nil {
			return zerrors.NewMissingRequiredError("could not read DV manifest: %v", err.Error())
		}

		if dv.Namespace != "" {
			c.DataVolumeNamespace = dv.Namespace
		} else {
			activeNamespace, err := env.GetActiveNamespace()
			if err != nil {
				return zerrors.NewMissingRequiredError("can't get active namespace: %v", err.Error())
			}
			c.DataVolumeNamespace = activeNamespace
		}
	}

	return nil
}

func (c *CLIOptions) trimSpaces() {
	for _, strVariablePtr := range []*string{&c.DataVolumeManifest, &c.DataVolumeNamespace, &c.WaitForSuccess} {
		*strVariablePtr = strings.TrimSpace(*strVariablePtr)
	}
}

func (c *CLIOptions) assertValidParams() error {
	if c.DataVolumeManifest == "" {
		return zerrors.NewMissingRequiredError("%s param has to be specified", dvManifestOptionName)
	}

	return nil
}

func (c *CLIOptions) assertValidTypes() error {
	if !output.IsOutputType(string(c.Output)) {
		return zerrors.NewMissingRequiredError("%v is not a valid output type", c.Output)
	}
	return nil
}
