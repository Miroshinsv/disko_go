package main

import (
	"disko/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)
	//Roles
	r.HandleFunc("/roles/get/all/", handlers.GetAllRoles).Methods(http.MethodGet)
	r.HandleFunc("/roles/get/{id}/", handlers.GetRoleById).Methods(http.MethodGet)
	r.HandleFunc("/roles/update/{id}/", handlers.UpdateRoleById).Methods(http.MethodPost)
	r.HandleFunc("/roles/disband/{id}/", handlers.DisbandRoleById).Methods(http.MethodPost)
	r.HandleFunc("/roles/add/", handlers.AddRole).Methods(http.MethodPost)
	//Directions
	r.HandleFunc("/direction/get/all/", handlers.GetAllDirections).Methods(http.MethodGet)
	r.HandleFunc("/direction/get/{id}/", handlers.GetDirectionById).Methods(http.MethodGet)
	r.HandleFunc("/direction/update/{id}/", handlers.UpdateDirectionById).Methods(http.MethodPost)
	r.HandleFunc("/direction/disband/{id}/", handlers.DisbandDirectionById).Methods(http.MethodPost)
	r.HandleFunc("/direction/add/", handlers.AddDirection).Methods(http.MethodPost)
	//Users
	r.HandleFunc("/user/get/all/", handlers.GetAllUsers).Methods(http.MethodGet)
	r.HandleFunc("/user/get/{id}/", handlers.GetUserById).Methods(http.MethodGet)
	r.HandleFunc("/user/update/{id}/", handlers.UpdateUserById).Methods(http.MethodPost)
	r.HandleFunc("/user/disband/{id}/", handlers.DisbandUserById).Methods(http.MethodPost)
	r.HandleFunc("/user/add/", handlers.AddUser).Methods(http.MethodPost)
	//Events
	r.HandleFunc("/events/get/all/", handlers.GetAllEvents).Methods(http.MethodGet)
	r.HandleFunc("/events/get/{id}/", handlers.GetEventById).Methods(http.MethodGet)
	r.HandleFunc("/events/add/", handlers.AddEvent).Methods(http.MethodPost)
	r.HandleFunc("/events/update/{id}/", handlers.GetAllEvents).Methods(http.MethodPost)
	r.HandleFunc("/events/activate/{id}/", handlers.ActivateEventById).Methods(http.MethodPost)
	//Schools
	r.HandleFunc("/schools/get/all/", handlers.GetAllSchools).Methods(http.MethodGet)
	r.HandleFunc("/schools/get/{id}/", handlers.GetSchoolById).Methods(http.MethodGet)
	r.HandleFunc("/schools/delete/{id}/", handlers.DeleteSchoolById).Methods(http.MethodDelete)
	r.HandleFunc("/schools/add/", handlers.AddSchool).Methods(http.MethodPost)
	r.HandleFunc("/schools/update/{id}/", handlers.UpdateSchoolById).Methods(http.MethodPost)

	http.Handle("/", r)
	http.ListenAndServe(":8081", nil)
}
