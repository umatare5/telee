package repository

import (
	"telee/internal/config"
	"telee/pkg/telnet"
	"time"

	x "github.com/google/goexpect"
)

const (
	protocol string = "tcp"
)

// ServerRepository struct
type ServerRepository struct {
	Config *config.Config
}

// Fetch returns stdout from telnet session
func (r *ServerRepository) Fetch(x *[]x.Batcher) (string, error) {
	client := telnet.New(
		r.Config.Hostname, r.Config.Port, protocol, time.Duration(r.Config.Timeout)*time.Second,
	)
	data, err := client.Fetch(x)
	if err != nil {
		return "", err
	}
	return data, nil
}
