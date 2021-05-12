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
	Hostname     string
	Port         int
	Timeout      int
	ExecPlatform string
	EnableMode   bool
	HAMode       bool
	Command      string
	Username     string
	Password     string
	PrivPassword string
}

// New returns Config struct
func New(ctx *cli.Context) Config {
	cfg := Config{
		Hostname:     ctx.String(domain.HostnameFlagName),
		Port:         ctx.Int(domain.PortFlagName),
		Timeout:      ctx.Int(domain.TimeoutFlagName),
		Command:      ctx.String(domain.CommandFlagName),
		ExecPlatform: ctx.String(domain.ExecPlatformFlagName),
		EnableMode:   ctx.Bool(domain.EnableModeFlagName),
		HAMode:       ctx.Bool(domain.HAModeFlagName),
		Username:     ctx.String(domain.UsernameFlagName),
		Password:     ctx.String(domain.PasswordFlagName),
		PrivPassword: ctx.String(domain.PrivPasswordFlagName),
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
	if cfg.Username == domain.UsernameFlagDefaultValue {
		return errors.ErrMissingUsername
	}
	if cfg.Password == domain.PasswordFlagDefaultValue {
		return errors.ErrMissingPassword
	}
	if !hasPrivPassword(cfg.EnableMode, cfg.PrivPassword) {
		return errors.ErrMissingPrivPassword
	}
	if !isExpandableTermLength(cfg.EnableMode, cfg.ExecPlatform) {
		return errors.ErrTermLengthIsEnforced
	}
	if cfg.Hostname == domain.EmptyString {
		return errors.ErrMissingHostname
	}
	if cfg.Command == domain.EmptyString {
		return errors.ErrMissingCommand
	}
	if isUsableHAMode(cfg.HAMode, cfg.ExecPlatform) {
		return errors.ErrUnsupportedHAMode
	}
	if !isValidExecPlatform(cfg.ExecPlatform) {
		return errors.ErrInvalidPlatform
	}
	if !isUsableEnableMode(cfg.EnableMode, cfg.ExecPlatform) {
		fmt.Println(infoEnableModeIgnored)
	}
	if !isUsableUsername(cfg.Username, cfg.ExecPlatform) {
		fmt.Println(infoUserameIgnored)
	}

	return nil
}

func isExpandableTermLength(mode bool, platform string) bool {
	if !mode && platform == domain.ASASoftwarePlatformName {
		return false
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

func hasPrivPassword(value bool, password string) bool {
	if value {
		if password == domain.PrivPasswordFlagDefaultValue {
			return false
		}
	}
	return true
}

func isUsableEnableMode(mode bool, platform string) bool {
	if mode {
		if platform == domain.AireOSPlatformName {
			return false
		}
		if platform == domain.AlliedWarePlatformName {
			return false
		}
		if platform == domain.ScreenOSPlatformName {
			return false
		}
	}
	return true
}

func isUsableHAMode(mode bool, platform string) bool {
	if mode {
		if platform == domain.ASASoftwarePlatformName {
			return true
		}
		if platform == domain.ScreenOSPlatformName {
			return true
		}
	}
	return false
}

func isUsableUsername(username string, platform string) bool {
	if username != domain.EmptyString {
		if platform == domain.YamahaOSPlatformName {
			return false
		}
	}
	return true
}
