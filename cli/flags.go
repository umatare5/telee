package cli

import (
	"telee/internal/domain"

	"github.com/urfave/cli/v2"
)

// Register all flags
func registerFlags() []cli.Flag {
	flags := []cli.Flag{}
	flags = append(flags, registerHostNameFlag()...)
	flags = append(flags, registerPortFlag()...)
	flags = append(flags, registerTimeoutFlag()...)
	flags = append(flags, registerCommandFlag()...)
	flags = append(flags, registerPlatformFlag()...)
	flags = append(flags, registerEnableModeFlag()...)
	flags = append(flags, registerUsernameFlag()...)
	flags = append(flags, registerPasswordFlag()...)
	flags = append(flags, registerPrivPasswordFlag()...)
	return flags
}

// Preset parameters for hostname flag
const (
	optHostnameName     string = "hostname"
	optHostnameUsage    string = "Set hostname or IP address."
	optHostnameRequired bool   = true
)

var (
	optHostnameAlias            = []string{"H"}
	optHostnameEnvVars []string = []string{"TELEE_HOSTNAME"}
)

// Declare hostname flag
func registerHostNameFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     optHostnameName,
			Usage:    optHostnameUsage,
			Aliases:  optHostnameAlias,
			EnvVars:  optHostnameEnvVars,
			Required: optHostnameRequired,
		},
	}
}

// Preset parameters for port flag
const (
	optPortName  string = "port"
	optPortUsage string = "Set port number."
	optPortValue int    = 23
)

var optPortAlias = []string{"P"}

// Declare port flag
func registerPortFlag() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:    optPortName,
			Usage:   optPortUsage,
			Aliases: optPortAlias,
			Value:   optPortValue,
		},
	}
}

// Preset parameters for timeout flag
const (
	optTimeoutName  string = "timeout"
	optTimeoutUsage string = "Set timeout seconds."
	optTimeoutValue int    = 5
)

var optTimeoutAlias = []string{"t"}

// Declare timeout flag
func registerTimeoutFlag() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:    optTimeoutName,
			Usage:   optTimeoutUsage,
			Aliases: optTimeoutAlias,
			Value:   optTimeoutValue,
		},
	}
}

// Preset parameters for command flag
const (
	optCommandName     string = "command"
	optCommandUsage    string = "Set a command."
	optCommandRequired bool   = true
)

var (
	optCommandAlias            = []string{"C"}
	optCommandEnvVars []string = []string{"TELEE_COMMAND"}
)

// Declare command flag
func registerCommandFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     optCommandName,
			Usage:    optCommandUsage,
			Aliases:  optCommandAlias,
			EnvVars:  optCommandEnvVars,
			Required: optCommandRequired,
		},
	}
}

// Preset parameters for platform flag
const (
	optPlatformName  string = "platform"
	optPlatformUsage string = "Set platform. Refer to README.md what to be set."
	optPlatformValue string = "ios"
)

var optPlatformAlias = []string{"x"}

// Declare platform flag
func registerPlatformFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    optPlatformName,
			Usage:   optPlatformUsage,
			Aliases: optPlatformAlias,
			Value:   optPlatformValue,
		},
	}
}

// Preset parameters for enable-mode flag
const (
	optEnableModeName  string = "enable-mode"
	optEnableModeUsage string = "Log in to privileged EXEC mode."
	optEnableModeValue bool   = false
)

var optEnableModeAlias = []string{"e", "ena", "enable"}

// Declare enable-mode flag
func registerEnableModeFlag() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    optEnableModeName,
			Usage:   optEnableModeUsage,
			Aliases: optEnableModeAlias,
			Value:   optEnableModeValue,
		},
	}
}

// Preset parameters for username flag
const (
	optUsernameName  string = "username"
	optUsernameUsage string = "Set username."
	optUsernameValue string = domain.DefaultUsernameValue
)

var (
	optUsernameAlias   []string = []string{"u"}
	optUsernameEnvVars []string = []string{"TELEE_USERNAME"}
)

// Declare username flag
func registerUsernameFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    optUsernameName,
			Usage:   optUsernameUsage,
			Aliases: optUsernameAlias,
			Value:   optUsernameValue,
			EnvVars: optUsernameEnvVars,
		},
	}
}

// Preset parameters for password flag
const (
	optPasswordName  string = "password"
	optPasswordUsage string = "Set password."
	optPasswordValue string = domain.DefaultPasswordValue
)

var (
	optPasswordAlias   []string = []string{"p"}
	optPasswordEnvVars []string = []string{"TELEE_PASSWORD"}
)

// Declare password flag
func registerPasswordFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    optPasswordName,
			Usage:   optPasswordUsage,
			Aliases: optPasswordAlias,
			Value:   optPasswordValue,
			EnvVars: optPasswordEnvVars,
		},
	}
}

// Preset parameters for priv-password flag
const (
	optPrivPasswordName  string = "priv-password"
	optPrivPasswordUsage string = "Set password to change to privileged EXEC mode."
	optPrivPasswordValue string = domain.DefaultPrivPasswordValue
)

var (
	optPrivPasswordAlias   []string = []string{"pp"}
	optPrivPasswordEnvVars []string = []string{"TELEE_PRIVPASSWORD"}
)

// Declare priv-password flag
func registerPrivPasswordFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    optPrivPasswordName,
			Usage:   optPrivPasswordUsage,
			Aliases: optPrivPasswordAlias,
			Value:   optPrivPasswordValue,
			EnvVars: optPrivPasswordEnvVars,
		},
	}
}
