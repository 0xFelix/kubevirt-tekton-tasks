package datavolumes

import (
	v1 "k8s.io/api/core/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

func hasDataVolumeFailedToImport(dv *cdiv1beta1.DataVolume) bool {
	conditions := getConditionMap(dv)
	return dv.Status.Phase == cdiv1beta1.ImportInProgress &&
		conditions[cdiv1beta1.DataVolumeBound].Status == v1.ConditionTrue &&
		conditions[cdiv1beta1.DataVolumeRunning].Status == v1.ConditionFalse &&
		conditions[cdiv1beta1.DataVolumeRunning].Reason == "Error"
}

func isDataVolumeImportStatusSuccessful(dv *cdiv1beta1.DataVolume) bool {
	conditions := getConditionMap(dv)
	return dv.Status.Phase == cdiv1beta1.Succeeded &&
		conditions[cdiv1beta1.DataVolumeBound].Status == v1.ConditionTrue
}

func getConditionMap(dv *cdiv1beta1.DataVolume) map[cdiv1beta1.DataVolumeConditionType]cdiv1beta1.DataVolumeCondition {
	result := map[cdiv1beta1.DataVolumeConditionType]cdiv1beta1.DataVolumeCondition{}
	for _, cond := range dv.Status.Conditions {
		result[cond.Type] = cond
	}
	return result
}
