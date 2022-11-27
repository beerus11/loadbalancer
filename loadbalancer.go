package loadbalancer

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Server struct {
	URL          *url.URL
	IsDead       bool
	lock         *sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

func (s *Server) SetHealth(status bool) {
	s.lock.Lock()
	s.IsDead = status
	s.lock.Unlock()
}

func (s *Server) GetHealth() bool {
	s.lock.Lock()
	heath := s.IsDead
	s.lock.Unlock()
	return heath
}

type ServerPool struct {
	Servers        []*Server
	ServerCount    int64
	LastServerUsed int64
}

func (sp *ServerPool) GetNextAvailableServer() int64 {
	return ((sp.LastServerUsed + 1) % sp.ServerCount)
}

func (sp *ServerPool) UpdateLastServerUsed(idx int64) {
	sp.LastServerUsed = idx
	fmt.Println("LastServer Used is :: ", idx)
}

func (sp *ServerPool) RegisterServer(server *Server) {
	sp.Servers = append(sp.Servers, server)
	sp.ServerCount++
}

func GetNewServerPool() (pool *ServerPool) {
	return new(ServerPool)
}

func GetNewServer(path string) *Server {
	server := new(Server)
	server.URL, _ = url.Parse(path)
	server.IsDead = false
	server.ReverseProxy = httputil.NewSingleHostReverseProxy(server.URL)
	server.lock = new(sync.RWMutex)
	return server
}
