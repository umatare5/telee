// Package telnet dials a device over Telnet and runs one session on the connection.
package telnet

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	x "github.com/umatare5/telee/pkg/expect"
)

const (
	errTelnetSpawnFailed = "TelnetClient was failed at spawn(). You can troubleshoot using wireshark.\n"
	errTelnetBatchFailed = "TelnetClient was failed at ExpectBatch(). You can troubleshoot using wireshark.\n"
)

// Commands and options from RFC 854, 857, 858 and 1073.
const (
	cmdSE   byte = 240
	cmdSB   byte = 250
	cmdWill byte = 251
	cmdWont byte = 252
	cmdDo   byte = 253
	cmdDont byte = 254
	cmdIAC  byte = 255

	optEcho byte = 1
	optSGA  byte = 3
	optNAWS byte = 31

	bufSize = 32 << 10
)

type Telnet struct {
	host     string
	port     int
	protocol string
	timeout  time.Duration
}

func New(host string, port int, protocol string, timeout time.Duration) *Telnet {
	return &Telnet{
		host:     host,
		port:     port,
		protocol: protocol,
		timeout:  timeout,
	}
}

// Fetch dials, runs login then commands and returns the transcript of commands.
func (t *Telnet) Fetch(login, commands []x.Batcher) (string, error) {
	conn, err := net.DialTimeout(t.protocol, net.JoinHostPort(t.host, strconv.Itoa(t.port)), t.timeout)
	if err != nil {
		fmt.Fprint(os.Stderr, errTelnetSpawnFailed)
		return "", err
	}
	defer conn.Close() //nolint:errcheck

	out, err := x.Run(newConn(conn), login, commands, t.timeout)
	if err != nil {
		fmt.Fprint(os.Stderr, errTelnetBatchFailed)
		return "", err
	}
	return out, nil
}

// conn strips the protocol out of the stream and answers its option negotiation.
type conn struct {
	net.Conn
	r    *bufio.Reader
	echo bool
	sga  bool
}

func newConn(c net.Conn) *conn {
	return &conn{Conn: c, r: bufio.NewReaderSize(c, bufSize)}
}

// Read returns application bytes only. It returns as soon as the buffered input is used up,
// so a prompt is delivered without waiting for the segment after it.
func (c *conn) Read(p []byte) (int, error) {
	n := 0
	for n < len(p) && (n == 0 || c.r.Buffered() > 0) {
		b, data, err := c.next()
		if err != nil {
			return n, err
		}
		if data {
			p[n] = b
			n++
		}
	}
	return n, nil
}

// next consumes one data byte or one protocol sequence, answering the latter.
func (c *conn) next() (b byte, data bool, err error) {
	if b, err = c.r.ReadByte(); err != nil {
		return 0, false, err
	}
	if b != cmdIAC {
		return b, true, nil
	}
	cmd, err := c.r.ReadByte()
	if err != nil {
		return 0, false, err
	}
	switch cmd {
	case cmdIAC:
		return cmdIAC, true, nil
	case cmdWill, cmdWont, cmdDo, cmdDont:
		opt, err := c.r.ReadByte()
		if err != nil {
			return 0, false, err
		}
		return 0, false, c.negotiate(cmd, opt)
	case cmdSB:
		return 0, false, c.skipSubnegotiation()
	default:
		// NOP, DM, BRK, IP, AO, AYT, EC, EL and GA carry no argument.
		return 0, false, nil
	}
}

// negotiate takes ECHO and SGA, answers NAWS, and refuses every other option. A refusal is
// never answered, per RFC 854.
func (c *conn) negotiate(cmd, opt byte) error {
	switch {
	case opt == optEcho || opt == optSGA:
		return c.toggle(cmd, opt)
	case cmd == cmdDo && opt == optNAWS:
		if err := c.reply(cmdWill, opt); err != nil {
			return err
		}
		// 65535 columns by 65535 rows, each 0xFF doubled. Paging on a server that sizes it
		// from the window stays off, as it did under ziutek/telnet.
		_, err := c.Conn.Write([]byte{
			cmdIAC, cmdSB, optNAWS,
			cmdIAC, cmdIAC, cmdIAC, cmdIAC, cmdIAC, cmdIAC, cmdIAC, cmdIAC,
			cmdIAC, cmdSE,
		})
		return err
	case cmd == cmdDo:
		return c.reply(cmdWont, opt)
	case cmd == cmdWill:
		return c.reply(cmdDont, opt)
	}
	return nil
}

// toggle agrees to ECHO and SGA in either direction, answering only a change of state so the
// exchange cannot loop.
func (c *conn) toggle(cmd, opt byte) error {
	on := &c.echo
	if opt == optSGA {
		on = &c.sga
	}
	want := cmd == cmdDo || cmd == cmdWill
	if *on == want {
		return nil
	}
	*on = want

	var ack byte
	switch cmd {
	case cmdDo:
		ack = cmdWill
	case cmdDont:
		ack = cmdWont
	case cmdWill:
		ack = cmdDo
	default:
		ack = cmdDont
	}
	return c.reply(ack, opt)
}

func (c *conn) reply(cmd, opt byte) error {
	_, err := c.Conn.Write([]byte{cmdIAC, cmd, opt})
	return err
}

func (c *conn) skipSubnegotiation() error {
	for {
		b, err := c.r.ReadByte()
		if err != nil {
			return err
		}
		if b != cmdIAC {
			continue
		}
		if b, err = c.r.ReadByte(); err != nil {
			return err
		}
		if b == cmdSE {
			return nil
		}
	}
}

// Write doubles the IAC byte, the only escape the protocol needs. A line ending goes as given.
func (c *conn) Write(p []byte) (int, error) {
	if _, err := c.Conn.Write(bytes.ReplaceAll(p, []byte{cmdIAC}, []byte{cmdIAC, cmdIAC})); err != nil {
		return 0, err
	}
	return len(p), nil
}
