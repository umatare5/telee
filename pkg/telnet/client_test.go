package telnet_test

import (
	"bytes"
	"io"
	"net"
	"testing"
	"time"

	"github.com/umatare5/telee/pkg/telnet"
)

const (
	iac  = 255
	dont = 254
	do   = 253
	wont = 252
	will = 251
	sb   = 250
	se   = 240
	nop  = 241
)

// exchange feeds wire to the client, then returns what Read produced and what the client sent back.
func exchange(t *testing.T, wire []byte, replyLen int) (data, replies []byte) {
	t.Helper()
	client, server := net.Pipe()
	t.Cleanup(func() {
		client.Close()
		server.Close()
	})
	server.SetDeadline(time.Now().Add(2 * time.Second))

	got := make(chan []byte, 1)
	go func() {
		server.Write(wire)
		replies := make([]byte, replyLen)
		io.ReadFull(server, replies)
		got <- replies
	}()

	buf := make([]byte, 64)
	n, err := telnet.NewConn(client).Read(buf)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	return buf[:n], <-got
}

func TestReadNegotiation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wire    []byte
		data    []byte
		replies []byte
	}{
		{
			name:    "cisco greeting",
			wire:    []byte{iac, will, 1, iac, will, 3, iac, do, 24, iac, do, 31, 'l', 'o', 'g', 'i', 'n', ':'},
			data:    []byte("login:"),
			replies: []byte{iac, do, 1, iac, do, 3, iac, wont, 24, iac, will, 31, iac, sb, 31, iac, iac, iac, iac, iac, iac, iac, iac, iac, se},
		},
		{
			name:    "repeated will echo answered once",
			wire:    []byte{iac, will, 1, iac, will, 1, 'x'},
			data:    []byte("x"),
			replies: []byte{iac, do, 1},
		},
		{
			name:    "do echo taken and dont echo released",
			wire:    []byte{iac, do, 1, iac, dont, 1, 'x'},
			data:    []byte("x"),
			replies: []byte{iac, will, 1, iac, wont, 1},
		},
		{
			name: "refused option not answered on wont",
			wire: []byte{iac, wont, 24, iac, dont, 24, 'x'},
			data: []byte("x"),
		},
		{
			name: "escaped iac and subnegotiation and nop",
			wire: []byte{'a', iac, iac, iac, sb, 24, 0, iac, iac, 'b', iac, se, iac, nop, 'c'},
			data: []byte{'a', iac, 'c'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data, replies := exchange(t, tt.wire, len(tt.replies))
			if !bytes.Equal(data, tt.data) {
				t.Errorf("Read() = %q; want %q", data, tt.data)
			}
			if !bytes.Equal(replies, tt.replies) {
				t.Errorf("replies = %v; want %v", replies, tt.replies)
			}
		})
	}
}

func TestWriteEscapesIAC(t *testing.T) {
	t.Parallel()
	client, server := net.Pipe()
	t.Cleanup(func() {
		client.Close()
		server.Close()
	})
	got := make(chan []byte, 1)
	go func() {
		buf := make([]byte, 4)
		io.ReadFull(server, buf)
		got <- buf
	}()
	n, err := telnet.NewConn(client).Write([]byte{'a', iac, '\n'})
	if err != nil || n != 3 {
		t.Fatalf("Write() = %d, %v; want 3, nil", n, err)
	}
	if wire, want := <-got, []byte{'a', iac, iac, '\n'}; !bytes.Equal(wire, want) {
		t.Fatalf("wire = %v; want %v", wire, want)
	}
}
