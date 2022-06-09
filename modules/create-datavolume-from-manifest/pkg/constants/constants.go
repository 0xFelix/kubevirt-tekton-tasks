package constants

import "time"

// Exit codes
const (
	InvalidCLIInputExitCode    = -1
	DataVolumeCreatorErrorCode = 1
	CreateDataVolumeErrorCode  = 3
	WriteResultsExitCode       = 4
)

// Result names
const (
	NameResultName      = "name"
	NamespaceResultName = "namespace"
)

// Polling
const (
	PollInterval = 1 * time.Second
)
