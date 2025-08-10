package ssh

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	x "github.com/google/goexpect"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
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
func GenerateClientConfig(username string, password string, hostKeyPath string, hostname string) (*ssh.ClientConfig, error) {
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

// createHostKeyCallback creates appropriate HostKeyCallback based on hostKeyPath
func createHostKeyCallback(hostKeyPath string, hostname string) (ssh.HostKeyCallback, error) {
	if hostKeyPath != "" {
		return createFixedHostKeyCallback(hostKeyPath)
	}
	return createKnownHostsCallback(hostname)
}

// createFixedHostKeyCallback creates HostKeyCallback from specific host key file
func createFixedHostKeyCallback(hostKeyPath string) (ssh.HostKeyCallback, error) {
	publicKeyBytes, err := os.ReadFile(hostKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read host key file: %w", err)
	}

	publicKey, err := ssh.ParsePublicKey(publicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse host key: %w", err)
	}

	return ssh.FixedHostKey(publicKey), nil
}

// createKnownHostsCallback creates HostKeyCallback using known_hosts file
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
		return nil, fmt.Errorf("~/.ssh/known_hosts not found. Please create it by running: ssh-keyscan %s >> ~/.ssh/known_hosts", hostOnly)
	}

	knownHostsCallback, err := knownhosts.New(knownHostsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load known_hosts file: %w", err)
	}

	return createFallbackCallback(knownHostsCallback), nil
}

// getKnownHostsPath returns the path to the known_hosts file
func getKnownHostsPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + "/.ssh/known_hosts", nil
}

// createFallbackCallback creates a callback that tries known_hosts first, then provides guidance
func createFallbackCallback(knownHostsCallback ssh.HostKeyCallback) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		err := knownHostsCallback(hostname, remote, key)
		if err != nil {
			return handleHostKeyVerificationFailure(hostname, err)
		}
		return nil
	}
}

// handleHostKeyVerificationFailure handles host key verification failure and provides user guidance
func handleHostKeyVerificationFailure(hostname string, originalErr error) error {
	// Extract hostname without port
	host := extractHostFromAddress(hostname)

	fmt.Fprintf(os.Stderr, "\n[ERROR] Host key verification failed for %s: %v\n", hostname, originalErr)
	fmt.Fprintln(os.Stderr, "\nTo resolve this issue, you can add the host key to your known_hosts file using one of these methods:")
	fmt.Fprintf(os.Stderr, "\n1. Run the following command to add the host key:\n")

	// Check if it's a standard SSH port or custom port
	if isStandardSSHPort(hostname) {
		fmt.Fprintf(os.Stderr, "   ssh-keyscan %s >> ~/.ssh/known_hosts\n", host)
		fmt.Fprintf(os.Stderr, "\n2. Or connect manually first with ssh to accept the host key:\n")
		fmt.Fprintf(os.Stderr, "   ssh %s\n", host)
		fmt.Fprintf(os.Stderr, "\n   For older Cisco IOS devices, you may need additional SSH options:\n")
		fmt.Fprintf(os.Stderr, "   ssh -o HostKeyAlgorithms=+ssh-rsa -o KexAlgorithms=+diffie-hellman-group14-sha1 %s\n", host)
	} else {
		port := extractPortFromAddress(hostname)
		fmt.Fprintf(os.Stderr, "   ssh-keyscan -p %s %s >> ~/.ssh/known_hosts\n", port, host)
		fmt.Fprintf(os.Stderr, "\n2. Or connect manually first with ssh to accept the host key:\n")
		fmt.Fprintf(os.Stderr, "   ssh -p %s %s\n", port, host)
		fmt.Fprintf(os.Stderr, "\n   For older Cisco IOS devices, you may need additional SSH options:\n")
		fmt.Fprintf(os.Stderr, "   ssh -p %s -o HostKeyAlgorithms=+ssh-rsa -o KexAlgorithms=+diffie-hellman-group14-sha1 %s\n", port, host)
	}

	fmt.Fprintf(os.Stderr, "\n3. Or use the --host-key-path flag to specify a specific host key file\n")
	fmt.Fprintf(os.Stderr, "\n4. Or set TELEE_HOSTKEYPATH environment variable to specify the host key file path\n")
	fmt.Fprintln(os.Stderr, "\nConnection cancelled for security reasons.")

	return fmt.Errorf("host key verification failed for %s", hostname)
}

// extractHostFromAddress extracts hostname from "hostname:port" format
func extractHostFromAddress(address string) string {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		// If address does not contain a port, return as is
		return address
	}
	return host
}

// extractPortFromAddress extracts port from "hostname:port" format
func extractPortFromAddress(address string) string {
	_, port, err := net.SplitHostPort(address)
	if err != nil || port == "" {
		return "22"
	}
	return port
}

// isStandardSSHPort checks if the address uses standard SSH port (22)
func isStandardSSHPort(address string) bool {
	port := extractPortFromAddress(address)
	return port == "22"
}

// Fetch starts the expect process.
func (c *SSH) Fetch(batchers *[]x.Batcher, config *ssh.ClientConfig) (string, error) {
	conn, err := c.dial(config)
	if err != nil {
		fmt.Println(errSSHSpawnFailed)
		return "", err
	}
	defer conn.Close() //nolint: errcheck

	expecter, _, err := x.SpawnSSH(conn, c.timeout)
	if err != nil {
		fmt.Print(errSSHSpawnFailed)
		return "", err
	}
	defer expecter.Close() //nolint: errcheck

	stdout, err := expecter.ExpectBatch(*batchers, c.timeout)
	if err != nil {
		fmt.Print(errSSHBatchFailed)
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
