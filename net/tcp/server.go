package tcp

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/blushift-io/hl7v2"
	"github.com/blushift-io/hl7v2/mllp"
)

type Server struct {
	opts  *ServerOptions
	h     Handler
	l     net.Listener
	conns map[net.Conn]Conn
	connL sync.RWMutex
	close chan struct{}
}

func NewServer(h Handler, opts ...ServerOption) (*Server, error) {
	o := NewServerOptions(opts...)

	return &Server{
		opts: o,
		h:    h,
	}, nil
}

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

func (s *Server) Serve(l net.Listener) error {
	s.l = l

	s.conns = make(map[net.Conn]Conn)
	s.close = make(chan struct{})

	defer func() {
		s.l.Close()
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			select {
			case <-s.close:
				return nil
			default:
			}

			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				time.Sleep(time.Millisecond * 5)
				continue
			}

			return err
		}

		nc, err := newConn(conn, DefaultConnOptions())
		if err != nil {
			return err
		}

		ctx := NewContext(nc)

		if err := s.h.OnConnect(ctx); err != nil {
			return err
		}

		s.connL.Lock()
		s.conns[conn] = nc
		s.connL.Unlock()

		go s.handleConn(nc)
	}
}

func (s *Server) Shutdown() error {
	if s.l != nil {
		if err := s.l.Close(); err != nil {
			return err
		}
	}

	select {
	case s.close <- struct{}{}:
	default:
	}

	return s.releaseConns()
}

func (s *Server) handleConn(conn Conn) {
	scanner := mllp.NewBufferedScanner(conn, 2048)

	ctx := NewContext(conn)

	defer s.h.OnClose(ctx)

	for {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				ctx.Apply(SetError(err))

				if herr := s.h.OnError(ctx); herr != nil {
					break
				}

				continue
			}
		}

		b := scanner.Bytes()
		if len(b) == 0 {
			continue
		}

		msg, err := hl7v2.NewRawMessageFromBytes(b)
		if err != nil {
			ctx.Apply(SetError(err))

			if herr := s.h.OnError(ctx); herr != nil {
				break
			}

			continue
		}

		ctx.Apply(SetData(msg))
		if err := s.h.OnMessage(ctx); err != nil {
			ctx.Apply(SetError(err))
			if herr := s.h.OnError(ctx); herr != nil {
				break
			}

			continue
		}
	}

}

func (s *Server) releaseConns() error {
	var errs []error

	s.connL.Lock()
	defer s.connL.Unlock()

	for _, c := range s.conns {
		if cerr := c.Close(); cerr != nil {
			errs = append(errs, cerr)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	var errMsg string
	for _, err := range errs {
		errMsg += err.Error() + "\n"
	}

	return fmt.Errorf("errors releasing connections: %s", errMsg)
}
