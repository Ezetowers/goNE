package net_test

import (
	"ne/net"
	"reflect"
	"strings"
	"testing"
)

func TestValidIPAddress(t *testing.T) {
	data := []byte{1, 2, 3, 4}
	ipAddr := net.NewIPAddress(data)

	if ipAddr == nil {
		t.Error("Expected a valid IPAddress. Got nil")
	}
}

func TestInvalidIPAddressNotEnoughData(t *testing.T) {
	data := []byte{0, 0, 0, 0, 0}
	ipAddr := net.NewIPAddress(data)

	if ipAddr != nil {
		t.Error("Expected a invalid IPAddress. Got something different than nil")
	}
}

func TestValidIPAddressStringRepresentation(t *testing.T) {
	data := []byte{1, 2, 3, 4}
	ipAddr := net.NewIPAddress(data)

	expected := "1.2.3.4"
	actual := ipAddr.String()
	if strings.Compare(actual, expected) != 0 {
		t.Errorf("Invalid IPAddress representation. Actual: %v - Expected: %v", actual, expected)
	}
}

func TestIPAddressIncreaseNonBorderCase(t *testing.T) {
	data := []byte{192, 168, 0, 254}
	ipAddr := net.NewIPAddress(data)
	result := ipAddr.Increase()

	expected := "192.168.0.255"
	actual := ipAddr.String()
	if result != nil || strings.Compare(actual, expected) != 0 {
		t.Errorf("IPAddress was not increased. Actual: %v - Expected: %v", actual, expected)
	}
}

func TestIPAddressIncreaseBorderCase(t *testing.T) {
	data := []byte{191, 255, 255, 255}
	ipAddr := net.NewIPAddress(data)
	result := ipAddr.Increase()

	expected := "192.0.0.0"
	actual := ipAddr.String()
	if result != nil || strings.Compare(actual, expected) != 0 {
		t.Errorf("IPAddress was not increased. Actual: %v - Expected: %v", actual, expected)
	}
}

func TestIncreaseMaxIPAddress(t *testing.T) {
	data := []byte{255, 255, 255, 255}
	ipAddr := net.NewIPAddress(data)
	result := ipAddr.Increase()

	expected := "255.255.255.255"
	actual := ipAddr.String()
	if result == nil || strings.Compare(actual, expected) != 0 {
		t.Errorf("IPAddress was not increased. Actual: %v - Expected: %v", actual, expected)
	}
}

func TestCompareIPAddressCaseEqual(t *testing.T) {
	ipAddr1 := net.NewIPAddress([]byte{1, 2, 3, 4})
	ipAddr2 := net.NewIPAddress([]byte{1, 2, 3, 4})

	if ret, _ := ipAddr1.Compare(ipAddr2); ret != 0 {
		t.Errorf("IPAddresses are not equal!!. Actual: %v - Expected: %v", ipAddr1, ipAddr2)
	}
}

func TestCompareIPAddressCaseGreater(t *testing.T) {
	ipAddr1 := net.NewIPAddress([]byte{1, 3, 3, 2})
	ipAddr2 := net.NewIPAddress([]byte{1, 2, 3, 4})

	if ret, _ := ipAddr1.Compare(ipAddr2); ret == 0 {
		t.Errorf("IPAddresses are not equal!!. Actual: %v - Expected: %v", ipAddr1, ipAddr2)
	}
}

func TestCompareIPAddressCaseLesser(t *testing.T) {
	ipAddr1 := net.NewIPAddress([]byte{1, 2, 2, 2})
	ipAddr2 := net.NewIPAddress([]byte{1, 2, 4, 1})

	if ret, _ := ipAddr1.Compare(ipAddr2); ret == 0 {
		t.Errorf("IPAddresses are not equal!!. Actual: %v - Expected: %v", ipAddr1, ipAddr2)
	}
}

func BenchmarkTestCompareIPAddressCaseEqual(b *testing.B) {
	ipAddr1 := net.NewIPAddress([]byte{1, 2, 3, 4})
	ipAddr2 := net.NewIPAddress([]byte{1, 2, 3, 4})

	for i := 0; i < b.N; i++ {
		ipAddr1.Compare(ipAddr2)
	}
}

func BenchmarkDeepEqualIPAddressCaseEqual(b *testing.B) {
	ipAddr1 := net.NewIPAddress([]byte{1, 2, 3, 4})
	ipAddr2 := net.NewIPAddress([]byte{1, 2, 3, 4})

	for i := 0; i < b.N; i++ {
		reflect.DeepEqual(ipAddr1, ipAddr2)
	}
}

func BenchmarkTestCompareIPAddressCaseNotEqual(b *testing.B) {
	ipAddr1 := net.NewIPAddress([]byte{1, 2, 3, 5})
	ipAddr2 := net.NewIPAddress([]byte{1, 2, 3, 4})

	for i := 0; i < b.N; i++ {
		ipAddr1.Compare(ipAddr2)
	}
}

func BenchmarkDeepEqualIPAddressCaseNotEqual(b *testing.B) {
	ipAddr1 := net.NewIPAddress([]byte{1, 2, 3, 5})
	ipAddr2 := net.NewIPAddress([]byte{1, 2, 3, 4})

	for i := 0; i < b.N; i++ {
		reflect.DeepEqual(ipAddr1, ipAddr2)
	}
}
