package main

import (
	"disko/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/roles/get/all/", handlers.GetAllRoles).Methods(http.MethodGet)
	r.HandleFunc("/roles/get/{id}/", handlers.GetRoleById).Methods(http.MethodGet)
	r.HandleFunc("/roles/update/{id}/", handlers.UpdateRoleById).Methods(http.MethodPost)
	r.HandleFunc("/roles/disband/{id}/", handlers.DisbandRoleById).Methods(http.MethodPost)
	r.HandleFunc("/roles/add/", handlers.AddRole).Methods(http.MethodPost)
	r.HandleFunc("/events/get/all/", handlers.GetAllEvents).Methods(http.MethodGet)
	http.Handle("/", r)
	http.ListenAndServe(":8081", nil)
}
