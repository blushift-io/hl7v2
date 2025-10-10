package tcp

type HandlerFunc func(*Context) error

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

type HandlerOption func(*handler)

func OnConnect(fn HandlerFunc) HandlerOption {
	return func(h *handler) {
		h.onConnect = append(h.onConnect, fn)
	}
}

func OnMessage(fn HandlerFunc) HandlerOption {
	return func(h *handler) {
		h.onMessage = append(h.onMessage, fn)
	}
}

func OnError(fn HandlerFunc) HandlerOption {
	return func(h *handler) {
		h.onError = append(h.onError, fn)
	}
}

func OnClose(fn HandlerFunc) HandlerOption {
	return func(h *handler) {
		h.onClose = append(h.onClose, fn)
	}
}

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
