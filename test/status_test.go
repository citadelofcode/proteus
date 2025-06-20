package test

import (
	"testing"
	"strings"
	"github.com/citadelofcode/proteus/internal"
)

// Test case to validate the status message received through the GetStatusMessage() of StatusCode.
func Test_GetStatusMessage(t *testing.T) {
	testCases := []struct {
		Name string
		IpStatus internal.StatusCode
		ExpOutput string
	} {
		{ "A valid status code with an associated message", internal.StatusOK, "OK" },
		{ "An invalid status code with no available message", internal.StatusCode(600), "" },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			message := testCase.IpStatus.GetStatusMessage()
			if strings.EqualFold(message, testCase.ExpOutput) {
				t.Logf("The expected status message [%s] matches the received status message [%s].", testCase.ExpOutput, message)
			} else {
				t.Errorf(internal.TextColor.Red("The expected status message [%s] does not match the received status message [%s]."), testCase.ExpOutput, message)
			}
		})
	}
}
