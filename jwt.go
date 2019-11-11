package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"hash"
	"strings"
)

const (
	// HS256 encryption
	HS256  = "HS256"
	secret = "test"
)

// Header rgd
type Header struct {
	alg string
}

func (h *Header) build() string {
	val := `{"typ": "JWT"` +
		`"alg": "` + h.alg + `"}`

	return base64.StdEncoding.EncodeToString([]byte(val))
}

// NewHeader gfgtg
func NewHeader(alg string) *Header {
	header := &Header{alg: alg}
	return header
}

// Payload rgsrg
type Payload struct {
	iss string
	exp string
	sub string
	aud string
}

func (p *Payload) build() string {
	val := ""
	return base64.StdEncoding.EncodeToString([]byte(val))
}

// NewPayload rgesr
func NewPayload(iss string, exp string, sub string, aud string) string {
	payload := &Payload{iss: iss, exp: exp, sub: sub, aud: aud}
	return payload.build()
}

// JWT generates jwt
type JWT struct {
	header    Header
	payload   Payload
	signature string
	secret    string
}

// NewJWT creates new JWT
func NewJWT(header Header, payload Payload, signature string) *JWT {
	return &JWT{header: header, payload: payload, signature: signature, secret: secret}
}

func GetJwt(jwtVal string) *JWT {
	jwt := strings.Split(jwtVal, ".")
	header := Header{}
	json.Unmarshal([]bytes(jwt[0]), &header)

	payload := Payload{}
	json.Unmarshal([]bytes(jwt[1]), &payload)

	signature := ""
	json.Unmarshal([]bytes(jwt[2]), &signature)
	return NewJWT(header, payload, signature)
}

func (j *JWT) generateSignature() string {
	result := header + "." + payload

	var h hash.Hash
	if j.header.alg == HS256 {
		h = hmac.New(sha256.New, []byte(secret))
	}
	h.Write([]byte(result))

	signature := hex.EncodeToString(h.Sum(nil))

	return signature
}

// Build the payload
func (j *JWT) Build() string {
	header := j.header.build()
	payload := j.payload.build()
	signature := j.generateSignature()

	return header + "." + payload + "." + signature
}

// Verify if JWT is valid
func (j *JWT) Verify() bool {
	signature := j.generateSignature()

	return signature == j.signature
}
