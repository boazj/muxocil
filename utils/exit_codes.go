package utils

//go:generate stringer -type ExitCode
type ExitCode int

const (
	ExitOk ExitCode = iota
	ExitUnknownOs
	ExitGeneralError
	ExitConfigFailure
	ExitConfigBindError
	ExitOpenEditorError
	ExitProviderDataError
	ExitProviderUnknownLayout
	ExitCmdUseBadCommand
	ExitCmdUseBadFile
	ExitProviderFailFast
	ExitProviderValidation
	ExitProviderFailure
)
