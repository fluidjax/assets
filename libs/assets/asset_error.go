package assets

import (
	"fmt"
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
	CodeTypeBadNonce          TransactionCode = 2
	CodeTypeUnauthorized      TransactionCode = 3
	CodeAlreadyExists         TransactionCode = 4
	CodeFailVerfication       TransactionCode = 6
	CodeTypeHTTPError         TransactionCode = 7
	CodeConsensusBalanceError TransactionCode = 8
	CodeConsensusError        TransactionCode = 9
	CodeInsufficientFunds     TransactionCode = 10
	CodeSerializationError    TransactionCode = 11
	CodeIsNil                 TransactionCode = 12

	CodeFailToRebuildAsset       TransactionCode = 100
	CodeDatabaseFail             TransactionCode = 101
	CodeCantUpdateImmutableAsset TransactionCode = 102
	CodeTypeEncodingError        TransactionCode = 103
	CodePayloadEncodingError     TransactionCode = 104

	CodeConsensusErrorFailtoVerifySignature TransactionCode = 200
	CodeConsensusErrorEmptyPayload          TransactionCode = 201
	CodeConsensusMissingFields              TransactionCode = 202
	CodeConsensusIndexNotZero               TransactionCode = 203
	CodeConsensusUnderlyingTXExists         TransactionCode = 204
	CodeConsensusInsufficientFunds          TransactionCode = 205
	CodeConsensusSignedAssetFailtoVerify    TransactionCode = 206
	CodeConsensusWalletNoTransferRules      TransactionCode = 207

	CodeConsensusBalanceFailToAddUnderlying TransactionCode = 205
	CodeTypeTendermintInternalError         TransactionCode = 999

	CodeCLIError TransactionCode = 300

	CodeTendermintInternalError TransactionCode = 999
)

//AssetsError -
type AssetsError struct {
	err  string
	code TransactionCode
}

func (a *AssetsError) Error() string {
	return fmt.Sprintf("%s : %d", a.err, a.code)
}

func NewAssetsError(code TransactionCode, message string) *AssetsError {
	return &AssetsError{
		err:  message,
		code: code,
	}
}

func (a *AssetsError) Code() TransactionCode {
	return a.code
}

// //Wrap - wrap existing AssetsError
// func (ae *AssetsError) Wrap(assetsError *AssetsError, errorString string) {
// 	if ae.Err == nil {
// 		err := errors.New(errorString)
// 		ae.Err = err
// 	}
// 	err := errors.Wrap(assetsError.Err, errorString)
// 	ae.Err = err
// }

// func NewAssetsErrorWithError(code TransactionCode, newDescription string, existingError error) *AssetsError {
// 	return &AssetsError{
// 		Code: code,
// 		Err:  errors.Wrap(existingError, newDescription),
// 	}
// }

// func NewAssetsError(code TransactionCode, description string) Error {
// 	return &AssetsError{
// 		Code: code,
// 		Err:  errors.New(description),
// 	}
// }
