package network

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestValidMac(t *testing.T) {
	data := []byte{0, 0, 0, 0, 0, 0}
	mac := NewMac(data)

	if mac == nil {
		t.Error("Expected a valid Mac. Got nil")
	}
}

func TestInvalidMacNotEnoughData(t *testing.T) {
	data := []byte{0, 0, 0, 0, 0}
	mac := NewMac(data)

	if mac != nil {
		t.Error("Expected a invalid Mac. Got something different than nil")
	}
}

func TestValidStringRepresentation(t *testing.T) {
	data := []byte{0x00, 0x20, 0x40, 0x60, 0x80, 0xa0}
	mac := NewMac(data)

	expected := "00:20:40:60:80:a0"
	actual := mac.String()
	if strings.Compare(actual, expected) != 0 {
		t.Errorf("Invalid Mac representation. Actual: %v - Expected: %v", actual, expected)
	}
}

func TestMacIncreaseNonBorderCase(t *testing.T) {
	data := []byte{0x00, 0x20, 0x40, 0x60, 0x80, 0xa0}
	mac := NewMac(data)
	result := mac.Increase()

	expected := "00:20:40:60:80:a1"
	actual := mac.String()
	if result != nil || strings.Compare(actual, expected) != 0 {
		t.Errorf("Mac was not increased. Actual: %v - Expected: %v", actual, expected)
	}
}

func TestMacIncreaseBorderCase(t *testing.T) {
	data := []byte{0x00, 0xff, 0xff, 0xff, 0xff, 0xff}
	mac := NewMac(data)
	result := mac.Increase()

	expected := "01:00:00:00:00:00"
	actual := mac.String()
	if result != nil || strings.Compare(actual, expected) != 0 {
		t.Errorf("Mac was not increased. Actual: %v - Expected: %v", actual, expected)
	}
}

func TestMacIncreaseMaxMac(t *testing.T) {
	data := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	mac := NewMac(data)
	result := mac.Increase()

	expected := "ff:ff:ff:ff:ff:ff"
	actual := mac.String()
	if result == nil || strings.Compare(actual, expected) != 0 {
		t.Errorf("Mac was increased!!. Actual: %v - Expected: %v", actual, expected)
	}
}

func TestMacCompareMacCaseEqual(t *testing.T) {
	mac1 := NewMac([]byte{0x00, 0x20, 0x40, 0x60, 0x80, 0xa0})
	mac2 := NewMac([]byte{0x00, 0x20, 0x40, 0x60, 0x80, 0xa0})

	if mac1.Compare(mac2) != 0 {
		t.Errorf("Mac addresses are not equal!!. Actual: %v - Expected: %v", mac1, mac2)
	}
}

func TestMacCompareMacCaseGreater(t *testing.T) {
	mac1 := NewMac([]byte{0x00, 0x20, 0x40, 0x60, 0x80, 0xa1})
	mac2 := NewMac([]byte{0x00, 0x20, 0x40, 0x60, 0x80, 0xa0})

	if mac1.Compare(mac2) != 1 {
		t.Errorf("Mac addresses are not equal!!. Actual: %v - Expected: %v", mac1, mac2)
	}
}

func TestMacCompareMacCaseLesser(t *testing.T) {
	mac1 := NewMac([]byte{0x00, 0x20, 0x40, 0x60, 0x80, 0xa1})
	mac2 := NewMac([]byte{0x01, 0x20, 0x40, 0x60, 0x80, 0xa0})

	if mac1.Compare(mac2) != -1 {
		t.Errorf("Mac addresses are not equal!!. Actual: %v - Expected: %v", mac1, mac2)
	}
}

func BenchmarkTestCompareMacCaseEqual(b *testing.B) {
	mac1 := NewMac([]byte{0x00, 0x20, 0x40, 0x60, 0x80, 0xa0})
	mac2 := NewMac([]byte{0x00, 0x20, 0x40, 0x60, 0x80, 0xa0})

	for i := 0; i < b.N; i++ {
		mac1.Compare(mac2)
	}
}

func BenchmarkDeepEqualMacCaseEqual(b *testing.B) {
	mac1 := NewMac([]byte{0x00, 0x20, 0x40, 0x60, 0x80, 0xa0})
	mac2 := NewMac([]byte{0x00, 0x20, 0x40, 0x60, 0x80, 0xa0})

	for i := 0; i < b.N; i++ {
		reflect.DeepEqual(mac1, mac2)
	}
}

func BenchmarkTestCompareMacCaseNotEqual(b *testing.B) {
	mac1 := NewMac([]byte{0x00, 0x20, 0x40, 0x60, 0x80, 0xa1})
	mac2 := NewMac([]byte{0x00, 0x20, 0x40, 0x60, 0x80, 0xa0})

	for i := 0; i < b.N; i++ {
		mac1.Compare(mac2)
	}
}

func BenchmarkDeepEqualMacCaseNotEqual(b *testing.B) {
	mac1 := NewMac([]byte{0x00, 0x20, 0x40, 0x60, 0x80, 0xa1})
	mac2 := NewMac([]byte{0x00, 0x20, 0x40, 0x60, 0x80, 0xa0})

	for i := 0; i < b.N; i++ {
		reflect.DeepEqual(mac1, mac2)
	}
}

func ExampleMac_Compare() {
	mac1 := NewMac([]byte{0x00, 0x20, 0x40, 0x60, 0x80, 0xa0})
	mac2 := NewMac([]byte{0x01, 0x20, 0x40, 0x60, 0x80, 0xa0})

	if mac1.Compare(mac2) != 0 {
		fmt.Println("Macs are not equivalents")
	}
	// Output:
	// Macs are not equivalents
}

func ExampleMac_Compare_second() {
	mac1 := NewMac([]byte{0x00, 0x20, 0x40, 0x60, 0x80, 0xa0})
	mac2 := NewMac([]byte{0x00, 0x20, 0x40, 0x60, 0x80, 0xa0})

	if mac1.Compare(mac2) == 0 {
		fmt.Println("Macs are equivalents")
	}
}

// Output:
// Macs are equivalents
