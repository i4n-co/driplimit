package api

import "github.com/gofiber/fiber/v2"

type rpc struct {
	Action        string
	Namespace     string
	Documentation RPCDocumentation
	Handler       fiber.Handler
}

func (h *rpc) path() string {
	return "/" + h.Namespace + "." + h.Action
}

type rpcs []*rpc

func (s *Server) RegisterRPC(router fiber.Router, rpc *rpc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rpcs = append(s.rpcs, rpc)
	router.Post(rpc.path(), rpc.Handler)
}
