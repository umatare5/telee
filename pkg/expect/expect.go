// Package expect runs one device session as a list of steps over a transport.
package expect

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"time"
)

const (
	readSize = 32 << 10
	// A new read is scanned together with this much of what preceded it, so a prompt split
	// across two reads is still found without rescanning the whole output. Every BExp is a
	// fixed-length literal; an anchored or unbounded pattern could match differently here.
	window = 4 << 10
)

// Batcher is one step of a session, either a BExp or a BSnd.
type Batcher interface {
	step()
}

// BExp waits until the output received since the previous match matches R.
type BExp struct {
	R string
}

// BSnd writes S to the device unchanged.
type BSnd struct {
	S string
}

func (*BExp) step() {}
func (*BSnd) step() {}

// TimeoutError is returned when a step makes no progress for its whole window.
type TimeoutError time.Duration

func (t TimeoutError) Error() string {
	return fmt.Sprintf("expect: timer expired after %d seconds", time.Duration(t)/time.Second)
}

// Run executes the batch in order and returns the output the last BExp captured. A BExp
// returns everything received since the previous match, and discards it for the next step.
func Run(rw io.ReadWriter, batch []Batcher, timeout time.Duration) (string, error) {
	s := &session{chunks: make(chan []byte), errc: make(chan error, 1), done: make(chan struct{}), timeout: timeout}
	defer close(s.done)
	go s.read(rw)

	var out string
	for _, b := range batch {
		switch b := b.(type) {
		case *BExp:
			re, err := regexp.Compile(b.R)
			if err != nil {
				return "", err
			}
			if out, err = s.expect(re); err != nil {
				return "", err
			}
		case *BSnd:
			if err := s.send(rw, b.S); err != nil {
				return "", err
			}
		}
	}
	return out, nil
}

type session struct {
	chunks  chan []byte
	errc    chan error
	done    chan struct{}
	timeout time.Duration
}

// read feeds the transport into chunks until it fails or Run returns. A chunk is handed over
// before its error, so the bytes that arrived with a close are matched first.
func (s *session) read(r io.Reader) {
	buf := make([]byte, readSize)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			select {
			case s.chunks <- bytes.Clone(buf[:n]):
			case <-s.done:
				return
			}
		}
		if err != nil {
			select {
			case s.errc <- err:
			case <-s.done:
			}
			return
		}
	}
}

// send bounds the write, because a peer that grants no window (RFC 4254 section 5.2) blocks it
// with nothing else to wake it; closing the transport releases the goroutine.
func (s *session) send(w io.Writer, str string) error {
	errc := make(chan error, 1)
	go func() {
		_, err := io.WriteString(w, str)
		errc <- err
	}()
	timer := time.NewTimer(s.timeout)
	defer timer.Stop()
	select {
	case err := <-errc:
		return err
	case <-timer.C:
		return TimeoutError(s.timeout)
	}
}

// expect appends chunks until re matches, giving up after timeout of silence.
func (s *session) expect(re *regexp.Regexp) (string, error) {
	timer := time.NewTimer(s.timeout)
	defer timer.Stop()

	var buf []byte
	for {
		var chunk []byte
		select {
		case chunk = <-s.chunks:
		case err := <-s.errc:
			return "", fmt.Errorf("expect: connection closed before a match: %w", err)
		case <-timer.C:
			// A chunk that arrived with the expiry still counts, as it did under goexpect.
			select {
			case chunk = <-s.chunks:
			default:
				return "", TimeoutError(s.timeout)
			}
		}
		buf = append(buf, chunk...)
		if re.Match(buf[max(0, len(buf)-len(chunk)-window):]) {
			return string(buf), nil
		}
		timer.Reset(s.timeout)
	}
}
