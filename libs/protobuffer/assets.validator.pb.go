// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: assets.proto

package protobuffer

import (
	fmt "fmt"
	math "math"
	proto "github.com/golang/protobuf/proto"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *PBSignedAsset) Validate() error {
	// Validation of proto3 map<> fields is unsupported.
	if this.Asset != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Asset); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Asset", err)
		}
	}
	return nil
}
func (this *PBAsset) Validate() error {
	if oneOfNester, ok := this.GetPayload().(*PBAsset_Wallet); ok {
		if oneOfNester.Wallet != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.Wallet); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Wallet", err)
			}
		}
	}
	if oneOfNester, ok := this.GetPayload().(*PBAsset_TrusteeGroup); ok {
		if oneOfNester.TrusteeGroup != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.TrusteeGroup); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("TrusteeGroup", err)
			}
		}
	}
	if oneOfNester, ok := this.GetPayload().(*PBAsset_Iddoc); ok {
		if oneOfNester.Iddoc != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.Iddoc); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Iddoc", err)
			}
		}
	}
	// Validation of proto3 map<> fields is unsupported.
	return nil
}
func (this *PBTransfer) Validate() error {
	// Validation of proto3 map<> fields is unsupported.
	return nil
}
func (this *PBWallet) Validate() error {
	return nil
}
func (this *PBTrusteeGroup) Validate() error {
	if this.TrusteeGroup != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.TrusteeGroup); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("TrusteeGroup", err)
		}
	}
	return nil
}
func (this *PBIDDoc) Validate() error {
	return nil
}
