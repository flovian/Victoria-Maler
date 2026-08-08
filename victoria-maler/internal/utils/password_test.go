package utils

import "testing"

func TestPasswordRoundTrip(t *testing.T) {
	hash, err := HashPassword("s3cret-pass")
	if err != nil {
		t.Fatal(err)
	}
	if hash == "s3cret-pass" {
		t.Fatal("hash must not equal plaintext")
	}
	if !CheckPassword(hash, "s3cret-pass") {
		t.Fatal("correct password should verify")
	}
	if CheckPassword(hash, "wrong-pass") {
		t.Fatal("wrong password must not verify")
	}
}

func TestPasswordUniqueSalt(t *testing.T) {
	h1, _ := HashPassword("same")
	h2, _ := HashPassword("same")
	if h1 == h2 {
		t.Fatal("bcrypt should salt each hash uniquely")
	}
}
