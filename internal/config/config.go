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
	Command      string
	Username     string
	Password     string
	PrivPassword string
}

// New returns Config struct
func New(ctx *cli.Context) Config {
	cfg := Config{
		Hostname:     ctx.String("hostname"),
		Port:         ctx.Int("port"),
		Timeout:      ctx.Int("timeout"),
		ExecPlatform: ctx.String("exec-platform"),
		EnableMode:   ctx.Bool("enable-mode"),
		Command:      ctx.String("command"),
		Username:     ctx.String("username"),
		Password:     ctx.String("password"),
		PrivPassword: ctx.String("priv-password"),
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
	if cfg.Username == domain.DefaultUsernameValue {
		return errors.ErrMissingUsername
	}
	if cfg.Password == domain.DefaultPasswordValue {
		return errors.ErrMissingPassword
	}
	if !hasPrivPassword(cfg.EnableMode, cfg.PrivPassword) {
		return errors.ErrMissingPrivPassword
	}
	if !isAbleToExpandTermLength(cfg.EnableMode, cfg.ExecPlatform) {
		return errors.ErrTermLengthIsEnforced
	}
	if cfg.Hostname == "" {
		return errors.ErrMissingHostname
	}
	if cfg.Command == "" {
		return errors.ErrMissingCommand
	}
	if !isValidExpectMode(cfg.ExecPlatform) {
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

func isAbleToExpandTermLength(mode bool, platform string) bool {
	if platform == domain.ASASoftwarePlatformName && !mode {
		return false
	}
	if platform == domain.ASASoftwareHAPlatformName && !mode {
		return false
	}
	return true
}

func isValidExpectMode(platform string) bool {
	for _, p := range domain.CmdPlatforms {
		if platform == p {
			return true
		}
	}
	return false
}

func hasPrivPassword(value bool, password string) bool {
	if value {
		if password == domain.DefaultPrivPasswordValue {
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
		if platform == domain.ScreenOSHAPlatformName {
			return false
		}
	}
	fmt.Println(true)
	return true
}

func isUsableUsername(username string, platform string) bool {
	if username != "" {
		if platform == domain.YamahaOSPlatformName {
			return false
		}
	}
	return true
}
