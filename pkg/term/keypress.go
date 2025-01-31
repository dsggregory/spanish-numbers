package term

import (
	"context"
	"golang.org/x/sys/unix"
	"log"
	"os"
)

// KeypressReader the data for the handler
type KeypressReader struct {
	kpchan        chan byte
	origTermState *unix.Termios
	cancel        context.CancelFunc
}

// NewKeypressHandler start listening for keypresses and returns them through the channel h.KeyEvent().
// Caller MUST call h.Reset() before exiting the program.
func NewKeypressReader(cancel context.CancelFunc) (*KeypressReader, error) {
	h := KeypressReader{cancel: cancel}
	if err := h.handleKeypress(); err != nil {
		return nil, err
	}
	return &h, nil
}

// KeyEvent caller uses this to select on the keypress event channel
func (h *KeypressReader) KeyEvent() chan byte {
	return h.kpchan
}

// Reset reset terminal
func (h *KeypressReader) Reset() error {
	return unix.IoctlSetTermios(int(os.Stdin.Fd()), unix.TIOCSETA, h.origTermState)
}

// handleKeypress to run a goroutine that returns keypress data to the caller through a channel.
func (h *KeypressReader) handleKeypress() error {
	ostate, err := makeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	h.origTermState = ostate

	h.kpchan = make(chan byte)

	go func() {
		var b = make([]byte, 1)
		for {
			_, err := os.Stdin.Read(b) // read one byte from raw terminal
			if err != nil {
				log.Println(err.Error())
				h.cancel()
				return
			}
			h.kpchan <- b[0]
			if b[0] == 'q' {
				h.cancel()
				return
			}
		}
	}()
	return err
}

// our own b/c we don't want to muck with output or signal reception
func makeRaw(fd int) (*unix.Termios, error) {
	termios, err := unix.IoctlGetTermios(fd, unix.TIOCGETA)
	if err != nil {
		return nil, err
	}

	oldState := *termios

	// This attempts to replicate the behaviour documented for cfmakeraw in
	// the termios(3) manpage.
	termios.Iflag &^= unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON
	//termios.Oflag &^= unix.OPOST
	termios.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | /*unix.ISIG |*/ unix.IEXTEN
	termios.Cflag &^= unix.CSIZE | unix.PARENB
	termios.Cflag |= unix.CS8
	termios.Cc[unix.VMIN] = 1
	termios.Cc[unix.VTIME] = 0
	if err := unix.IoctlSetTermios(fd, unix.TIOCSETA, termios); err != nil {
		return nil, err
	}

	return &oldState, nil
}
