// Package ssh dials a device over SSH and runs one session on an interactive shell.
package ssh

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"

	x "github.com/umatare5/telee/pkg/expect"
)

const (
	errSSHSpawnFailed = "SSH was failed at spawn(). You can troubleshoot using wireshark.\n"
	errSSHBatchFailed = "SSH was failed at ExpectBatch(). You can troubleshoot using wireshark.\n"
)

// The terminal goexpect requested: its 34 flags minus CS7, which CS8 carries. Every line-discipline
// flag is off, so a server that applies pty modes neither echoes nor rewrites line endings, and the
// 16 control characters and two speeds it also sent mean nothing with ICANON, ISIG and IXON off.
const (
	ptyTerm   = "xterm"
	ptyWidth  = 132
	ptyHeight = 43
)

var ptyModes = ssh.TerminalModes{
	ssh.IGNPAR: 0, ssh.PARMRK: 0, ssh.INPCK: 0, ssh.ISTRIP: 0, ssh.INLCR: 0, ssh.IGNCR: 0, ssh.ICRNL: 0,
	ssh.IUCLC: 0, ssh.IXON: 0, ssh.IXANY: 0, ssh.IXOFF: 0, ssh.IMAXBEL: 0,
	ssh.ISIG: 0, ssh.ICANON: 0, ssh.XCASE: 0, ssh.ECHO: 0, ssh.ECHOE: 0, ssh.ECHOK: 0, ssh.ECHONL: 0,
	ssh.NOFLSH: 0, ssh.TOSTOP: 0, ssh.IEXTEN: 0, ssh.ECHOCTL: 0, ssh.ECHOKE: 0,
	ssh.OPOST: 0, ssh.OLCUC: 0, ssh.ONLCR: 0, ssh.OCRNL: 0, ssh.ONOCR: 0, ssh.ONLRET: 0,
	ssh.CS8: 1, ssh.PARENB: 0, ssh.PARODD: 0,
}

type SSH struct {
	host     string
	port     int
	protocol string
	timeout  time.Duration
}

func New(host string, port int, protocol string, timeout time.Duration) *SSH {
	return &SSH{
		host:     host,
		port:     port,
		protocol: protocol,
		timeout:  timeout,
	}
}

func GenerateClientConfig(username, password, hostKeyPath, hostname string) (*ssh.ClientConfig, error) {
	hostKeyCallback, err := createHostKeyCallback(hostKeyPath, hostname)
	if err != nil {
		return nil, err
	}

	return &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: hostKeyCallback,
	}, nil
}

func createHostKeyCallback(hostKeyPath, hostname string) (ssh.HostKeyCallback, error) {
	if hostKeyPath != "" {
		return createFixedHostKeyCallback(hostKeyPath)
	}
	return createKnownHostsCallback(hostname)
}

func createFixedHostKeyCallback(hostKeyPath string) (ssh.HostKeyCallback, error) {
	publicKeyBytes, err := os.ReadFile(hostKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read host key file: %w", err)
	}

	// A known_hosts line parses here too: its leading host field is taken as the
	// authorized_keys options field, so a .pub file and a scanned line both work.
	publicKey, _, _, _, err := ssh.ParseAuthorizedKey(publicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse host key: %w", err)
	}

	return ssh.FixedHostKey(publicKey), nil
}

func createKnownHostsCallback(hostname string) (ssh.HostKeyCallback, error) {
	knownHostsPath, err := getKnownHostsPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	if _, err := os.Stat(knownHostsPath); err != nil {
		hostOnly := hostname
		if h, _, err := net.SplitHostPort(hostname); err == nil {
			hostOnly = h
		}
		return nil, fmt.Errorf("~/.ssh/known_hosts not found. Please create it by running: ssh %s and accepting the key", hostOnly)
	}

	knownHostsCallback, err := knownhosts.New(knownHostsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load known_hosts file: %w", err)
	}

	return createFallbackCallback(knownHostsCallback), nil
}

func getKnownHostsPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + "/.ssh/known_hosts", nil
}

// Wraps the known_hosts callback so a failure prints what to do before it is returned.
func createFallbackCallback(knownHostsCallback ssh.HostKeyCallback) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		err := knownHostsCallback(hostname, remote, key)
		if err != nil {
			return handleHostKeyVerificationFailure(hostname, key, err)
		}
		return nil
	}
}

func handleHostKeyVerificationFailure(hostname string, key ssh.PublicKey, originalErr error) error {
	fmt.Fprint(os.Stderr, hostKeyFailureMessage(hostname, key, originalErr))
	return fmt.Errorf("host key verification failed for %s", hostname)
}

