package telnet

import (
	"fmt"
	"strconv"
	"time"

	x "github.com/google/goexpect"
	"github.com/ziutek/telnet"
)

const (
	errTelnetSpawnFailed  = "TelnetClient was failed at spawn(). You can troubleshoot using wireshark.\n"
	errTelnetBatchFailed  = "TelnetClient was failed at ExpectBatch(). You can troubleshoot using wireshark."
	hintTelnetBatchFailed = "[Hint] Invalid username or password. mismatch hostnames your set and configured one.\n"
)

// Telnet struct
type Telnet struct {
	host     string
	port     int
	protocol string
	timeout  time.Duration
}

// New returns Telnet struct
func New(host string, port int, protocol string, timeout time.Duration) *Telnet {
	return &Telnet{
		host:     host,
		port:     port,
		protocol: protocol,
		timeout:  timeout,
	}
}

// Fetch starts the expect process.
func (t *Telnet) Fetch(x *[]x.Batcher) (string, error) {
	conn, _, err := t.spawn()
	if err != nil {
		fmt.Println(errTelnetSpawnFailed)
		return "", err
	}
	defer conn.Close() // nolint: errcheck

	stdout, err := conn.ExpectBatch(*x, t.timeout)
	if err != nil {
		fmt.Println(errTelnetBatchFailed)
		fmt.Println(hintTelnetBatchFailed)
		return "", err
	}

	return stdout[len(stdout)-1].Output, nil
}

func (t *Telnet) spawn() (x.Expecter, <-chan error, error) {
	conn, err := telnet.Dial(t.protocol, t.host+":"+strconv.Itoa(t.port))
	if err != nil {
		return nil, nil, err
	}

	ch := make(chan error)

	return x.SpawnGeneric(
		&x.GenOptions{
			In:    conn,
			Out:   conn,
			Check: func() bool { return true },
			Wait:  func() error { return <-ch },
			Close: func() error {
				close(ch)
				return conn.Close()
			},
		}, t.timeout)
}
