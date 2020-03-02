package assets

import (
	"github.com/pkg/errors"
)

type TransactionCode uint32

/*
	After adding new error codes run stringer in the directory of this source file
	go build golang.org/x/tools/cmd/stringer
	stringer -type=TransactionCode
	This builds the source file transactioncode_string.go, which adds .String() to resolve numeric
	codes to their descriptions
*/
const (
	CodeTypeOK                TransactionCode = 0
	CodeTypeEncodingError     TransactionCode = 1
	CodeTypeBadNonce          TransactionCode = 2
	CodeTypeUnauthorized      TransactionCode = 3
	CodeAlreadyExists         TransactionCode = 4
	CodeDatabaseFail          TransactionCode = 5
	CodeFailVerfication       TransactionCode = 6
	CodeTypeHTTPError         TransactionCode = 7
	CodeConsensusBalanceError TransactionCode = 8
	CodeConsensusError        TransactionCode = 9
	CodeInsufficientFunds     TransactionCode = 10
	CodeSerializationError    TransactionCode = 11
	CodeIsNil                 TransactionCode = 12
	CodeMarshallError         TransactionCode = 13
	CodeUnMarshallError       TransactionCode = 14

	CodeTendermintInternalError TransactionCode = 999
)

//AssetsError -
type AssetsError struct {
	Err  error
	Code TransactionCode
}

//Wrap - wrap existing AssetsError
func (ae *AssetsError) Wrap(assetsError *AssetsError, errorString string) {
	if ae.Err == nil {
		err := errors.New(errorString)
		ae.Err = err
	}
	err := errors.Wrap(assetsError.Err, errorString)
	ae.Err = err
}

func NewAssetsErrorWithError(code TransactionCode, newDescription string, existingError error) *AssetsError {
	return &AssetsError{
		Code: code,
		Err:  errors.Wrap(existingError, newDescription),
	}
}

func NewAssetsError(code TransactionCode, description string) *AssetsError {
	return &AssetsError{
		Code: code,
		Err:  errors.New(description),
	}
}
