package ssh

import (
	"fmt"
	"strconv"
	"time"

	expect "github.com/google/goexpect"
	x "github.com/google/goexpect"
	"golang.org/x/crypto/ssh"
)

const (
	errSSHSpawnFailed = "SSH was failed at spawn(). You can troubleshoot using wireshark.\n"
	errSSHBatchFailed = "SSH was failed at ExpectBatch(). You can troubleshoot using wireshark.\n"
)

// SSH struct
type SSH struct {
	host     string
	port     int
	protocol string
	timeout  time.Duration
}

// New returns SSH struct
func New(host string, port int, protocol string, timeout time.Duration) *SSH {
	return &SSH{
		host:     host,
		port:     port,
		protocol: protocol,
		timeout:  timeout,
	}
}

// GenerateClientConfig returns client config
func GenerateClientConfig(username string, password string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // nolint: gosec
	}
}

// Fetch starts the expect process.
func (c *SSH) Fetch(x *[]x.Batcher, config *ssh.ClientConfig) (string, error) {
	conn, err := c.dial(config)
	if err != nil {
		fmt.Println(errSSHSpawnFailed)
		return "", err
	}
	defer conn.Close() // nolint: errcheck

	expecter, _, err := expect.SpawnSSH(conn, c.timeout)
	if err != nil {
		fmt.Println(errSSHSpawnFailed)
		return "", err
	}
	defer expecter.Close() // nolint: errcheck

	stdout, err := expecter.ExpectBatch(*x, c.timeout)
	if err != nil {
		fmt.Println(errSSHBatchFailed)
		return "", err
	}

	return stdout[len(stdout)-1].Output, nil
}

func (c *SSH) dial(config *ssh.ClientConfig) (*ssh.Client, error) {
	conn, err := ssh.Dial(c.protocol, c.host+":"+strconv.Itoa(c.port), config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
