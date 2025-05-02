package routes

import "net/http"

func ServerFilesApi(mux *http.ServeMux) {
	mux.Handle("/api/v1/uploads/", http.StripPrefix("/api/v1/uploads/", http.FileServer(http.Dir("./uploads"))))
}