// hostKeyFailureMessage tells a key that fails its known_hosts record from a host the file does not
// know. The former is what the file exists to catch, so it gets the fingerprint and the recorded line
// and no instruction to accept it, whatever its type; the latter keeps the onboarding steps.
func hostKeyFailureMessage(hostname string, key ssh.PublicKey, originalErr error) string {
	var revoked *knownhosts.RevokedError
	if errors.As(originalErr, &revoked) {
		return fmt.Sprintf("\n[ERROR] Host key for %s is revoked at %s:%d\n\nConnection canceled for security reasons.\n",
			hostname, revoked.Revoked.Filename, revoked.Revoked.Line)
	}

	var keyErr *knownhosts.KeyError
	if errors.As(originalErr, &keyErr) && len(keyErr.Want) > 0 {
		if i := slices.IndexFunc(keyErr.Want, func(k knownhosts.KnownKey) bool { return k.Key.Type() == key.Type() }); i >= 0 {
			return fmt.Sprintf("\n[ERROR] Host key for %s has changed: %s\nThe recorded key is at %s:%d.\n\nConnection canceled for security reasons.\n",
				hostname, ssh.FingerprintSHA256(key), keyErr.Want[i].Filename, keyErr.Want[i].Line)
		}
		// Another type on record cannot be told from another device, so it is refused the same way.
		var recorded []string
		for _, k := range keyErr.Want {
			if !slices.Contains(recorded, k.Key.Type()) {
				recorded = append(recorded, k.Key.Type())
			}
		}
		return fmt.Sprintf("\n[ERROR] Host key for %s does not match known_hosts: it holds %s at %s:%d and the device presented %s %s\n\nConnection canceled for security reasons.\n",
			hostname, strings.Join(recorded, ", "), keyErr.Want[0].Filename, keyErr.Want[0].Line, key.Type(), ssh.FingerprintSHA256(key))
	}

	var b strings.Builder
	fmt.Fprintf(&b, "\n[ERROR] Host key verification failed for %s: %v\n", hostname, originalErr)
	fmt.Fprintln(&b, "\nTo resolve this issue, you can add the host key to your known_hosts file using one of these methods:")
	fmt.Fprintf(&b, "\n1. Connect once with ssh and accept the key:\n")
	fmt.Fprintf(&b, "   %s\n", sshCommand(hostname, ""))
	fmt.Fprintf(&b, "\n2. Or, if that reports no matching host key type or key exchange method:\n")
	fmt.Fprintf(&b, "   %s\n", sshCommand(hostname, "-o HostKeyAlgorithms=+ssh-rsa -o KexAlgorithms=+diffie-hellman-group14-sha1"))
	fmt.Fprintf(&b, "\n3. Or use the --host-key-path flag to specify a specific host key file\n")
	fmt.Fprintf(&b, "\n4. Or set TELEE_HOSTKEYPATH environment variable to specify the host key file path\n")
	fmt.Fprintln(&b, "\nConnection canceled for security reasons.")

	return b.String()
}

// sshCommand is the ssh line that reaches the address, carrying -p for a non-default port.
func sshCommand(address, options string) string {
	parts := []string{"ssh"}
	if !isStandardSSHPort(address) {
		parts = append(parts, "-p", extractPortFromAddress(address))
	}
	if options != "" {
		parts = append(parts, options)
	}
	return strings.Join(append(parts, extractHostFromAddress(address)), " ")
}

func extractHostFromAddress(address string) string {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return address
	}
	return host
}

func extractPortFromAddress(address string) string {
	_, port, err := net.SplitHostPort(address)
	if err != nil || port == "" {
		return "22"
	}
	return port
}

func isStandardSSHPort(address string) bool {
	port := extractPortFromAddress(address)
	return port == "22"
}

// Fetch dials, opens a shell on a pseudo-terminal, runs login then commands and returns the
// transcript of commands.
func (c *SSH) Fetch(login, commands []x.Batcher, config *ssh.ClientConfig) (string, error) {
	client, conn, err := c.dial(config)
	if err != nil {
		fmt.Fprint(os.Stderr, errSSHSpawnFailed)
		return "", err
	}
	defer client.Close() //nolint:errcheck

	session, err := client.NewSession()
	if err != nil {
		fmt.Fprint(os.Stderr, errSSHSpawnFailed)
		return "", err
	}
	defer session.Close() //nolint:errcheck

	pr, pw := io.Pipe()
	defer pw.Close() //nolint:errcheck
	stdin, err := openShell(session, pw)
	if err != nil {
		fmt.Fprint(os.Stderr, errSSHSpawnFailed)
		return "", err
	}
	// The deadline dial set is lifted once the shell is up, because the multiplexer reads
	// this socket for the rest of the session.
	if err := conn.SetDeadline(time.Time{}); err != nil {
		fmt.Fprint(os.Stderr, errSSHSpawnFailed)
		return "", err
	}
	// Ends the batch's input when the shell ends, so a device closing the session fails the
	// step at once instead of running out the timeout.
	go func() {
		session.Wait() //nolint:errcheck,gosec
		pw.Close()     //nolint:errcheck,gosec
	}()

	out, err := x.Run(shell{pr, stdin}, login, commands, c.timeout)
	if err != nil {
		fmt.Fprint(os.Stderr, errSSHBatchFailed)
		return "", err
	}
	return out, nil
}

// shell is the batch's view of a session: its joined output and its input.
type shell struct {
	io.Reader
	io.Writer
}

// openShell starts an interactive shell whose output streams both feed out, as goexpect read
// both, and returns its input.
func openShell(session *ssh.Session, out io.Writer) (io.Writer, error) {
	session.Stdout, session.Stderr = out, out
	stdin, err := session.StdinPipe()
	if err != nil {
		return nil, err
	}
	if err := session.RequestPty(ptyTerm, ptyHeight, ptyWidth, ptyModes); err != nil {
		return nil, err
	}
	if err := session.Shell(); err != nil {
		return nil, err
	}
	return stdin, nil
}

// dial connects under the timeout and leaves one more timeout as the socket deadline, which the
// handshake, the authentication and the shell request share and Fetch lifts afterwards.
func (c *SSH) dial(config *ssh.ClientConfig) (*ssh.Client, net.Conn, error) {
	addr := net.JoinHostPort(c.host, strconv.Itoa(c.port))
	conn, err := net.DialTimeout(c.protocol, addr, c.timeout)
	if err != nil {
		return nil, nil, err
	}
	if err := conn.SetDeadline(time.Now().Add(c.timeout)); err != nil {
		conn.Close() //nolint:errcheck,gosec
		return nil, nil, err
	}
	// NewClientConn closes the socket itself when the handshake or the authentication fails.
	sshConn, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
	if err != nil {
		return nil, nil, err
	}
	return ssh.NewClient(sshConn, chans, reqs), conn, nil
}
