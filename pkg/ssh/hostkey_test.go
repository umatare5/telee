package ssh_test

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"strings"
	"testing"

	cryptossh "golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"

	"github.com/umatare5/telee/pkg/ssh"
)

func ed25519Key(t *testing.T) cryptossh.PublicKey {
	t.Helper()
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	key, err := cryptossh.NewPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}
	return key
}

func ecdsaKey(t *testing.T) cryptossh.PublicKey {
	t.Helper()
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	key, err := cryptossh.NewPublicKey(&priv.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	return key
}

func TestHostKeyFailureMessage(t *testing.T) {
	t.Parallel()
	presented := ed25519Key(t)
	recorded := knownhosts.KnownKey{Key: ed25519Key(t), Filename: "known_hosts", Line: 3}
	other := knownhosts.KnownKey{Key: ecdsaKey(t), Filename: "known_hosts", Line: 7}
	tests := []struct {
		name     string
		hostname string
		err      error
		want     []string
		unwanted string
	}{
		{
			name: "unknown host", hostname: "sw01:22", err: &knownhosts.KeyError{},
			want: []string{"1. Connect once with ssh", "   ssh sw01\n", "   ssh -o HostKeyAlgorithms=+ssh-rsa -o KexAlgorithms=+diffie-hellman-group14-sha1 sw01\n"},
		},
		{
			name: "unknown host on another port", hostname: "sw01:2222", err: &knownhosts.KeyError{},
			want: []string{"   ssh -p 2222 sw01\n"},
		},
		{
			name: "changed key", hostname: "sw01:22", err: &knownhosts.KeyError{Want: []knownhosts.KnownKey{other, recorded}},
			want: []string{"has changed: " + cryptossh.FingerprintSHA256(presented), "known_hosts:3"}, unwanted: "accept",
		},
		{
			name: "other key type recorded", hostname: "sw01:22", err: &knownhosts.KeyError{Want: []knownhosts.KnownKey{other}},
			want: []string{"does not match known_hosts: it holds ecdsa-sha2-nistp256 at known_hosts:7 and the device presented ssh-ed25519 " + cryptossh.FingerprintSHA256(presented)}, unwanted: "accept",
		},
		{
			name: "revoked key", hostname: "sw01:22", err: &knownhosts.RevokedError{Revoked: recorded},
			want: []string{"is revoked at known_hosts:3"}, unwanted: "accept",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ssh.HostKeyFailureMessage(tt.hostname, presented, tt.err)
			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Errorf("hostKeyFailureMessage() = %q; want it to contain %q", got, want)
				}
			}
			if tt.unwanted != "" && strings.Contains(got, tt.unwanted) {
				t.Errorf("hostKeyFailureMessage() = %q; must not contain %q", got, tt.unwanted)
			}
		})
	}
}
