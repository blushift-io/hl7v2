package tcp

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

type keepAliveListener struct {
	d time.Duration
	*net.TCPListener
}

func (l keepAliveListener) Accept() (net.Conn, error) {
	c, err := l.AcceptTCP()
	if err != nil {
		return nil, err
	}

	if err := c.SetKeepAlive(true); err != nil {
		return nil, err
	}

	if l.d > 0 {
		if err := c.SetKeepAlivePeriod(l.d); err != nil {
			return nil, err
		}
	}

	return c, nil
}

type listener struct {
	l         []net.Listener
	close     chan struct{}
	closeOnce sync.Once
	conns     chan accepter
}

type accepter struct {
	conn net.Conn
	err  error
}

func newListener(l ...net.Listener) (*listener, error) {
	if len(l) == 0 {
		return nil, fmt.Errorf("no listeners")
	}

	li := &listener{
		l:     l,
		close: make(chan struct{}),
		conns: make(chan accepter),
	}

	li.start()

	return li, nil
}

func (l *listener) Addr() net.Addr {
	if l == nil || len(l.l) == 0 {
		return nil
	}
	return l.l[0].Addr()
}

func (l *listener) Close() error {
	if l == nil {
		return nil
	}

	l.closeOnce.Do(func() {
		close(l.close)
	})

	var errs []error
	for i := range l.l {
		if err := l.l[i].Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (l *listener) Accept() (net.Conn, error) {
	select {
	case a, ok := <-l.conns:
		if ok {
			return a.conn, a.err
		}
		return nil, net.ErrClosed
	case <-l.close:
		return nil, net.ErrClosed
	}
}

func (l *listener) start() {
	for i := range l.l {
		go l.run(l.l[i])
	}
}

func (l *listener) run(nl net.Listener) {
	for {
		conn, err := nl.Accept()
		if err != nil {
			select {
			case <-l.close:
				return
			default:
				a := accepter{conn: conn, err: err}
				select {
				case l.conns <- a:
				case <-l.close:
					return
				}
				if errors.Is(err, net.ErrClosed) {
					return
				}
				continue
			}
		}

		a := accepter{
			conn: conn,
			err:  err,
		}

		select {
		case l.conns <- a:
		case <-l.close:
			if a.conn != nil {
				a.conn.Close()
			}
			return
		}
	}
}
