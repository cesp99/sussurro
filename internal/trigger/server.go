package trigger

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"path/filepath"
)

// Server listens for trigger events via UNIX socket
type Server struct {
	socket     string
	listener   net.Listener
	log        *slog.Logger
	done       chan struct{}
	onKeyDown  func()
	onKeyUp    func()
	isRecording bool
}

// NewServer creates a new trigger server
func NewServer(log *slog.Logger) (*Server, error) {
	// Create socket in user's runtime directory
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = "/tmp"
	}

	socketPath := filepath.Join(runtimeDir, "sussurro.sock")

	// Remove existing socket if present
	os.Remove(socketPath)

	return &Server{
		socket: socketPath,
		log:    log,
		done:   make(chan struct{}),
	}, nil
}

// Start starts listening for trigger events
func (s *Server) Start(onKeyDown, onKeyUp func()) error {
	s.onKeyDown = onKeyDown
	s.onKeyUp = onKeyUp

	listener, err := net.Listen("unix", s.socket)
	if err != nil {
		return fmt.Errorf("failed to create socket: %w", err)
	}
	s.listener = listener

	// Set permissions so user can access it
	os.Chmod(s.socket, 0600)

	s.log.Debug("Trigger server started", "socket", s.socket)

	go s.listen()

	return nil
}

// Stop stops the server
func (s *Server) Stop() {
	close(s.done)
	if s.listener != nil {
		s.listener.Close()
	}
	os.Remove(s.socket)
}

func (s *Server) listen() {
	for {
		select {
		case <-s.done:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				select {
				case <-s.done:
					return
				default:
					s.log.Error("Failed to accept connection", "error", err)
					continue
				}
			}

			go s.handleConnection(conn)
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}

	cmd := string(buf[:n])
	s.log.Debug("Received trigger command", "cmd", cmd)

	// Toggle recording state
	if !s.isRecording {
		s.log.Info("Recording started - press hotkey again when done speaking")
		s.isRecording = true
		if s.onKeyDown != nil {
			s.onKeyDown()
		}
		conn.Write([]byte("RECORDING\n"))
	} else {
		s.log.Info("Recording stopped - processing...")
		s.isRecording = false
		if s.onKeyUp != nil {
			s.onKeyUp()
		}
		conn.Write([]byte("STOPPED\n"))
	}

	// Try to send desktop notification if notify-send is available
	if !s.isRecording {
		// Just finished recording
		exec.Command("notify-send", "-t", "2000", "Sussurro", "Processing your speech...").Start()
	}
}

// GetSocketPath returns the socket path for external triggering
func (s *Server) GetSocketPath() string {
	return s.socket
}
