package testsuite

import "testing"

// AssertEquals that 2 interfaces are equals
func AssertByteArraysEquals(t *testing.T, message string, expected, given []byte) {
	if len(expected) != len(given) {
		t.Fatalf("Error during assertion: %s, expected byte array length <%v> but given length <%v>\n", message, len(expected), len(given))
	}
	for idx, g := range given {
		e := expected[idx]
		if e != g {
			t.Fatalf("Error during assertion: %s, byte array at index %v expected value <%v> but given <%v>\n", message, idx, e, g)
		}
	}
}

// AssertEquals that 2 interfaces are equals
func AssertEquals(t *testing.T, message string, expected, given interface{}) {
	if expected != given {
		t.Fatalf("Error during assertion: %s, expected value <%v> but given <%v>\n", message, expected, given)
	}
}

// AssertNotEquals that 2 interfaces are not equals
func AssertNotEquals(t *testing.T, message string, expected, given interface{}) {
	if expected == given {
		t.Fatalf("Error during assertion: %s, NOT expected value <%v> but given <%v>\n", message, expected, given)
	}
}

// AssertNotNil that 1 interfaces is not nil
func AssertNotNil(t *testing.T, message string, given interface{}) {
	if nil == given {
		t.Fatalf("Error during assertion: %s, NOT expected nil value but given <%v>\n", message, given)
	}
}

// AssertNil that 1 interface is nil
func AssertNil(t *testing.T, message string, given interface{}) {
	if nil != given {
		t.Fatalf("Error during assertion: %s, expected nil value but given <%v>\n", message, given)
	}
}


