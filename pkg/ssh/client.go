// Package ssh dials a device over SSH and runs one batch on an interactive shell.
package ssh

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
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

// Wraps the known_hosts callback so a mismatch prints the onboarding steps before it fails.
func createFallbackCallback(knownHostsCallback ssh.HostKeyCallback) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		err := knownHostsCallback(hostname, remote, key)
		if err != nil {
			return handleHostKeyVerificationFailure(hostname, err)
		}
		return nil
	}
}

func handleHostKeyVerificationFailure(hostname string, originalErr error) error {
	host := extractHostFromAddress(hostname)

	fmt.Fprintf(os.Stderr, "\n[ERROR] Host key verification failed for %s: %v\n", hostname, originalErr)
	fmt.Fprintln(os.Stderr, "\nTo resolve this issue, you can add the host key to your known_hosts file using one of these methods:")
	fmt.Fprintf(os.Stderr, "\n1. Connect once with ssh and accept the key:\n")

	if isStandardSSHPort(hostname) {
		fmt.Fprintf(os.Stderr, "   ssh %s\n", host)
		fmt.Fprintf(os.Stderr, "\n2. Or, if that reports no matching host key type or key exchange method:\n")
		fmt.Fprintf(os.Stderr, "   ssh -o HostKeyAlgorithms=+ssh-rsa -o KexAlgorithms=+diffie-hellman-group14-sha1 %s\n", host)
	} else {
		port := extractPortFromAddress(hostname)
		fmt.Fprintf(os.Stderr, "   ssh -p %s %s\n", port, host)
		fmt.Fprintf(os.Stderr, "\n2. Or, if that reports no matching host key type or key exchange method:\n")
		fmt.Fprintf(os.Stderr, "   ssh -p %s -o HostKeyAlgorithms=+ssh-rsa -o KexAlgorithms=+diffie-hellman-group14-sha1 %s\n", port, host)
	}

	fmt.Fprintf(os.Stderr, "\n3. Or use the --host-key-path flag to specify a specific host key file\n")
	fmt.Fprintf(os.Stderr, "\n4. Or set TELEE_HOSTKEYPATH environment variable to specify the host key file path\n")
	fmt.Fprintln(os.Stderr, "\nConnection canceled for security reasons.")

	return fmt.Errorf("host key verification failed for %s", hostname)
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

// Fetch dials, opens a shell on a pseudo-terminal, runs the batch and returns the output the
// last prompt match captured.
func (c *SSH) Fetch(batchers *[]x.Batcher, config *ssh.ClientConfig) (string, error) {
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

	out, err := x.Run(shell{pr, stdin}, *batchers, c.timeout)
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
