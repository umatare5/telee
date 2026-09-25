// Package expect runs one device session as a list of steps over a transport.
package expect

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strings"
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

// TimeoutError is returned when a BExp receives nothing for its whole window, or a BSnd is
// not written within it.
type TimeoutError time.Duration

func (t TimeoutError) Error() string {
	return fmt.Sprintf("expect: timer expired after %d seconds", time.Duration(t)/time.Second)
}

// Run executes login, then commands, and returns the transcript of commands: the prompt the
// last BExp of login matched, then everything each BExp of commands captured. A BExp captures
// everything received since the previous match, and discards it for the next step.
func Run(rw io.ReadWriter, login, commands []Batcher, timeout time.Duration) (string, error) {
	s := &session{chunks: make(chan []byte), errc: make(chan error, 1), done: make(chan struct{}), timeout: timeout}
	defer close(s.done)
	go s.read(rw)

	var out strings.Builder
	for i, b := range slices.Concat(login, commands) {
		switch b := b.(type) {
		case *BExp:
			re, err := regexp.Compile(b.R)
			if err != nil {
				return "", err
			}
			buf, at, err := s.expect(re)
			if err != nil {
				return "", err
			}
			if i < len(login) {
				out.Reset()
				buf = buf[at:]
			}
			out.Write(buf)
		case *BSnd:
			if err := s.send(rw, b.S); err != nil {
				return "", err
			}
		}
	}
	return out.String(), nil
}

// Commands sends each command with suffix and waits for prompt after it.
func Commands(prompt, suffix string, cmds []string) []Batcher {
	var batch []Batcher
	for _, c := range cmds {
		batch = append(batch, &BSnd{S: c + suffix}, &BExp{R: prompt})
	}
	return batch
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

// expect appends chunks until re matches, giving up after timeout of silence. It returns the
// buffer and the offset the match starts at.
func (s *session) expect(re *regexp.Regexp) (buf []byte, at int, err error) {
	timer := time.NewTimer(s.timeout)
	defer timer.Stop()

	for {
		var chunk []byte
		select {
		case chunk = <-s.chunks:
		case err = <-s.errc:
			return nil, 0, fmt.Errorf("expect: connection closed before a match: %w", err)
		case <-timer.C:
			// A chunk that arrived with the expiry still counts, as it did under goexpect.
			select {
			case chunk = <-s.chunks:
			default:
				return nil, 0, TimeoutError(s.timeout)
			}
		}
		buf = append(buf, chunk...)
		off := max(0, len(buf)-len(chunk)-window)
		if loc := re.FindIndex(buf[off:]); loc != nil {
			return buf, off + loc[0], nil
		}
		timer.Reset(s.timeout)
	}
}
