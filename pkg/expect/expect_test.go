package expect_test

import (
	"bufio"
	"errors"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	x "github.com/umatare5/telee/pkg/expect"
)

const timeout = 300 * time.Millisecond

// dial returns the client end of a pipe whose server end is driven by device.
func dial(t *testing.T, device func(c net.Conn)) net.Conn {
	t.Helper()
	client, server := net.Pipe()
	t.Cleanup(func() {
		client.Close()
		server.Close()
	})
	go device(server)
	return client
}

func TestRunLogin(t *testing.T) {
	t.Parallel()
	conn := dial(t, func(c net.Conn) {
		r := bufio.NewReader(c)
		c.Write([]byte("Username: "))
		if line, _ := r.ReadString('\n'); line != "admin\n" {
			return
		}
		c.Write([]byte("Password: "))
		if line, _ := r.ReadString('\n'); line != "secret\n" {
			return
		}
		c.Write([]byte("\r\nsw01>"))
	})
	batch := []x.Batcher{
		&x.BExp{R: "Username:"},
		&x.BSnd{S: "admin\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: "secret\n"},
		&x.BExp{R: "sw01>"},
	}
	out, err := x.Run(conn, batch, timeout)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if out != "\r\nsw01>" {
		t.Fatalf("Run() = %q; want %q", out, "\r\nsw01>")
	}
}

func TestRunSplitPromptAndSilenceReset(t *testing.T) {
	t.Parallel()
	conn := dial(t, func(c net.Conn) {
		// Six writes spread over twice the timeout, each inside it, then a prompt in two halves.
		for range 6 {
			time.Sleep(timeout / 3)
			c.Write([]byte("line\n"))
		}
		c.Write([]byte("sw0"))
		c.Write([]byte("1>"))
	})
	out, err := x.Run(conn, []x.Batcher{&x.BExp{R: "sw01>"}}, timeout)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if want := strings.Repeat("line\n", 6) + "sw01>"; out != want {
		t.Fatalf("Run() = %q; want %q", out, want)
	}
}

func TestRunDiscardsBytesAfterMatch(t *testing.T) {
	t.Parallel()
	conn := dial(t, func(c net.Conn) {
		c.Write([]byte("sw01>tail"))
		bufio.NewReader(c).ReadString('\n')
	})
	batch := []x.Batcher{&x.BExp{R: "sw01>"}, &x.BSnd{S: "next\n"}, &x.BExp{R: "tail"}}
	_, err := x.Run(conn, batch, timeout)
	var timeoutErr x.TimeoutError
	if !errors.As(err, &timeoutErr) {
		t.Fatalf("Run() error = %v; want TimeoutError", err)
	}
}

func TestRunSendTimeout(t *testing.T) {
	t.Parallel()
	conn := dial(t, func(net.Conn) {})
	_, err := x.Run(conn, []x.Batcher{&x.BSnd{S: "show version\n"}}, timeout)
	var timeoutErr x.TimeoutError
	if !errors.As(err, &timeoutErr) {
		t.Fatalf("Run() error = %v; want TimeoutError from a peer that never reads", err)
	}
}

func TestRunConnectionClosed(t *testing.T) {
	t.Parallel()
	conn := dial(t, func(c net.Conn) {
		c.Write([]byte("Username: "))
		c.Close()
	})
	_, err := x.Run(conn, []x.Batcher{&x.BExp{R: "Password:"}}, timeout)
	if !errors.Is(err, io.EOF) {
		t.Fatalf("Run() error = %v; want io.EOF", err)
	}
}

func TestRunInvalidPattern(t *testing.T) {
	t.Parallel()
	conn := dial(t, func(net.Conn) {})
	if _, err := x.Run(conn, []x.Batcher{&x.BExp{R: "("}}, timeout); err == nil {
		t.Fatal("Run() error = nil; want a regexp error")
	}
}

func TestTimeoutErrorMessage(t *testing.T) {
	t.Parallel()
	got := x.TimeoutError(5 * time.Second).Error()
	if want := "expect: timer expired after 5 seconds"; got != want {
		t.Fatalf("Error() = %q; want %q", got, want)
	}
}
