package config

import (
	"cfa/native/tunnel"

	"github.com/metacubex/mihomo/hub/route"
	"github.com/metacubex/mihomo/log"

	"github.com/metacubex/chi"
	"github.com/metacubex/chi/render"
	"github.com/metacubex/http"
)

func init() {
	// mihomo never mounts its own /restart in embed mode (see
	// hub/route/server.go): its implementation re-execs the executable and
	// exits, which would kill the whole app process. Register an
	// embed-friendly implementation through the external router hook
	// instead; it is mounted inside the authenticated route group.
	route.Register(func(r chi.Router) {
		r.Mount("/restart", restartRouter())
	})
}

func restartRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/", restart)
	return r
}

func restart(w http.ResponseWriter, r *http.Request) {
	if LoadedPath() == "" {
		render.Status(r, http.StatusServiceUnavailable)
		render.JSON(w, r, render.M{"message": "no profile loaded"})
		return
	}

	render.JSON(w, r, render.M{"status": "ok"})
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	// The core is embedded in the app process and its lifecycle belongs to
	// the app/VpnService, so "restart" means rebuilding the core state in
	// place: drop all connections, reset statistics and re-apply the
	// current profile from disk. The TUN device is attached by TunService
	// independently of hub.ApplyConfig and survives the reload.
	go func() {
		log.Infoln("[APP] restarting core requested by RESTful API")

		tunnel.CloseAllConnections()
		tunnel.ResetStatistic()

		if err := Reload(); err != nil {
			log.Errorln("Restart core: %s", err.Error())
		}
	}()
}
