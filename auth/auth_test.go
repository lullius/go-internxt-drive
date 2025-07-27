package auth

import "testing"

func TestDecryptTextWithKey(t *testing.T) {
	want := "this was encrypted"
	got, _ := decryptTextWithKey("53616c7465645f5f78753b7d2417893c685726f8c5d5778e2e591e528f657270d3f94ff26ccc440112d94eea8308ce10", "secret")
	if want != got {
		t.Fatalf("Wanted " + want + ", but got " + got)
	}
}

func TestEncryptTextWithKey(t *testing.T) {
	want := "this was encrypted"
	got_encrypted, _ := encryptTextWithKey("this was encrypted", "secret")
	got_decrypted, _ := decryptTextWithKey(got_encrypted, "secret")
	if want != got_decrypted {
		t.Fatalf("Wanted " + want + ", but got " + got_decrypted)
	}
}
