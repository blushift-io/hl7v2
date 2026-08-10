package tcp

// HandlerFunc defines a function signature for processing connection events.
type HandlerFunc func(*Context) error

// Handler processes connection events including connect, message, error, and close callbacks.
type Handler interface {
	OnConnect(*Context) error
	OnMessage(*Context) error
	OnError(*Context) error
	OnClose(*Context) error
}

type handler struct {
	onConnect []HandlerFunc
	onMessage []HandlerFunc
	onError   []HandlerFunc
	onClose   []HandlerFunc
}

// HandlerOption defines a function signature for configuring a Handler.
type HandlerOption func(*handler)

// OnConnect returns a HandlerOption that registers a callback for connection establishment.
func OnConnect(fn HandlerFunc) HandlerOption {
	return func(h *handler) {
		h.onConnect = append(h.onConnect, fn)
	}
}

// OnMessage returns a HandlerOption that registers a callback for incoming messages.
func OnMessage(fn HandlerFunc) HandlerOption {
	return func(h *handler) {
		h.onMessage = append(h.onMessage, fn)
	}
}

// OnError returns a HandlerOption that registers a callback for connection errors.
func OnError(fn HandlerFunc) HandlerOption {
	return func(h *handler) {
		h.onError = append(h.onError, fn)
	}
}

// OnClose returns a HandlerOption that registers a callback for connection closure.
func OnClose(fn HandlerFunc) HandlerOption {
	return func(h *handler) {
		h.onClose = append(h.onClose, fn)
	}
}

// NewHandler constructs a Handler initialized with the provided options.
func NewHandler(opts ...HandlerOption) Handler {
	h := &handler{}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *handler) OnConnect(ctx *Context) error {
	for _, fn := range h.onConnect {
		if err := fn(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (h *handler) OnMessage(ctx *Context) error {
	for _, fn := range h.onMessage {
		if err := fn(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (h *handler) OnError(ctx *Context) error {
	for _, fn := range h.onError {
		if err := fn(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (h *handler) OnClose(ctx *Context) error {
	for _, fn := range h.onClose {
		if err := fn(ctx); err != nil {
			return err
		}
	}

	return nil
}
