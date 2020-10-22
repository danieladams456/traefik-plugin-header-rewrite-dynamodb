package dynamodb

import "testing"

type lookupTest struct {
	key         string
	expected    string
	expectError bool
}

// currently an intergration test againt a live DynamoDB instance
func TestGet(t *testing.T) {
	// setup
	r := Repository{
		TableName:      "traefik_header_lookups",
		KeyAttribute:   "key",
		ValueAttribute: "value",
	}
	r.InitSdk()

	// define tests
	var tests = []lookupTest{
		{"key1", "ddb-value1", false},
		{"key2", "ddb-value2", false},
		{"missingValue", "", true},
		{"wrongTypeValue", "", true},
	}

	for _, test := range tests {
		ans, err := r.Lookup(test.key)
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
