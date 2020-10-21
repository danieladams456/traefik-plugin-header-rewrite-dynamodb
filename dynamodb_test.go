package traefik_plugin_header_rewrite_dynamodb

import "testing"

type getTest struct {
	key         string
	expected    string
	expectError bool
}

func TestGet(t *testing.T) {
	var tests = []getTest{
		{"key1", "ddb-value1", false},
		{"key2", "ddb-value2", false},
		{"missingValue", "", true},
		{"wrongTypeValue", "", true},
	}

	for _, test := range tests {
		ans, err := get(test.key)
		if err != nil && !test.expectError {
			t.Errorf("got unexpected error: %v", err)
		}
		if err == nil && test.expectError {
			t.Errorf("did not get expected error: %v", err)
		}
		if ans != test.expected {
			t.Errorf("got: %v, want: %v", ans, test.expected)
		}
	}
}
