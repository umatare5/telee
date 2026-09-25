package ssh_test

import (
	"bufio"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"io"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	cryptossh "golang.org/x/crypto/ssh"

	x "github.com/umatare5/telee/pkg/expect"
	"github.com/umatare5/telee/pkg/ssh"
)

const timeout = 10 * time.Second

// serve runs an SSH server on the loopback whose shell is driven by handler, and returns its port
// and a client configuration pinned to its host key.
func serve(t *testing.T, handler func(ch cryptossh.Channel)) (int, *cryptossh.ClientConfig) {
	t.Helper()
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	signer, err := cryptossh.NewSignerFromKey(priv)
	if err != nil {
		t.Fatal(err)
	}
	hostKey, err := cryptossh.NewPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "host.pub")
	if err := os.WriteFile(path, cryptossh.MarshalAuthorizedKey(hostKey), 0o600); err != nil {
		t.Fatal(err)
	}
	config, err := ssh.GenerateClientConfig("operator", "secret", path, "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}

	serverConfig := &cryptossh.ServerConfig{
		PasswordCallback: func(_ cryptossh.ConnMetadata, password []byte) (*cryptossh.Permissions, error) {
			if string(password) != "secret" {
				return nil, errors.New("denied")
			}
			return nil, nil
		},
	}
	serverConfig.AddHostKey(signer)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { listener.Close() })
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		serverConn, channels, requests, err := cryptossh.NewServerConn(conn, serverConfig)
		if err != nil {
			return
		}
		defer serverConn.Close()
		go cryptossh.DiscardRequests(requests)
		for newChannel := range channels {
			ch, channelRequests, err := newChannel.Accept()
			if err != nil {
				return
			}
			go func() {
				for r := range channelRequests {
					r.Reply(r.Type == "pty-req" || r.Type == "shell", nil)
					if r.Type == "shell" {
						handler(ch)
						ch.Close()
					}
				}
			}()
		}
	}()
	return listener.Addr().(*net.TCPAddr).Port, config
}

func TestFetchDialogue(t *testing.T) {
	t.Parallel()
	port, config := serve(t, func(ch cryptossh.Channel) {
		r := bufio.NewReader(ch)
		ch.Write([]byte("Username: "))
		if line, _ := r.ReadString('\n'); line != "operator\n" {
			return
		}
		ch.Write([]byte("\r\nsw01>"))
	})
	batch := []x.Batcher{&x.BExp{R: "Username:"}, &x.BSnd{S: "operator\n"}, &x.BExp{R: "sw01>"}}
	out, err := ssh.New("127.0.0.1", port, "tcp", timeout).Fetch(batch, nil, config)
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if out != "sw01>" {
		t.Fatalf("Fetch() = %q; want %q", out, "sw01>")
	}
}

func TestFetchRemoteClose(t *testing.T) {
	t.Parallel()
	port, config := serve(t, func(ch cryptossh.Channel) {
		ch.Write([]byte("line\r\n"))
	})
	batch := []x.Batcher{&x.BExp{R: "never"}}
	start := time.Now()
	_, err := ssh.New("127.0.0.1", port, "tcp", timeout).Fetch(batch, nil, config)
	if !errors.Is(err, io.EOF) {
		t.Fatalf("Fetch() error = %v; want io.EOF", err)
	}
	if elapsed := time.Since(start); elapsed >= time.Second {
		t.Fatalf("Fetch() took %v; want the close to end the step at once", elapsed)
	}
}

func TestGenerateClientConfig(t *testing.T) {
	t.Parallel()
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	key, err := cryptossh.NewPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "host.pub")
	if err := os.WriteFile(path, cryptossh.MarshalAuthorizedKey(key), 0o600); err != nil {
		t.Fatal(err)
	}
	config, err := ssh.GenerateClientConfig("operator", "secret", path, "sw01")
	if err != nil {
		t.Fatalf("GenerateClientConfig() error = %v", err)
	}
	if len(config.Auth) != 1 || config.HostKeyCallback == nil {
		t.Fatalf("GenerateClientConfig() = %d auth methods, callback %v; want 1 and a callback", len(config.Auth), config.HostKeyCallback)
	}
}
