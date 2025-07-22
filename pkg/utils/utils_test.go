package utils

import "testing"

func TestStringToBytesGigabytes(t *testing.T) {
	sizeBytes, err := StringToBytes("4G")
	if err != nil {
		t.Fatal(err)
	}
	if sizeBytes != 4*1024*1024*1024 {
		t.Fatal("incorrectly parsed size")
	}
}

func TestStringToBytesMegabytes(t *testing.T) {
	sizeBytes, err := StringToBytes("256M")
	if err != nil {
		t.Fatal(err)
	}
	if sizeBytes != 256*1024*1024 {
		t.Fatal("incorrectly parsed size")
	}
}

func TestStringToBytesTerabytes(t *testing.T) {
	_, err := StringToBytes("2T")
	if err == nil {
		t.Fatal("Terabytes should not be supported")
	}
}

func TestStringToBytesKilobytes(t *testing.T) {
	_, err := StringToBytes("1024K")
	if err == nil {
		t.Fatal("Kilobytes should not be supported")
	}
}

func TestStringToBytesUnknown(t *testing.T) {
	_, err := StringToBytes("1024A")
	if err == nil {
		t.Fatal("Unknown unit should not be parsed")
	}
}

func TestStringToBytesExponent(t *testing.T) {
	// GO only supports exponents for floats
	_, err := StringToBytes("10E4")
	if err == nil {
		t.Fatal("Exponents should not be supported")
	}
}

func TestStringToBytes(t *testing.T) {
	sizeBytes, err := StringToBytes("256")
	if err != nil {
		t.Fatal(err)
	}
	if sizeBytes != 256 {
		t.Fatal("incorrectly parsed size")
	}
}

func TestIsPrimitive(t *testing.T) {
	if !IsPrimitive(1) {
		t.Fatal("int should be primitive")
	}
	if !IsPrimitive("test") {
		t.Fatal("string should be primitive")
	}
	if !IsPrimitive(true) {
		t.Fatal("boolean should be primitive")
	}
	if IsPrimitive([]string{"test"}) {
		t.Fatal("string slice should not be primitive")
	}
}
