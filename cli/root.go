package cli

import (
	"context"
	"log"
	"os"

	"github.com/umatare5/telee/internal/application"
	"github.com/umatare5/telee/internal/config"
	"github.com/umatare5/telee/internal/framework"
	"github.com/umatare5/telee/internal/infrastructure"

	"github.com/urfave/cli/v3"
)

// Preset parameters for this command
const (
	cmdName      string = "telee"
	cmdUsage     string = "One-line command executor"
	cmdUsageText string = "telee -H HOSTNAME -C COMMAND [options...]"
)

// Start executes this command
func Start() {
	cmd := &cli.Command{
		Name:      cmdName,
		Usage:     cmdUsage,
		UsageText: cmdUsageText,
		Version:   getVersion(),
		Flags:     registerFlags(),
		Action: func(ctx context.Context, cli *cli.Command) error {
			c := config.New(cli)
			r := infrastructure.New(&c)
			u := application.New(&c, &r)
			exec := framework.New(&c, &r, &u)
			exec.Run()
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
