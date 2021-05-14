package config

import (
	"fmt"
	"log"
	"telee/internal/domain"
	"telee/pkg/errors"

	"github.com/jinzhu/configor"
	"github.com/urfave/cli/v2"
)

const (
	infoUserameIgnored    = "[INFO] username is ignored. It's not supported."
	infoEnableModeIgnored = "[INFO] enable-mode is ignored. It's not supported."
)

// Config struct
type Config struct {
	Hostname        string
	Port            int
	Timeout         int
	ExecPlatform    string
	EnableMode      bool
	HAMode          bool
	SecureMode      bool
	DefaultPrivMode bool
	Command         string
	Username        string
	Password        string
	PrivPassword    string
}

// New returns Config struct
func New(ctx *cli.Context) Config {
	cfg := Config{
		Hostname:        ctx.String(domain.HostnameFlagName),
		Port:            ctx.Int(domain.PortFlagName),
		Timeout:         ctx.Int(domain.TimeoutFlagName),
		Command:         ctx.String(domain.CommandFlagName),
		ExecPlatform:    ctx.String(domain.ExecPlatformFlagName),
		EnableMode:      ctx.Bool(domain.EnableModeFlagName),
		HAMode:          ctx.Bool(domain.HAModeFlagName),
		DefaultPrivMode: ctx.Bool(domain.DefaultPrivModeFlagName),
		SecureMode:      ctx.Bool(domain.SecureModeFlagName),
		Username:        ctx.String(domain.UsernameFlagName),
		Password:        ctx.String(domain.PasswordFlagName),
		PrivPassword:    ctx.String(domain.PrivPasswordFlagName),
	}

	err := configor.New(&configor.Config{}).Load(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = checkArguments(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

func checkArguments(cfg *Config) error {
	if cfg.Port == 0 {
		completePortNumber(cfg)
	}
	if !isValidExecPlatform(cfg.ExecPlatform) {
		return errors.ErrInvalidPlatform
	}
	if cfg.EnableMode && cfg.DefaultPrivMode {
		return errors.ErrUnsupportedModeSet
	}
	if cfg.HAMode && !isUsableHAMode(cfg.ExecPlatform) {
		return errors.ErrUnsupportedHAMode
	}
	if cfg.SecureMode && !isUsableSecureMode(cfg.ExecPlatform) {
		return errors.ErrUnsupportedSecureMode
	}
	if !cfg.SecureMode && !isUsableUnsecureMode(cfg.ExecPlatform) {
		return errors.ErrUnsupportedUnsecureMode
	}
	if !cfg.EnableMode && !isExpandableTermLength(cfg.ExecPlatform) {
		return errors.ErrTermLengthIsEnforced
	}
	if cfg.Hostname == domain.EmptyString {
		return errors.ErrMissingHostname
	}
	if cfg.Command == domain.EmptyString {
		return errors.ErrMissingCommand
	}
	if cfg.Username == domain.UsernameFlagDefaultValue {
		return errors.ErrMissingUsername
	}
	if cfg.Password == domain.PasswordFlagDefaultValue {
		return errors.ErrMissingPassword
	}
	if cfg.EnableMode && !hasPrivPassword(cfg.PrivPassword) {
		return errors.ErrMissingPrivPassword
	}
	if cfg.EnableMode && !isUsableEnableMode(cfg.ExecPlatform) {
		fmt.Println(infoEnableModeIgnored)
	}
	if cfg.Username != domain.EmptyString && !isUsableUsername(cfg.ExecPlatform) {
		fmt.Println(infoUserameIgnored)
	}

	return nil
}

func completePortNumber(cfg *Config) bool {
	if cfg.SecureMode {
		cfg.Port = domain.SSHPort
	} else {
		cfg.Port = domain.TelnetPort
	}
	return true
}

func isValidExecPlatform(platform string) bool {
	for _, p := range domain.Platforms {
		if platform == p {
			return true
		}
	}
	return false
}

func isUsableHAMode(platform string) bool {
	if platform == domain.ASASoftwarePlatformName {
		return true
	}
	if platform == domain.ScreenOSPlatformName {
		return true
	}
	return false
}

func isUsableSecureMode(platform string) bool {
	if platform == domain.IOSPlatformName {
		return true
	}
	if platform == domain.JunOSPlatformName {
		return true
	}
	return false
}

func isUsableUnsecureMode(platform string) bool {
	return platform != domain.JunOSPlatformName
}

func isExpandableTermLength(platform string) bool {
	return platform != domain.ASASoftwarePlatformName
}

func hasPrivPassword(password string) bool {
	return password != domain.PrivPasswordFlagDefaultValue
}

func isUsableEnableMode(platform string) bool {
	if platform == domain.AireOSPlatformName {
		return false
	}
	if platform == domain.AlliedWarePlatformName {
		return false
	}
	if platform == domain.ScreenOSPlatformName {
		return false
	}
	return true
}

func isUsableUsername(platform string) bool {
	return platform != domain.YamahaOSPlatformName
}
