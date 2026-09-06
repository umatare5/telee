package cli

import (
	"github.com/urfave/cli/v3"

	"github.com/umatare5/telee/internal/domain"
)

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
	flags = append(flags, registerHostKeyPathFlag()...)
	return flags
}

func registerHostNameFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     domain.HostnameFlagName,
			Usage:    domain.HostnameFlagUsage,
			Aliases:  domain.HostnameFlagAliases,
			Sources:  cli.EnvVars(domain.HostnameFlagEnvVars),
			Required: domain.HostnameFlagRequired,
		},
	}
}

func registerPortFlag() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:    domain.PortFlagName,
			Usage:   domain.PortFlagUsage,
			Aliases: domain.PortFlagAliases,
			Value:   int(domain.PortFlagDefaultValue),
		},
	}
}

func registerTimeoutFlag() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:    domain.TimeoutFlagName,
			Usage:   domain.TimeoutFlagUsage,
			Aliases: domain.TimeoutFlagAliases,
			Value:   int(domain.TimeoutFlagDefaultValue),
		},
	}
}

func registerCommandFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     domain.CommandFlagName,
			Usage:    domain.CommandFlagUsage,
			Aliases:  domain.CommandFlagAliases,
			Sources:  cli.EnvVars(domain.CommandFlagEnvVars),
			Required: domain.CommandFlagRequired,
		},
	}
}

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

func registerUsernameFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    domain.UsernameFlagName,
			Usage:   domain.UsernameFlagUsage,
			Aliases: domain.UsernameFlagAliases,
			Value:   domain.UsernameFlagDefaultValue,
			Sources: cli.EnvVars(domain.UsernameFlagEnvVars),
		},
	}
}

func registerPasswordFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    domain.PasswordFlagName,
			Usage:   domain.PasswordFlagUsage,
			Aliases: domain.PasswordFlagAliases,
			Value:   domain.PasswordFlagDefaultValue,
			Sources: cli.EnvVars(domain.PasswordFlagEnvVars),
		},
	}
}

func registerPrivPasswordFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    domain.PrivPasswordFlagName,
			Usage:   domain.PrivPasswordFlagUsage,
			Aliases: domain.PrivPasswordFlagAliases,
			Value:   domain.PrivPasswordFlagDefaultValue,
			Sources: cli.EnvVars(domain.PrivPasswordFlagEnvVars),
		},
	}
}

func registerHostKeyPathFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    domain.HostKeyPathFlagName,
			Usage:   domain.HostKeyPathFlagUsage,
			Aliases: domain.HostKeyPathFlagAliases,
			Value:   domain.HostKeyPathFlagDefaultValue,
			Sources: cli.EnvVars(domain.HostKeyPathFlagEnvVars),
		},
	}
}
