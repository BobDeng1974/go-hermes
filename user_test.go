package main

import (
	"bytes"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	const expectedLength = 64
	user := &User{}
	user.generateRandomToken()
	var generatedToken = user.Token

	if len(user.Token) != expectedLength {
		t.Errorf("Generated token should be %d, got: %d", expectedLength, len(user.Token))
	}

	if user.tokenGenerated != true {
		t.Error("tokenGenerated flag should be true after token is generated")
	}

	// make sure generateRandomToken() is idempotent
	user.generateRandomToken()

	if generatedToken != user.Token {
		t.Error("generateRandomToken() should not generate different outputs if called twice!")
	}
}

func TestEncryptPassword(t *testing.T) {
	var providedPassword = []byte("12345")
	user := &User{Password: providedPassword}
	user.encryptPassword()
	var encryptedPassword = user.Password

	if bytes.Equal(providedPassword, encryptedPassword) {
		t.Error("Password was not encrypted!")
	}

	if !user.passwordEncoded {
		t.Error("passwordEncoded must be true after encryptPassword() is called.")
	}

	// try encrypt password again, to make sure it is idempotent
	user.encryptPassword()
	if !bytes.Equal(encryptedPassword, user.Password) {
		t.Error("password must not be encrypted twice!")
	}
}
