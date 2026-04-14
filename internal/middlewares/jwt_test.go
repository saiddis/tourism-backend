package middleware

import "testing"

func TestGenerateAndValidateAccessToken(t *testing.T) {
	InitSecrets("access-secret", "refresh-secret")

	token, err := GenerateAccessToken(42, "alice@example.com", "admin")
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	claims, err := ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("ValidateAccessToken() error = %v", err)
	}

	if claims.UserID != 42 {
		t.Fatalf("UserID = %d, want 42", claims.UserID)
	}
	if claims.Email != "alice@example.com" {
		t.Fatalf("Email = %q, want alice@example.com", claims.Email)
	}
	if claims.Role != "admin" {
		t.Fatalf("Role = %q, want admin", claims.Role)
	}
	if claims.Type != accessTokenType {
		t.Fatalf("Type = %q, want %q", claims.Type, accessTokenType)
	}
}

func TestGenerateAndValidateRefreshToken(t *testing.T) {
	InitSecrets("access-secret", "refresh-secret")

	token, err := GenerateRefreshToken(99)
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}

	claims, err := ValidateRefreshToken(token)
	if err != nil {
		t.Fatalf("ValidateRefreshToken() error = %v", err)
	}

	if claims.UserID != 99 {
		t.Fatalf("UserID = %d, want 99", claims.UserID)
	}
	if claims.Type != refreshTokenType {
		t.Fatalf("Type = %q, want %q", claims.Type, refreshTokenType)
	}
}

func TestValidateAccessTokenRejectsRefreshToken(t *testing.T) {
	InitSecrets("access-secret", "refresh-secret")

	refreshToken, err := GenerateRefreshToken(99)
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}

	if _, err := ValidateAccessToken(refreshToken); err == nil {
		t.Fatal("ValidateAccessToken() unexpectedly accepted a refresh token")
	}
}

func TestGenerateTokenPair(t *testing.T) {
	InitSecrets("access-secret", "refresh-secret")

	pair, err := GenerateTokenPair(7, "bob@example.com", "client")
	if err != nil {
		t.Fatalf("GenerateTokenPair() error = %v", err)
	}

	if pair.AccessToken == "" {
		t.Fatal("AccessToken is empty")
	}
	if pair.RefreshToken == "" {
		t.Fatal("RefreshToken is empty")
	}
	if pair.TokenType != "Bearer" {
		t.Fatalf("TokenType = %q, want Bearer", pair.TokenType)
	}
	if pair.AccessExpiresIn <= 0 {
		t.Fatalf("AccessExpiresIn = %d, want positive value", pair.AccessExpiresIn)
	}
	if pair.RefreshExpiresIn <= pair.AccessExpiresIn {
		t.Fatalf("RefreshExpiresIn = %d, want greater than AccessExpiresIn = %d", pair.RefreshExpiresIn, pair.AccessExpiresIn)
	}
}
