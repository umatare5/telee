package cli

import (
	"fmt"
	"os"
	"telee/internal/application"
	"telee/internal/config"
	"telee/internal/framework"
	"telee/internal/infrastructure"

	"github.com/urfave/cli/v2"
)

// Preset parameters for this command
const (
	cmdName      string = "telee"
	cmdUsage     string = "One-line pseudo terminal"
	cmdUsageText string = "telee -H HOSTNAME -C COMMAND [options...]"
	cmdVersion   string = "1.5.0"
)

// Start executes this command
func Start() {
	cmd := &cli.App{
		Name:      cmdName,
		HelpName:  cmdName,
		Usage:     cmdUsage,
		UsageText: cmdUsageText,
		Version:   cmdVersion,
		Flags:     registerFlags(),
		Action: func(ctx *cli.Context) error {
			c := config.New(ctx)
			r := infrastructure.New(&c)
			u := application.New(&c, &r)
			exec := framework.New(&c, &r, &u)
			exec.Run()
			return nil
		},
	}

	err := cmd.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}
}
