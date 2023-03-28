package http

import "net/http"

// routes returns defined routes on the muxtiplexer
func (a *WebApp) Routes() *http.ServeMux {

	mux := http.NewServeMux()

	return mux

}
