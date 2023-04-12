package char

import (
	"testing"
)

func TestIsPrint(t *testing.T) {
	if IsPrint(' ') {
		t.Errorf("IsPrint(' ') != false")
	}

	if !IsPrint('a') {
		t.Errorf("IsPrint('a') != true")
	}
}

func TestIsDigit(t *testing.T) {
	if IsDigit('a') {
		t.Errorf("IsDigit('a') != false")
	}

	if !IsDigit('0') {
		t.Errorf("IsDigit('0') != true")
	}
}

func TestIsHexDigit(t *testing.T) {
	if IsHexDigit('Z') {
		t.Errorf("IsHexDigit('Z') != false")
	}

	if !IsHexDigit('A') {
		t.Errorf("IsHexDigit('A') != true")
	}
}

func TestIsOctalDigit(t *testing.T) {
	if IsOctalDigit('9') {
		t.Errorf("IsOctalDigit('9') != false")
	}

	if !IsOctalDigit('7') {
		t.Errorf("IsOctalDigit('7') != true")
	}
}

func TestIsIdentStart(t *testing.T) {
	if IsIdentStart('0') {
		t.Errorf("IsIdentStart('0') != false")
	}

	if !IsIdentStart('a') {
		t.Errorf("IsIdentStart('a') != true")
	}
}

func TestIsIdentPart(t *testing.T) {
	if IsIdentPart('-') {
		t.Errorf("IsIdentPart('-') != false")
	}

	if !IsIdentPart('0') {
		t.Errorf("IsIdentPart('0') != true")
	}
}
