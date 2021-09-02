package cobra

import (
	"context"
	"fmt"
	"time"

	"github.com/krostar/cli/app"
)

type commandVersion struct{}

func (cmd commandVersion) Description() string { return "Print version and exit" }

func (cmd commandVersion) Execute(_ context.Context, _ []string, _ []string) error {
	fmt.Println(fmt.Sprintf("%s, compiled %s", app.Version(), app.BuiltAt().Local().Format(time.RFC3339)))
	return nil
}
