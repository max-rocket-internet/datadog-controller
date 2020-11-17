package utils

import (
	"os"
	"testing"
)

func testEqualSlices(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestContainsString(t *testing.T) {
	strings := []string{
		"one",
		"two",
		"three",
	}

	expected := true
	actual := ContainsString(strings, "two")
	if actual != true {
		t.Errorf("Got %t, expected %t, given %v", actual, expected, strings)
	}

	expected = false
	actual = ContainsString(strings, "four")
	if actual != expected {
		t.Errorf("Got %t, expected %t, given %v", actual, expected, strings)
	}
}

func TestRemoveString(t *testing.T) {
	a := []string{
		"one",
		"two",
		"three",
	}

	expected := []string{"one", "three"}
	actual := RemoveString(a, "two")

	if !testEqualSlices(actual, expected) {
		t.Errorf("Got %v, expected %v, given %v", actual, expected, a)
	}

	expected = a
	actual = RemoveString(a, "four")

	if !testEqualSlices(actual, expected) {
		t.Errorf("Got %v, expected %v, given %v", actual, expected, a)
	}
}

func TestGetEnvString(t *testing.T) {
	os.Setenv("TEST_ENV0", "1")

	expected := "1"
	actual, _ := GetEnvString("TEST_ENV0")

	if actual != expected {
		t.Errorf("Got %v, expected %v", actual, expected)
	}

	expectedError := "Environment variable not set: TEST_ENV_DOES_NOT_EXIST"
	_, actualError := GetEnvString("TEST_ENV_DOES_NOT_EXIST")

	if actualError.Error() != expectedError {
		t.Errorf("Got %v, expected %v", actualError, expectedError)
	}
}

func TestGetEnvInt(t *testing.T) {
	os.Setenv("TEST_ENV2", "1")

	expected := 1
	actual, _ := GetEnvInt("TEST_ENV2")

	if actual != expected {
		t.Errorf("Got %v, expected %v", actual, expected)
	}

	expectedError := "Environment variable not set: TEST_ENV_DOES_NOT_EXIST"
	_, actualError := GetEnvInt("TEST_ENV_DOES_NOT_EXIST")

	if actualError.Error() != expectedError {
		t.Errorf("Got %v, expected %v", actualError, expectedError)
	}
}

func TestCheckRequiredEnvVars(t *testing.T) {
	os.Setenv("TEST_ENV3", "test")
	os.Setenv("TEST_ENV4", "test")

	actual := CheckRequiredEnvVars([]string{"TEST_ENV3", "TEST_ENV4"})

	if actual != nil {
		t.Errorf("Got %v, expected nil", actual)
	}

	expected := "Environment variable not set: TEST_ENV5"
	actual = CheckRequiredEnvVars([]string{"TEST_ENV3", "TEST_ENV5"})

	if actual.Error() != expected {
		t.Errorf("Got %v, expected %v", actual, expected)
	}
}
