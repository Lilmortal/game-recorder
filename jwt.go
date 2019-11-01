package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"hash"
)

const (
	// HS256 encryption
	HS256 = "HS256"
)

// Header rgd
type Header struct {
	alg string
}

// Build rdgr
func (h *Header) Build() string {
	val := `{"typ": "JWT"` +
		`"alg": "` + h.alg + `"}`

	return base64.StdEncoding.EncodeToString([]byte(val))
}

// Payload rgsrg
type Payload struct {
	iss string
	exp string
	sub string
	aud string
}

// New rf
func (p *Payload) New() {

}

// Build rgr
func (p *Payload) Build() string {
	val := ""
	return base64.StdEncoding.EncodeToString([]byte(val))
}

// Jwt generates jwt
type Jwt struct {
	header  Header
	payload Payload
}

// Build the payload
func (j *Jwt) Build() string {
	secret := ""

	header := j.header.Build()
	payload := j.payload.Build()
	result := header + "." + payload

	var h hash.Hash
	if j.header.alg == HS256 {
		h = hmac.New(sha256.New, []byte(secret))
	}
	h.Write([]byte(result))

	signature := hex.EncodeToString(h.Sum(nil))

	return result + "." + signature
}
