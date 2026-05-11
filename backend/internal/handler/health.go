package handler

import "net/http"

func Health(w http.ResponseWriter, r *http.Request) {
	writeData(w, map[string]string{"status": "ok"})
}
