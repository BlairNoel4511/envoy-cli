package envfile

import (
	"testing"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	plaintext := "super-secret-value"
	passphrase := "my-passphrase"

	encrypted, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}
	if encrypted == plaintext {
		t.Error("encrypted output should not equal plaintext")
	}

	decrypted, err := Decrypt(encrypted, passphrase)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}
	if decrypted != plaintext {
		t.Errorf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestEncrypt_ProducesUniqueOutputs(t *testing.T) {
	plaintext := "value"
	passphrase := "pass"
	a, _ := Encrypt(plaintext, passphrase)
	b, _ := Encrypt(plaintext, passphrase)
	if a == b {
		t.Error("two encryptions of the same value should differ due to random nonce")
	}
}

func TestDecrypt_WrongPassphrase(t *testing.T) {
	encrypted, _ := Encrypt("secret", "correct-pass")
	_, err := Decrypt(encrypted, "wrong-pass")
	if err == nil {
		t.Error("expected error when decrypting with wrong passphrase")
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	_, err := Decrypt("not-valid-base64!!!", "pass")
	if err != ErrInvalidCiphertext {
		t.Errorf("expected ErrInvalidCiphertext, got %v", err)
	}
}

func TestDecrypt_TruncatedCiphertext(t *testing.T) {
	_, err := Decrypt("dG9vc2hvcnQ=", "pass") // base64 of "tooshort"
	if err != ErrInvalidCiphertext {
		t.Errorf("expected ErrInvalidCiphertext, got %v", err)
	}
}
