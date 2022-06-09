package main

import (
	"net/http"

	goarg "github.com/alexflint/go-arg"
	. "github.com/kubevirt/kubevirt-tekton-tasks/modules/create-datavolume-from-manifest/pkg/constants"
	datavolumecreator "github.com/kubevirt/kubevirt-tekton-tasks/modules/create-datavolume-from-manifest/pkg/datavolumes"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-datavolume-from-manifest/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/exit"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/log"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/output"
	res "github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/results"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/shared/pkg/zerrors"
	"go.uber.org/zap"
)

func main() {
	defer exit.HandleExit()

	cliOptions := &parse.CLIOptions{}
	goarg.MustParse(cliOptions)

	logger := log.InitLogger(cliOptions.GetDebugLevel())
	defer logger.Sync()

	err := cliOptions.Init()
	if err != nil {
		exit.ExitOrDieFromError(InvalidCLIInputExitCode, err)
	}

	log.Logger().Debug("parsed arguments", zap.Reflect("cliOptions", cliOptions))

	dataVolumeCreator, err := datavolumecreator.NewDataVolumeCreator(cliOptions)
	if err != nil {
		exit.ExitOrDieFromError(DataVolumeCreatorErrorCode, err)
	}

	newDataVolume, err := dataVolumeCreator.CreateDataVolume()
	if err != nil {
		exit.ExitOrDieFromError(CreateDataVolumeErrorCode, err,
			zerrors.IsStatusError(err, http.StatusNotFound, http.StatusConflict, http.StatusUnprocessableEntity),
		)
	}

	results := map[string]string{
		NameResultName:      newDataVolume.Name,
		NamespaceResultName: newDataVolume.Namespace,
	}

	log.Logger().Debug("recording results", zap.Reflect("results", results))
	err = res.RecordResults(results)
	if err != nil {
		exit.ExitOrDieFromError(WriteResultsExitCode, err)
	}

	output.PrettyPrint(newDataVolume, cliOptions.Output)
}
