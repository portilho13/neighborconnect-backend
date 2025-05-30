package routes

import "net/http"

func ServerFilesApi(mux *http.ServeMux) {
	mux.Handle("/api/v1/uploads/listing/", http.StripPrefix("/api/v1/uploads/listing/", http.FileServer(http.Dir("./uploads/listing"))))
	mux.Handle("/api/v1/uploads/events/", http.StripPrefix("/api/v1/uploads/events/", http.FileServer(http.Dir("./uploads/events"))))
	mux.Handle("/api/v1/uploads/users/", http.StripPrefix("/api/v1/uploads/users/", http.FileServer(http.Dir("./uploads/users"))))
}
