package http

import "testing"

func TestHappyReadStr(t *testing.T) {
	testBS := []byte("hello")
	str, err := readBody(testBS)
	if err != nil {
		t.Fatalf("unexpected error %s", err.Error())
	}
	if str != "hello" {
		t.Fatalf("Read Body failed, expected %s, got %s", "hello", str)
	}

}

// empty byte slice
func TestSadRead0(t *testing.T) {
	testBS := []byte("")
	_, err := readBody(testBS)
	if err == nil {
        t.Fatalf("expected error %s", "Empty byte slice")
	}
}

// invalid byte sequence 
func TestSadInvalidByteSeq(t *testing.T) {
	invalidBytes := []byte{0xC0, 0x80}
	_, err := readBody(invalidBytes)
    if err == nil {
        t.Fatalf("expected error %s", "Invalid byte slice")
    }
}
