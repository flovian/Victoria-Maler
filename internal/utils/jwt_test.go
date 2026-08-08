package utils

import "testing"

const testSecret = "unit-test-secret"

func TestTokenRoundTrip(t *testing.T) {
	token, err := GenerateToken(testSecret, 42, "ngo@test.org", "ngo")
	if err != nil {
		t.Fatal(err)
	}

	claims, err := ParseToken(testSecret, token)
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != 42 {
		t.Fatalf("expected user 42, got %d", claims.UserID)
	}
	if claims.Email != "ngo@test.org" || claims.Role != "ngo" {
		t.Fatalf("claims mismatch: %+v", claims)
	}
}

func TestTokenRejectedWithWrongSecret(t *testing.T) {
	token, _ := GenerateToken(testSecret, 1, "a@b.c", "donor")
	if _, err := ParseToken("different-secret", token); err == nil {
		t.Fatal("token must fail with wrong secret")
	}
}

func TestTokenRejectedWhenInvalid(t *testing.T) {
	if _, err := ParseToken(testSecret, "not.a.jwt"); err == nil {
		t.Fatal("garbage token must be rejected")
	}
}
