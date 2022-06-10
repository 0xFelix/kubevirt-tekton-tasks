# Create DataVolume from Manifest Task

This task creates a DataVolume with oc client.

### Service Account

This task should be run with `create-datavolume-from-manifest-task` serviceAccount.
Please see [RBAC permissions for running the tasks](../../docs/tasks-rbac-permissions.md) for more details.

### Parameters

- **manifest**: YAML manifest of a DataVolume resource to be created.
- **namespace**: Namespace where to create the DV. (defaults to manifest namespace or active namespace)
- **waitForSuccess**: Set to `true` or `false` if container should wait for Ready condition of a DataVolume.
  
### Results

- **name**: The name of DataVolume that was created.
- **namespace**: The namespace of DataVolume that was created.

### Usage

Please see [examples](examples) on how to create DataVolumes.
