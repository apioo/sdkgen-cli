package cmd

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestValidateConventions(t *testing.T) {
	tests := []struct {
		name     string
		spec     string
		expected []string
	}{
		{
			name: "valid",
			spec: `{
				"operations": {
					"records.create": {
						"path": "/v0/:baseId", "method": "POST", "return": {}, "throws": [], "description": "", "stability": 1, "security": [], "authorization": true,
						"arguments": {
							"baseId": {"in": "path", "schema": {"type": "string"}},
							"payload": {"in": "body", "schema": {"type": "reference", "target": "Record"}}
						}
					}
				},
				"definitions": {"Record": {}, "Record_Collection": {}}
			}`,
			expected: nil,
		},
		{
			name: "invalid",
			spec: `{
				"operations": {
					"getRecord": {
						"path": "/v0/{baseId}/:recordId", "method": "GET",
						"arguments": {
							"tableId": {"in": "path"},
							"a": {"in": "body"},
							"b": {"in": "body"}
						}
					},
					"records.get": {
						"path": "/v0/:recordId", "method": "GET", "return": {}, "throws": [], "description": "", "stability": 1, "security": [], "authorization": true,
						"arguments": {"recordId": {"in": "query"}}
					}
				},
				"definitions": {"badName": {}}
			}`,
			expected: []string{
				"operations.getRecord: name should be 'resource.action' in camelCase (e.g. records.getAll)",
				"operations.getRecord: missing 'return'",
				"operations.getRecord: missing 'throws'",
				"operations.getRecord: missing 'description'",
				"operations.getRecord: missing 'stability'",
				"operations.getRecord: missing 'security'",
				"operations.getRecord: missing 'authorization'",
				"operations.getRecord: path uses {param} style, use :param instead (/v0/{baseId}/:recordId)",
				"operations.getRecord: path parameter ':recordId' has no argument",
				"operations.getRecord: path argument 'tableId' does not appear in path",
				"operations.getRecord: more than one body argument (a, b)",
				"operations.getRecord: GET operation has a body argument",
				"operations.records.get: argument 'recordId' must have \"in\": \"path\"",
				"types.badName: name should be Pascal_Snake case (e.g. Record_Collection)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var apiMap map[string]interface{}
			if err := json.Unmarshal([]byte(tt.spec), &apiMap); err != nil {
				t.Fatal(err)
			}

			actual := validateConventions(apiMap)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("unexpected problems:\nactual:   %q\nexpected: %q", actual, tt.expected)
			}
		})
	}
}
