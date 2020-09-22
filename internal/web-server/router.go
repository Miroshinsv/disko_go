package web_server

import (
	auth_service "github.com/Miroshinsv/disko_go/internal/auth-service"
	directionService "github.com/Miroshinsv/disko_go/internal/direction-service"
	eventService "github.com/Miroshinsv/disko_go/internal/event-service"
	poll_service "github.com/Miroshinsv/disko_go/internal/poll-service"
	roleService "github.com/Miroshinsv/disko_go/internal/role-service"
	schedule_service "github.com/Miroshinsv/disko_go/internal/schedule-service"
	schoolService "github.com/Miroshinsv/disko_go/internal/school-service"
	userService "github.com/Miroshinsv/disko_go/internal/user-service"
	"github.com/Miroshinsv/disko_go/internal/web-server/middleware"

	"github.com/gorilla/mux"
	"net/http"
)

var WebRouter *mux.Router

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func RegisterHandlers() {
	WebRouter = mux.NewRouter()
	WebRouter.StrictSlash(true)

	WebRouter.Use(jsonMiddleware)
	WebRouter.Use(middleware.AuthMiddleware)

	//Directions
	hDirection := directionService.MustNewHandlerDirection()
	WebRouter.HandleFunc("/direction/get/all/", hDirection.GetAllDirections).Methods(http.MethodGet)
	WebRouter.HandleFunc("/direction/get/{id}/", hDirection.GetDirectionById).Methods(http.MethodGet)
	WebRouter.HandleFunc("/direction/update/{id}/", hDirection.UpdateDirectionById).Methods(http.MethodPost)
	WebRouter.HandleFunc("/direction/disband/{id}/", hDirection.DisbandDirectionById).Methods(http.MethodPost)
	WebRouter.HandleFunc("/direction/add/", hDirection.AddDirection).Methods(http.MethodPost)

	//Events
	hEvents := eventService.MustNewHandlerEvent()
	WebRouter.HandleFunc("/events/get/all/", hEvents.GetAllEvents).Methods(http.MethodGet)
	WebRouter.HandleFunc("/events/get/{id}/", hEvents.GetEventById).Methods(http.MethodGet)
	WebRouter.HandleFunc("/events/add/", hEvents.AddEvent).Methods(http.MethodPost)
	WebRouter.HandleFunc("/events/update/{id}/", hEvents.GetAllEvents).Methods(http.MethodPost)
	WebRouter.HandleFunc("/events/activate/{id}/", hEvents.ActivateEventById).Methods(http.MethodPost)

	//Roles
	hRoles := roleService.MustNewHandlerRole()
	WebRouter.HandleFunc("/roles/get/all/", hRoles.GetAllRoles).Methods(http.MethodGet)
	WebRouter.HandleFunc("/roles/get/{id}/", hRoles.GetRoleById).Methods(http.MethodGet)
	WebRouter.HandleFunc("/roles/update/{id}/", hRoles.UpdateRoleById).Methods(http.MethodPost)
	WebRouter.HandleFunc("/roles/disband/{id}/", hRoles.DisbandRoleById).Methods(http.MethodPost)
	WebRouter.HandleFunc("/roles/add/", hRoles.AddRole).Methods(http.MethodPost)

	//School
	hSchool := schoolService.MustNewHandlerSchool()
	WebRouter.HandleFunc("/schools/get/all/", hSchool.GetAllSchools).Methods(http.MethodGet)
	WebRouter.HandleFunc("/schools/get/{id}/", hSchool.GetSchoolById).Methods(http.MethodGet)
	WebRouter.HandleFunc("/schools/delete/{id}/", hSchool.DeleteSchoolById).Methods(http.MethodDelete)
	WebRouter.HandleFunc("/schools/add/", hSchool.AddSchool).Methods(http.MethodPost)
	WebRouter.HandleFunc("/schools/update/{id}/", hSchool.UpdateSchoolById).Methods(http.MethodPost)

	//Users
	hUsers := userService.MustNewHandlerUser()
	WebRouter.HandleFunc("/user/get/all/", hUsers.GetAllUsers).Methods(http.MethodGet)
	WebRouter.HandleFunc("/user/get/{id}/", hUsers.GetUserById).Methods(http.MethodGet)
	WebRouter.HandleFunc("/user/update/{id}/", hUsers.UpdateUserById).Methods(http.MethodPost)
	WebRouter.HandleFunc("/user/disband/{id}/", hUsers.DisbandUserById).Methods(http.MethodPost)
	WebRouter.HandleFunc("/user/add/", hUsers.AddUser).Methods(http.MethodPost)

	//Schedule
	hSchedule := schedule_service.MustNewHandlerSchedule()
	WebRouter.HandleFunc("/schedule/today/", hSchedule.LoadEventsForToday).Methods(http.MethodGet)
	WebRouter.HandleFunc("/schedule/period/", hSchedule.LoadEventsForPeriod).Methods(http.MethodGet)

	//Auth
	hAuth := auth_service.MustNewHandlerAuth()
	WebRouter.HandleFunc("/auth/register/", hAuth.Register).Methods(http.MethodPost)
	WebRouter.HandleFunc("/auth/login/", hAuth.Login).Methods(http.MethodPost)
	WebRouter.HandleFunc("/auth/refresh/", hAuth.UpdateTokens).Methods(http.MethodGet)

	//Poll
	hPoll := poll_service.MustNewHandlerPoll()
	WebRouter.HandleFunc("/poll/add/", hPoll.Create).Methods(http.MethodPost)
	WebRouter.HandleFunc("/poll/update/{id}/", hPoll.Update).Methods(http.MethodPost)
	WebRouter.HandleFunc("/poll/vote/{id}/", hPoll.Vote).Methods(http.MethodGet).Name("protected_poll_vote")
	WebRouter.HandleFunc("/poll/view/{id}/", hPoll.View).Methods(http.MethodGet).Name("protected_poll_view")

	//Health
	WebRouter.HandleFunc("/events/health", hEvents.Health).Methods(http.MethodGet)
	WebRouter.HandleFunc("/events/health_protected/", hEvents.Health).Methods(http.MethodGet).Name("protected_health")
}
