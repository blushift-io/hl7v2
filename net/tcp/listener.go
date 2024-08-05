package tcp

import (
	"fmt"
	"net"
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

	if err := c.SetKeepAlivePeriod(l.d); err != nil {
		return nil, err
	}

	return c, nil
}

type listener struct {
	l     []net.Listener
	close chan struct{}
	conns chan accepter
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
	return l.l[0].Addr()
}

func (l *listener) Close() error {
	close(l.close)

	var errs []error
	for i := range l.l {
		if err := l.l[i].Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	var errMsg string
	for i := range errs {
		errMsg += fmt.Sprintf("%s\n", errs[i].Error())
	}

	return fmt.Errorf("error closing listener: %s", errMsg)
}

func (l *listener) Accept() (net.Conn, error) {
	select {
	case a, ok := <-l.conns:
		if ok {
			return a.conn, a.err
		}

		return nil, fmt.Errorf("inner listener channel closed")
	case <-l.close:
		return nil, fmt.Errorf("listener closed")
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
		a := accepter{
			conn: conn,
			err:  err,
		}

		select {
		case l.conns <- a:
		case <-l.close:
			if a.err == nil {
				a.conn.Close()
			}

			return
		}
	}
}
