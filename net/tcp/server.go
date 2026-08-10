package tcp

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/mllp"
)

// Server manages an MLLP TCP server that accepts client connections and handles messages.
type Server struct {
	opts      *ServerOptions
	h         Handler
	l         net.Listener
	conns     map[net.Conn]Conn
	connL     sync.RWMutex
	close     chan struct{}
	closeOnce sync.Once
}

// NewServer creates a new Server with the specified Handler and options.
func NewServer(h Handler, opts ...ServerOption) (*Server, error) {
	if h == nil {
		return nil, fmt.Errorf("tcp: handler cannot be nil")
	}

	o := NewServerOptions(opts...)

	return &Server{
		opts:  o,
		h:     h,
		conns: make(map[net.Conn]Conn),
		close: make(chan struct{}),
	}, nil
}

// Start opens listeners on configured addresses and begins serving incoming connections.
func (s *Server) Start() error {
	log.Printf("starting server on %s", s.opts.Addresses)

	var ls []net.Listener
	for _, addr := range s.opts.Addresses {
		l, err := net.Listen("tcp", addr)
		if err != nil {
			return err
		}

		ls = append(ls, l)
	}

	l, err := newListener(ls...)
	if err != nil {
		return err
	}

	return s.Serve(l)
}

// Serve accepts connections on the provided net.Listener and processes them.
func (s *Server) Serve(l net.Listener) error {
	s.connL.Lock()
	s.l = l
	if s.conns == nil {
		s.conns = make(map[net.Conn]Conn)
	}
	if s.close == nil {
		s.close = make(chan struct{})
	}
	s.connL.Unlock()

	defer func() {
		if s.l != nil {
			s.l.Close()
		}
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			select {
			case <-s.close:
				return nil
			default:
			}

			if errors.Is(err, net.ErrClosed) {
				return nil
			}

			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				time.Sleep(time.Millisecond * 5)
				continue
			}

			return err
		}

		nc, err := newConn(conn, DefaultConnOptions())
		if err != nil {
			conn.Close()
			return err
		}

		ctx := NewContext(nc)

		if s.h != nil {
			if err := s.h.OnConnect(ctx); err != nil {
				nc.Close()
				continue
			}
		}

		s.connL.Lock()
		s.conns[conn] = nc
		s.connL.Unlock()

		go s.handleConn(conn, nc)
	}
}

// Shutdown gracefully stops the server and closes all active client connections.
func (s *Server) Shutdown() error {
	s.closeOnce.Do(func() {
		if s.close != nil {
			close(s.close)
		}
	})

	s.connL.RLock()
	l := s.l
	s.connL.RUnlock()

	if l != nil {
		if err := l.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			return err
		}
	}

	return s.releaseConns()
}

func (s *Server) handleConn(rawConn net.Conn, conn Conn) {
	defer func() {
		s.connL.Lock()
		delete(s.conns, rawConn)
		s.connL.Unlock()

		if conn != nil {
			conn.Close()
		}

		if s.h != nil {
			ctx := NewContext(conn)
			s.h.OnClose(ctx)
		}
	}()

	scanner := mllp.NewBufferedScanner(conn, 1024*1024)
	ctx := NewContext(conn)

	for {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				ctx.Apply(SetError(err))

				if s.h != nil {
					if herr := s.h.OnError(ctx); herr != nil {
						break
					}
				}

				continue
			}
			break
		}

		b := scanner.Bytes()
		if len(b) == 0 {
			continue
		}

		msg, err := hl7v2.ParseRaw(b)
		if err != nil {
			ctx.Apply(SetError(err))

			if s.h != nil {
				if herr := s.h.OnError(ctx); herr != nil {
					break
				}
			}

			continue
		}

		ctx.Apply(SetData(msg))
		if s.h != nil {
			if err := s.h.OnMessage(ctx); err != nil {
				ctx.Apply(SetError(err))
				if herr := s.h.OnError(ctx); herr != nil {
					break
				}

				continue
			}
		}
	}
}

func (s *Server) releaseConns() error {
	var errs []error

	s.connL.Lock()
	defer s.connL.Unlock()

	for rawConn, c := range s.conns {
		if cerr := c.Close(); cerr != nil {
			errs = append(errs, cerr)
		}
		delete(s.conns, rawConn)
	}

	return errors.Join(errs...)
}
