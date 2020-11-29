package routing

import "sync"

type Router struct {
	shutdownMutex sync.RWMutex
}

func (r *Router) Shutdown() {
	r.shutdownMutex.Lock()
}
