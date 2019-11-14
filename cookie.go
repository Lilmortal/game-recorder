package main

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

const (
	// JWTTokenKey is used as the cookie name for steam authentication/authorization JWT
	JWTTokenKey = "gameRecordersToken"
	// AntiCSRFTokenKey is used as the cookie name for anti-CSRF-token
	AntiCSRFTokenKey = "anti-csrf-token"
)

func createCookie(name string, value string, maxAge int, secure bool, isHTTPOnly bool, isStrict http.SameSite) (*http.Cookie, error) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     "/",
		Secure:   secure,
		HttpOnly: isHTTPOnly,
		SameSite: http.SameSiteStrictMode,
	}

	return &cookie, nil
}

// GenerateJWTTokenCookie generates a cookie with JWT as the value
func GenerateJWTTokenCookie(cookieName string, value string) (*http.Cookie, error) {
	header := Header{Alg: HS256}
	payload := Payload{Sub: value}

	jwt := NewJWT(header, payload)
	return createCookie(cookieName, jwt.Build(), secondsInWeek, isNotSecure, isHTTPOnly, http.SameSiteDefaultMode)
}

// GenerateAntiCSRFTokenCookie generates a cookie with anti-CSRF token as the value
func GenerateAntiCSRFTokenCookie(cookieName string) (*http.Cookie, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token := base64.URLEncoding.EncodeToString(randomBytes)
	return createCookie(cookieName, token, secondsInHour, isNotSecure, isNotHTTPOnly, http.SameSiteDefaultMode)
}
