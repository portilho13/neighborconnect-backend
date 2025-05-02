package routes

import "net/http"

func ServerFilesApi(mux *http.ServeMux) {
	mux.Handle("/api/v1/uploads/listing/", http.StripPrefix("/api/v1/uploads/listing/", http.FileServer(http.Dir("./uploads/listing"))))
	mux.Handle("/api/v1/uploads/event/", http.StripPrefix("/api/v1/uploads/event/", http.FileServer(http.Dir("./uploads/event"))))
	mux.Handle("/api/v1/uploads/category/", http.StripPrefix("/api/v1/uploads/category/", http.FileServer(http.Dir("./uploads/category"))))
}
