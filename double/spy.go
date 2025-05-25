package double

import (
	"container/list"
	"strings"
	"sync"

	"github.com/krostar/test"

	"github.com/krostar/cli"
)

// SpyCLI creates a spy wrapper around a CLI instance to track method calls.
// It returns both the spy object for assertions and the wrapped CLI for use in tests.
// The spy records all method calls made on any command in the CLI tree, allowing
// verification of call sequences and parameters.
func SpyCLI(c *cli.CLI) (*Spy, *cli.CLI) {
	spy := &Spy{commands: make(map[*list.Element][]*cli.CLI)}
	return spy, wrapCLIWithSpy(spy, nil, c)
}

// wrapCLIWithSpy recursively wraps a CLI tree with spy interceptors.
// It creates a new CLI tree where each command is wrapped with spy functionality,
// maintaining the original structure but intercepting all method calls to record them.
func wrapCLIWithSpy(spy *Spy, tree []*cli.CLI, c *cli.CLI) *cli.CLI {
	spied := new(cli.CLI)

	if c.Command != nil {
		spied.Name = c.Name
		spied.Command = reduceWrappedToUnderlyingInterface(c.Command, &spyAllInterfaces{
			underlying: c.Command,
			saveRecord: func(record SpyCommandRecord) {
				spy.m.Lock()
				defer spy.m.Unlock()
				spy.commands[spy.records.PushBack(record)] = append(tree, spied)
			},
		})
	}

	if len(c.SubCommands) > 0 {
		spied.SubCommands = make([]*cli.CLI, len(c.SubCommands))
		for i, sub := range c.SubCommands {
			spied.SubCommands[i] = wrapCLIWithSpy(spy, append(tree, spied), sub)
		}
	}

	return spied
}

// SpyCommandRecord represents a single method call on a spy-wrapped CLI command.
// It captures the method name, input parameters, and output values to enable
// verification of method calls during testing.
type SpyCommandRecord struct {
	Method  string
	Inputs  []any
	Outputs []any
}

// Spy records and provides access to method call records on CLI commands.
// It maintains thread-safe access to records of all method calls made on
// commands in the wrapped CLI tree.
type Spy struct {
	m        sync.RWMutex                 // Mutex for thread-safe access to records
	records  list.List                    // Linked list of all method call records
	commands map[*list.Element][]*cli.CLI // Maps records to their command tree path
}

// DebugMethods logs all recorded method calls to the test output for debugging purposes.
// For each recorded method call, it prints the command path and method name in the format
// "[command.path] MethodName called". This is useful for understanding the sequence of
// method calls during test debugging.
func (spy *Spy) DebugMethods(t test.TestingT) {
	spy.ForEachCommandRecords(func(tree []*cli.CLI, record SpyCommandRecord) {
		var cmdPath []string
		for _, cmd := range tree {
			cmdPath = append(cmdPath, cmd.Name)
		}
		t.Logf("[%s] %s called", strings.Join(cmdPath, "."), record.Method)
	})
}

// ForEachCommandRecords iterates through all recorded method calls in chronological order.
// For each record, it calls the provided callback with the command tree path and the record.
// This allows for custom processing or validation of the recorded method calls.
func (spy *Spy) ForEachCommandRecords(cb func(tree []*cli.CLI, record SpyCommandRecord)) {
	spy.m.RLock()
	defer spy.m.RUnlock()

	for elem := spy.records.Front(); elem != nil; elem = elem.Next() {
		record, ok := elem.Value.(SpyCommandRecord)
		if !ok {
			panic("expected SpyCommandRecord")
		}
		cb(spy.commands[elem], record)
	}
}

// CountCommandMethodCalls returns the number of times a specific method was called on a command.
// The cmdPath parameter specifies the path to the command in the CLI tree (e.g., ["root", "sub"])
// and methodName is the name of the method to count (e.g., "Execute", "Description").
// Returns 0 if the method was never called or the command path doesn't exist.
func (spy *Spy) CountCommandMethodCalls(cmdPath []string, methodName string) uint {
	var count uint

	spy.ForEachCommandRecords(func(tree []*cli.CLI, record SpyCommandRecord) {
		if !spy.cmdPathMatchTree(cmdPath, tree) {
			return
		}

		if record.Method == methodName {
			count++
		}
	})

	return count
}

func (*Spy) cmdPathMatchTree(cmdPath []string, tree []*cli.CLI) bool {
	if len(tree) != len(cmdPath) {
		return false
	}

	match := true
	for i, name := range cmdPath {
		if tree[i].Name != name {
			match = false
			break
		}
	}
	return match
}
