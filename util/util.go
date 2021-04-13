package util

import "net/http"

// ContentType sets the ContentType header to type
func ContentType(h http.HandlerFunc, ctype string) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", ctype)
		h(w, r)
	}
	return http.HandlerFunc(fn)
}
