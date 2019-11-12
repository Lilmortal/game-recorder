package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

const (
	// HS256 encryption used to encrypt JWT tokens
	HS256 = "HS256"
	// TODO: Generate this randomly and get it from env variables
	secret = "test"
)

// Header is our JWT header, which specifies what algorithm to use
type Header struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
}

func (h *Header) build() string {
	val := `{"typ": "` + h.Typ + `", "alg": "` + h.Alg + `"}`

	return base64.StdEncoding.EncodeToString([]byte(val))
}

// NewHeader generates a new header
func NewHeader(alg string) *Header {
	header := &Header{Typ: "JWT", Alg: alg}
	return header
}

// Payload is our JWT payload, which contains the actual values
type Payload struct {
	Iss string `json:"iss"`
	Exp string `json:"exp"`
	Sub string `json:"sub"`
	Aud string `json:"aud"`
}

func (p *Payload) build() string {
	val := `{"sub": "` + p.Sub + `"}`
	return base64.StdEncoding.EncodeToString([]byte(val))
}

// NewPayload generates a new payload
func NewPayload(iss string, exp string, sub string, aud string) string {
	payload := &Payload{Iss: iss, Exp: exp, Sub: sub, Aud: aud}
	return payload.build()
}

// JWT is what we used to verify our users
type JWT struct {
	header    Header
	payload   Payload
	signature string
	secret    string
}

// NewJWT generates a new JWT
func NewJWT(header Header, payload Payload) *JWT {
	signature := generateSignature(header, payload)
	return &JWT{header: header, payload: payload, signature: signature, secret: secret}
}

// Signature appends signature to an existing JWT
func (j *JWT) Signature(signature string) *JWT {
	j.signature = signature
	return j
}

// GetJWT gets a json string and returns back a JWT struct
func GetJWT(jwtVal string) *JWT {
	jwt := strings.Split(jwtVal, ".")

	header := Header{}
	headerVal, _ := base64.StdEncoding.DecodeString(jwt[0])
	json.Unmarshal(headerVal, &header)

	payload := Payload{}
	payloadVal, _ := base64.StdEncoding.DecodeString(jwt[1])
	json.Unmarshal(payloadVal, &payload)

	signature := jwt[2]

	newJwt := NewJWT(header, payload)
	return newJwt.Signature(signature)
}

func generateSignature(h Header, p Payload) string {
	header := h.build()
	payload := p.build()

	result := header + "." + payload

	signature := ""
	if h.Alg == HS256 {
		hash := hmac.New(sha256.New, []byte(secret))
		_, err := hash.Write([]byte(result))
		if err != nil {
			log.Fatalf("main.go: failed to write to hash. %s", err.Error())
		}

		signature = hex.EncodeToString(hash.Sum(nil))
	}

	return signature
}

// Build returns the JWT in string format
func (j *JWT) Build() string {
	header := j.header.build()
	payload := j.payload.build()
	signature := j.signature

	return header + "." + payload + "." + signature
}

// Verify if JWT is valid via signatures
func (j *JWT) Verify() bool {
	signature := generateSignature(j.header, j.payload)

	fmt.Println(signature, j.signature)
	return signature == j.signature
}
