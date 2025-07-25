package termprobe

import (
	"errors"
	"fmt"

	"github.com/tiendc/gofn"
)

type ErrorCode int

const (
	UnknownError ErrorCode = iota
	EmptyResponse
	SendSequenceFailed
	ReadResponseFailed
	ResponseStatusBadRequest
	BadResponseUnexpectedEscapeCodes
	BadResponseMissingPrefixSequence
	BadResponseMissingTerminiationSequence
	BadResponseMalformedBody
)

var communicationError []ErrorCode = []ErrorCode{
	SendSequenceFailed,
	ReadResponseFailed,
}

var badResponseError []ErrorCode = []ErrorCode{
	BadResponseMalformedBody,
	BadResponseMissingPrefixSequence,
	BadResponseMissingTerminiationSequence,
	BadResponseUnexpectedEscapeCodes,
}

type ProbeError struct {
	Code    ErrorCode
	Action  ProbeActions
	Message string
	Err     error
}

func (e *ProbeError) Error() string {
	return e.Message
}

func (e *ProbeError) Unwrap() error {
	return e.Err
}

func (e *ProbeError) IsUnknownError() bool {
	return e.Code == UnknownError
}

func (e *ProbeError) IsCommunicationError() bool {
	return gofn.Contain(communicationError, e.Code)
}

func (e *ProbeError) IsEmptyResponse() bool {
	return e.Code == EmptyResponse
}

func (e *ProbeError) IsBadRequest() bool {
	return e.Code == ResponseStatusBadRequest
}

func (e *ProbeError) IsBadResponse() bool {
	return gofn.Contain(badResponseError, e.Code)
}

func IsProbeErrorOrUnknown(err error, action ...ProbeActions) *ProbeError {
	if err == nil {
		return &ProbeError{
			Action:  None,
			Code:    UnknownError,
			Message: "nil error",
			Err:     nil,
		}
	}
	eAction := None
	if len(action) == 1 {
		eAction = action[0]
	}
	var perr *ProbeError
	if errors.As(err, &perr) {
		return perr
	} else {
		return &ProbeError{
			Action:  eAction,
			Code:    UnknownError,
			Message: "unknown error",
			Err:     err,
		}
	}
}

func EmptyResponseError(action ProbeActions) error {
	return &ProbeError{
		Action:  action,
		Code:    EmptyResponse,
		Message: fmt.Sprintf("%s response sequence is empty", action),
	}
}

func SendFailedError(action ProbeActions, err error) error {
	return &ProbeError{
		Action:  action,
		Code:    SendSequenceFailed,
		Message: fmt.Sprintf("failed to send %s sequence", action),
		Err:     err,
	}
}

func ReadFailedError(action ProbeActions, err error) error {
	return &ProbeError{
		Action:  action,
		Code:    ReadResponseFailed,
		Message: fmt.Sprintf("failed to read %s response sequence", action),
		Err:     err,
	}
}

func UnexpectedEscapeCodesError(action ProbeActions) error {
	return &ProbeError{
		Action:  action,
		Code:    BadResponseUnexpectedEscapeCodes,
		Message: fmt.Sprintf("%s response sequence contains unexpected escape code", action),
	}
}

func MissingPrefixError(action ProbeActions, expectedPrefix string) error {
	return &ProbeError{
		Action:  action,
		Code:    BadResponseMissingPrefixSequence,
		Message: fmt.Sprintf("%s response sequence %s prefix missing", action, expectedPrefix),
	}
}

func MissingSuffixError(action ProbeActions, expectedSuffix string) error {
	return &ProbeError{
		Action:  action,
		Code:    BadResponseMissingTerminiationSequence,
		Message: fmt.Sprintf("%s response sequence %s termination suffix missing", action, expectedSuffix),
	}
}

func BadRequestError(action ProbeActions) error {
	return &ProbeError{
		Action:  action,
		Code:    ResponseStatusBadRequest,
		Message: fmt.Sprintf("%s illegal request", action),
	}
}

func BadResponseBodyError(action ProbeActions, msg string) error {
	return &ProbeError{
		Action:  action,
		Code:    BadResponseMalformedBody,
		Message: msg,
	}
}
