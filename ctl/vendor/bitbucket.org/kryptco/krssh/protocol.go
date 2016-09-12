package krssh

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"golang.org/x/crypto/ssh"
)

type Request struct {
	RequestID       string       `json:"request_id"`
	SignRequest     *SignRequest `json:"sign_request"`
	ListRequest     *ListRequest `json:"list_request"`
	MeRequest       *MeRequest   `json:"me_request"`
	timeoutOverride *int
}

func NewRequest() (request Request, err error) {
	id, err := Rand128Base62()
	if err != nil {
		return
	}
	request = Request{
		RequestID: id,
	}
	return
}

type Response struct {
	RequestID      string        `json:"request_id"`
	SignResponse   *SignResponse `json:"sign_response"`
	ListResponse   *ListResponse `json:"list_response"`
	MeResponse     *MeResponse   `json:"me_response"`
	SNSEndpointARN *string       `json:"sns_endpoint_arn"`
}

type SignRequest struct {
	//	N.B. []byte marshals to base64 encoding in JSON
	Digest []byte `json:"digest"`
	//	SHA256 hash of public key DER
	PublicKeyFingerprint []byte `json:"public_key_fingerprint"`
}

type SignResponse struct {
	Signature *[]byte `json:"signature"`
	Error     *string `json:"error"`
}

type ListRequest struct {
	EmailFilter *string `json:"email_filter"`
}

type ListResponse struct {
	Profiles []Profile `json:"profiles"`
}

type Profile struct {
	PublicKeyDER []byte `json:"public_key_der"`
	Email        string `json:"email"`
}

func (p Profile) DisplayString() string {
	pkFingerprint := sha256.Sum256(p.PublicKeyDER)
	return base64.StdEncoding.EncodeToString(pkFingerprint[:]) + " <" + p.Email + ">"
}
func (p Profile) SSHWireString() (wireString string, err error) {
	x509Pk, err := x509.ParsePKIXPublicKey(p.PublicKeyDER)
	if err != nil {
		return
	}
	sshPk, err := ssh.NewPublicKey(x509Pk)
	if err != nil {
		return
	}
	wireString = sshPk.Type() + " " + base64.StdEncoding.EncodeToString(sshPk.Marshal())
	return
}

type MeRequest struct{}

type MeResponse struct {
	Me Profile `json:"me"`
}