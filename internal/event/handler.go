package event

type HandleType int
type HandleFunction func() error

type Handler struct {
	HandleType
	HandledFunc map[HandleType]HandleFunction
}

func (h *Handler) HandleFunc(t HandleType) error {
	for e, f := range h.HandledFunc {
		if t == e {
			return f()
		}
	}

	return nil
}

func (h *Handler) NewHandler(f HandleFunction, t HandleType) {
	if h.HandledFunc == nil {
		h.HandledFunc = make(map[HandleType]HandleFunction)
	}

	h.HandledFunc[t] = f
}
