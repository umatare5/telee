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
	flags = append(flags, registerExecPlatformFlag()...)
	flags = append(flags, registerEnableModeFlag()...)
	flags = append(flags, registerRedundantModeFlag()...)
	flags = append(flags, registerSecureModeFlag()...)
	flags = append(flags, registerDefaultPrivModeFlag()...)
	flags = append(flags, registerUsernameFlag()...)
	flags = append(flags, registerPasswordFlag()...)
	flags = append(flags, registerPrivPasswordFlag()...)
	return flags
}

// Declare hostname flag
func registerHostNameFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     domain.HostnameFlagName,
			Usage:    domain.HostnameFlagUsage,
			Aliases:  domain.HostnameFlagAliases,
			EnvVars:  domain.HostnameFlagEnvVars,
			Required: domain.HostnameFlagRequired,
		},
	}
}

// Declare port flag
func registerPortFlag() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:    domain.PortFlagName,
			Usage:   domain.PortFlagUsage,
			Aliases: domain.PortFlagAliases,
			Value:   domain.PortFlagDefaultValue,
		},
	}
}

// Declare timeout flag
func registerTimeoutFlag() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:    domain.TimeoutFlagName,
			Usage:   domain.TimeoutFlagUsage,
			Aliases: domain.TimeoutFlagAliases,
			Value:   domain.TimeoutFlagDefaultValue,
		},
	}
}

// Declare command flag
func registerCommandFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     domain.CommandFlagName,
			Usage:    domain.CommandFlagUsage,
			Aliases:  domain.CommandFlagAliases,
			EnvVars:  domain.CommandFlagEnvVars,
			Required: domain.CommandFlagRequired,
		},
	}
}

// Declare exec-platform flag
func registerExecPlatformFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    domain.ExecPlatformFlagName,
			Usage:   domain.ExecPlatformFlagUsage,
			Aliases: domain.ExecPlatformFlagAliases,
			Value:   domain.ExecPlatformFlagDefaultValue,
		},
	}
}

// Declare enable-mode flag
func registerEnableModeFlag() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    domain.EnableModeFlagName,
			Usage:   domain.EnableModeFlagUsage,
			Aliases: domain.EnableModeFlagAliases,
			Value:   domain.EnableModeFlagDefaultValue,
		},
	}
}

// Declare redundant-mode flag
func registerRedundantModeFlag() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    domain.RedundantModeFlagName,
			Usage:   domain.RedundantModeFlagUsage,
			Aliases: domain.RedundantModeFlagAliases,
			Value:   domain.RedundantModeFlagDefaultValue,
		},
	}
}

// Declare secure-mode flag
func registerSecureModeFlag() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    domain.SecureModeFlagName,
			Usage:   domain.SecureModeFlagUsage,
			Aliases: domain.SecureModeFlagAliases,
			Value:   domain.SecureModeFlagDefaultValue,
		},
	}
}

// Declare priv-mode flag
func registerDefaultPrivModeFlag() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    domain.DefaultPrivModeFlagName,
			Usage:   domain.DefaultPrivModeFlagUsage,
			Aliases: domain.DefaultPrivModeFlagAliases,
			Value:   domain.DefaultPrivModeFlagDefaultValue,
		},
	}
}

// Declare username flag
func registerUsernameFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    domain.UsernameFlagName,
			Usage:   domain.UsernameFlagUsage,
			Aliases: domain.UsernameFlagAliases,
			Value:   domain.UsernameFlagDefaultValue,
			EnvVars: domain.UsernameFlagEnvVars,
		},
	}
}

// Declare password flag
func registerPasswordFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    domain.PasswordFlagName,
			Usage:   domain.PasswordFlagUsage,
			Aliases: domain.PasswordFlagAliases,
			Value:   domain.PasswordFlagDefaultValue,
			EnvVars: domain.PasswordFlagEnvVars,
		},
	}
}

// Declare priv-password flag
func registerPrivPasswordFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    domain.PrivPasswordFlagName,
			Usage:   domain.PrivPasswordFlagUsage,
			Aliases: domain.PrivPasswordFlagAliases,
			Value:   domain.PrivPasswordFlagDefaultValue,
			EnvVars: domain.PrivPasswordFlagEnvVars,
		},
	}
}
