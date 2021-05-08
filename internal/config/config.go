package config

import (
	"log"
	"telee/internal/domain"
	"telee/pkg/errors"

	"github.com/jinzhu/configor"
	"github.com/urfave/cli/v2"
)

// Config struct
type Config struct {
	Hostname     string
	Port         int
	Timeout      int
	Platform     string
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
		Platform:     ctx.String("platform"),
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
	if !isAbleToExpandTermLength(cfg.EnableMode, cfg.Platform) {
		return errors.ErrTermLengthIsEnforced
	}
	if cfg.Hostname == "" {
		return errors.ErrMissingHostname
	}
	if cfg.Command == "" {
		return errors.ErrMissingCommand
	}
	if !isValidExpectMode(cfg.Platform) {
		return errors.ErrInvalidPlatform
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
