package config

import (
	"github.com/metacubex/mihomo/component/updater"
	"github.com/metacubex/mihomo/hub/route"
	"github.com/metacubex/mihomo/log"

	"github.com/metacubex/chi"
	"github.com/metacubex/chi/render"
	"github.com/metacubex/http"
)

func init() {
	// In embed mode mihomo skips POST /upgrade (it downloads a new
	// executable and re-execs, impossible for a library inside an apk) and
	// both geo update routes (POST /upgrade/geo, POST /configs/geo), so
	// dashboards surface bare 404s. Replace them here: a clear error for
	// core upgrade and a real in-process implementation for geo updates.
	// These static routes take precedence over the mounted subrouters for
	// the exact method+path only.
	route.Register(func(r chi.Router) {
		r.Post("/upgrade", upgradeCore)
		r.Post("/upgrade/geo", updateGeoDatabases)
		r.Post("/configs/geo", updateGeoDatabases)
	})
}

func upgradeCore(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotImplemented)
	render.JSON(w, r, render.M{"message": "core is embedded in the app, please upgrade the CMFA apk instead"})
}

func updateGeoDatabases(w http.ResponseWriter, r *http.Request) {
	err := updater.UpdateGeoDatabases()
	if err != nil {
		log.Errorln("[GEO] update GEO databases failed: %v", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{"message": err.Error()})
		return
	}

	render.NoContent(w, r)
}
