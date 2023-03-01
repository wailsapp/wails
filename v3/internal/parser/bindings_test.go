package parser

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

const expectedGreetService = `function GreetService(method) {
    return {
        packageName: "main",
        serviceName: "GreetService",
        methodName: method,
        args: Array.prototype.slice.call(arguments, 1),
    };
}

/**
 * GreetService.Greet
 * Greet someone
 * @param name {string}
 * @returns {Promise<string>}
 */
function Greet(name) {
    return wails.Call(GreetService("Greet", name));
}

window.go = window.go || {};
Object.window.go.main = {
    GreetService: {
        Greet,
    }
};
`

func TestGenerateGreetService(t *testing.T) {
	parsedMethods := map[string]map[string][]*BoundMethod{
		"main": {
			"GreetService": {
				{
					Name:       "Greet",
					DocComment: "Greet someone\n",
					Inputs: []*Parameter{
						{
							Name: "name",
							Type: &ParameterType{
								Name: "string",
							},
						},
					},
					Outputs: []*Parameter{
						{
							Name: "",
							Type: &ParameterType{
								Name: "string",
							},
						},
					},
				},
			},
		},
	}
	got := GenerateBindings(parsedMethods)
	if diff := cmp.Diff(expectedGreetService, got); diff != "" {
		t.Fatalf("GenerateService() mismatch (-want +got):\n%s", diff)
	}
}
