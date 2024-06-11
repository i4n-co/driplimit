package api

import "github.com/gofiber/fiber/v2"

type rpcs []*rpc
type rpc struct {
	Action        string
	Namespace     string
	Documentation RPCDocumentation
	Handler       fiber.Handler
}

func (h *rpc) path() string {
	return "/" + h.Namespace + "." + h.Action
}

// RegisterRPC registers an RPC endpoint into the server. It
// appends the RPC to the list of RPCs and registers the handler
func (s *Server) RegisterRPC(router fiber.Router, rpc *rpc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rpcs = append(s.rpcs, rpc)
	router.Post(rpc.path(), rpc.Handler)
}
