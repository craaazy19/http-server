package routing

import (
	"fmt"
	"net/http"
	"time"
)

func (r *Router) Test(w http.ResponseWriter, req *http.Request) {
	r.shutdownMutex.RLock()
	defer r.shutdownMutex.RUnlock()

	time.Sleep(clientWaitTimeout - 100*time.Millisecond)
	_, _ = fmt.Fprintf(w, "test")
}
