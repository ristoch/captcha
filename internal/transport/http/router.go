package http

import (
	"net/http"

	wsTransport "captcha-service/internal/transport/websocket"
)

type Router struct {
	demoHandler *DemoHandler
	wsHandler   *wsTransport.DemoWebSocketHandler
}

func NewRouter(demoHandler *DemoHandler, wsHandler *wsTransport.DemoWebSocketHandler) *Router {
	return &Router{
		demoHandler: demoHandler,
		wsHandler:   wsHandler,
	}
}

func (r *Router) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/backgrounds/", http.StripPrefix("/backgrounds/", http.FileServer(http.Dir("./backgrounds/"))))
	mux.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates/"))))

	mux.HandleFunc("/ws", r.wsHandler.HandleWebSocket)

	mux.HandleFunc("/health", r.demoHandler.HandleHealth)
	mux.HandleFunc("/demo", r.demoHandler.HandleDemo)
	mux.HandleFunc("/performance", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Performance test not implemented", http.StatusNotImplemented)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/demo", http.StatusFound)
	})

	return mux
}
