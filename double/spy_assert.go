package double

import (
	"fmt"
	"strings"

	"github.com/krostar/test"

	"github.com/krostar/cli"
)

// AssertCommandRecords validates the recorded method calls for a specific command using a custom check function.
// It collects all SpyCommandRecord entries for the specified command path and passes them to the check function.
// The cmdPath parameter specifies the path to the command in the CLI tree (e.g., ["root", "sub"]).
// The check function receives all records for that command and should return an error if validation fails.
// This method is useful for complex assertions that require custom logic beyond simple method counting or sequencing.
func (spy *Spy) AssertCommandRecords(t test.TestingT, cmdPath []string, check func([]SpyCommandRecord) error) {
	var records []SpyCommandRecord
	spy.ForEachCommandRecords(func(tree []*cli.CLI, record SpyCommandRecord) {
		if spy.cmdPathMatchTree(cmdPath, tree) {
			records = append(records, record)
		}
	})

	err := check(records)
	test.Assert(t, err == nil, "record expectations failed for %s: %v", strings.Join(cmdPath, "."), err)
}

// AssertCommandMethodCalled verifies that a specific method was called on a command.
// The cmdPath parameter specifies the path to the command in the CLI tree (e.g., ["root", "sub"]).
// The methodName parameter is the name of the method to check (e.g., "Execute", "Description").
// If once is true, the method must have been called exactly once. If once is false, the method
// must have been called at least once. The test will fail if the assertion is not met.
func (spy *Spy) AssertCommandMethodCalled(t test.TestingT, cmdPath []string, methodName string, once bool) {
	count := spy.CountCommandMethodCalls(cmdPath, methodName)

	if once {
		test.Assert(t, count == 1, "expected %s method %s to be called once but called %d times", strings.Join(cmdPath, "."), methodName, count)
	} else {
		test.Assert(t, count >= 1, "expected %s method %s to be called at least once but was not called", strings.Join(cmdPath, "."), methodName)
	}
}

// AssertCommandMethodSequence verifies that methods were called in the specified order for a command path.
// The command path is specified by the names of the commands in the tree, and methodNames contains
// the sequence of method calls that is expected.
//
// This is useful for verifying the exact sequence of operations performed during command execution,
// such as ensuring hooks, flags, and execution happen in the correct order.
func (spy *Spy) AssertCommandMethodSequence(t test.TestingT, cmdPath []string, expectedMethods ...string) {
	spy.AssertCommandRecords(t, cmdPath, func(records []SpyCommandRecord) error {
		if len(records) != len(expectedMethods) {
			return fmt.Errorf("expected %d methods but found %d", len(expectedMethods), len(records))
		}

		for i, expectedMethod := range expectedMethods {
			if records[i].Method != expectedMethod {
				return fmt.Errorf("expected command at position %d for command %s to be method %s but found %s", i, strings.Join(cmdPath, "."), expectedMethod, records[i].Method)
			}
		}

		return nil
	})
}
