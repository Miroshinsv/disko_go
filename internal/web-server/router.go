package web_server

import (
	auth_service "github.com/Miroshinsv/disko_go/internal/auth-service"
	city_service "github.com/Miroshinsv/disko_go/internal/city-service"
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

func RegisterHandlers() {
	WebRouter = mux.NewRouter()
	WebRouter.StrictSlash(true)

	WebRouter.Use(middleware.CORSMethodMiddleware(WebRouter))
	WebRouter.Use(middleware.AuthMiddleware)
	WebRouter.Use(middleware.AuthAdminMiddleware)
	WebRouter.Use(middleware.AuthSchoolAdminMiddleware)

	WebRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet, http.MethodOptions)

	//Directions
	hDirection := directionService.MustNewHandlerDirection()
	WebRouter.HandleFunc("/direction/get/all/", hDirection.GetAllDirections).Methods(http.MethodGet, http.MethodOptions)
	WebRouter.HandleFunc("/direction/get/{id}/", hDirection.GetDirectionById).Methods(http.MethodGet, http.MethodOptions)
	WebRouter.HandleFunc("/direction/update/{id}/", hDirection.UpdateDirectionById).Methods(http.MethodPost, http.MethodOptions).Name("protected_admin")
	WebRouter.HandleFunc("/direction/disband/{id}/", hDirection.DisbandDirectionById).Methods(http.MethodPost, http.MethodOptions).Name("protected_admin")
	WebRouter.HandleFunc("/direction/add/", hDirection.AddDirection).Methods(http.MethodPost, http.MethodOptions).Name("protected_admin")

	//Events
	hEvents := eventService.MustNewHandlerEvent()
	WebRouter.HandleFunc("/events/get/all/", hEvents.GetAllEvents).Methods(http.MethodGet, http.MethodOptions).Name("protected_school")
	WebRouter.HandleFunc("/events/get/{id}/", hEvents.GetEventById).Methods(http.MethodGet, http.MethodOptions)
	WebRouter.HandleFunc("/events/add/", hEvents.AddEvent).Methods(http.MethodPost, http.MethodOptions).Name("protected_school")
	WebRouter.HandleFunc("/events/update/{id}/", hEvents.UpdateEventById).Methods(http.MethodPost, http.MethodOptions).Name("protected_school")
	WebRouter.HandleFunc("/events/activate/{id}/", hEvents.ActivateEventById).Methods(http.MethodPost, http.MethodOptions).Name("protected_school")
	WebRouter.HandleFunc("/events/deactivate/{id}/", hEvents.DeactivateEventById).Methods(http.MethodPost, http.MethodOptions).Name("protected_school")
	WebRouter.HandleFunc("/events/disband/{id}/", hEvents.DeleteEventById).Methods(http.MethodPost, http.MethodOptions).Name("protected_admin")
	WebRouter.HandleFunc("/events/types/all/", hEvents.GetEventsType).Methods(http.MethodGet, http.MethodOptions)

	//Roles
	hRoles := roleService.MustNewHandlerRole()
	WebRouter.HandleFunc("/roles/get/all/", hRoles.GetAllRoles).Methods(http.MethodGet, http.MethodOptions)
	WebRouter.HandleFunc("/roles/get/{id}/", hRoles.GetRoleById).Methods(http.MethodGet, http.MethodOptions)
	WebRouter.HandleFunc("/roles/update/{id}/", hRoles.UpdateRoleById).Methods(http.MethodPost, http.MethodOptions).Name("protected_admin")
	WebRouter.HandleFunc("/roles/disband/{id}/", hRoles.DisbandRoleById).Methods(http.MethodPost, http.MethodOptions).Name("protected_admin")
	WebRouter.HandleFunc("/roles/add/", hRoles.AddRole).Methods(http.MethodPost, http.MethodOptions).Name("protected_admin")

	//School
	hSchool := schoolService.MustNewHandlerSchool()
	WebRouter.HandleFunc("/schools/get/all/", hSchool.GetAllSchools).Methods(http.MethodGet, http.MethodOptions)
	WebRouter.HandleFunc("/schools/get/{id}/", hSchool.GetSchoolById).Methods(http.MethodGet, http.MethodOptions)
	WebRouter.HandleFunc("/schools/delete/{id}/", hSchool.DeleteSchoolById).Methods(http.MethodDelete, http.MethodOptions).Name("protected_admin")
	WebRouter.HandleFunc("/schools/add/", hSchool.AddSchool).Methods(http.MethodPost, http.MethodOptions).Name("protected_admin")
	WebRouter.HandleFunc("/schools/update/{id}/", hSchool.UpdateSchoolById).Methods(http.MethodPost, http.MethodOptions).Name("protected_admin")

	//Users
	hUsers := userService.MustNewHandlerUser()
	WebRouter.HandleFunc("/user/get/all/", hUsers.GetAllUsers).Methods(http.MethodGet, http.MethodOptions).Name("protected_admin")
	WebRouter.HandleFunc("/user/get/{id}/", hUsers.GetUserById).Methods(http.MethodGet, http.MethodOptions)
	WebRouter.HandleFunc("/user/update/{id}/", hUsers.UpdateUserById).Methods(http.MethodPost, http.MethodOptions).Name("protected_admin")
	WebRouter.HandleFunc("/user/disband/{id}/", hUsers.DisbandUserById).Methods(http.MethodPost, http.MethodOptions).Name("protected_admin")
	WebRouter.HandleFunc("/user/add/", hUsers.AddUser).Methods(http.MethodPost, http.MethodOptions)

	//Schedule
	hSchedule := schedule_service.MustNewHandlerSchedule()
	WebRouter.HandleFunc("/schedule/today/", hSchedule.LoadEventsForToday).Methods(http.MethodGet, http.MethodOptions)
	WebRouter.HandleFunc("/schedule/all/", hSchedule.LoadAllEvents).Methods(http.MethodGet, http.MethodOptions)
	WebRouter.HandleFunc("/schedule/period/", hSchedule.LoadEventsForPeriod).Methods(http.MethodGet, http.MethodOptions)

	//Auth
	hAuth := auth_service.MustNewHandlerAuth()
	WebRouter.HandleFunc("/auth/register/", hAuth.Register).Methods(http.MethodPost, http.MethodOptions)
	WebRouter.HandleFunc("/auth/login/", hAuth.Login).Methods(http.MethodPost, http.MethodOptions)
	WebRouter.HandleFunc("/auth/refresh/", hAuth.UpdateTokens).Methods(http.MethodGet, http.MethodOptions)

	//Poll
	hPoll := poll_service.MustNewHandlerPoll()
	WebRouter.HandleFunc("/poll/add/", hPoll.Create).Methods(http.MethodPost, http.MethodOptions)
	WebRouter.HandleFunc("/poll/update/{id}/", hPoll.Update).Methods(http.MethodPost, http.MethodOptions).Name("protected_admin")
	WebRouter.HandleFunc("/poll/vote/{id}/", hPoll.Vote).Methods(http.MethodGet, http.MethodOptions).Name("protected_poll_vote")
	WebRouter.HandleFunc("/poll/view/{id}/", hPoll.View).Methods(http.MethodGet, http.MethodOptions).Name("protected_poll_view")
	//deprecated
	WebRouter.HandleFunc("/poll/count/{id}/", hPoll.ShowCount).Methods(http.MethodGet, http.MethodOptions)

	//Cities
	hCity := city_service.MustNewHandlerCities()
	WebRouter.HandleFunc("/city/add/", hCity.AddCity).Methods(http.MethodPost, http.MethodOptions).Name("protected_admin")
	WebRouter.HandleFunc("/city/all/", hCity.GetAllCities).Methods(http.MethodGet, http.MethodOptions)
	WebRouter.HandleFunc("/city/update/{id}/", hCity.UpdateCityById).Methods(http.MethodPost, http.MethodOptions).Name("protected_admin")
	WebRouter.HandleFunc("/city/get/{id}/", hCity.GetCityById).Methods(http.MethodGet, http.MethodOptions)
	WebRouter.HandleFunc("/city/disband/{id}/", hCity.DisbandCityById).Methods(http.MethodGet, http.MethodOptions).Name("protected_poll_view")

	//Health
	WebRouter.HandleFunc("/events/health", hEvents.Health).Methods(http.MethodGet, http.MethodOptions)
	WebRouter.HandleFunc("/events/health_protected/", hEvents.Health).Methods(http.MethodGet, http.MethodOptions).Name("protected_health")
}
