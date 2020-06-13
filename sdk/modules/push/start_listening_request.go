package push

type StartListeningRequest struct {
	Type       string
	Parameters string
	Handler    Handler
}
