package auth_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/Miroshinsv/disko_go/internal/auth-service/models"
	dbConnector "github.com/Miroshinsv/disko_go/pkg/db-connector"
	loggerService "github.com/Miroshinsv/disko_go/pkg/logger-service"
)

var (
	errorUnknownRole            = errors.New("unknown role")
	errorExistingUser           = errors.New("email already registered")
	errorInvalidEmail           = errors.New("invalid email address")
	errorInvalidPassword        = errors.New("invalid or empty password")
	errorInvalidEmailOrPassword = errors.New("invalid email or password")

	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Handler struct {
	log     loggerService.ILogger
	conn    dbConnector.IConnector
	service *Service
}

func (h Handler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Print(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Print(err.Error())
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	dbUser, err := h.service.RegisterUser(user)
	if err != nil {
		fmt.Print(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())

		return
	}

	_ = json.NewEncoder(w).Encode(dbUser)
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	var user models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	dbUser, err := h.service.LoginUser(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())

		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{
		"auth":    h.service.GenerateAuthJWT(dbUser),
		"refresh": h.service.GenerateRefreshJWT(dbUser),
	})
}

func (h Handler) UpdateTokens(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query()["refresh"]
	if len(token) == 0 || token[0] == "" {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	res, err := h.service.UpdateTokens(token[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)

		return
	}

	_ = json.NewEncoder(w).Encode(res)
}

func (h Handler) SocialAuth(w http.ResponseWriter, r *http.Request) {
	// получаем код от API VK из квери стринга
	authCode := r.URL.Query()["code"]
	if authCode[0] == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid code",
		})
	}

	user, err := h.service.LoginSocial(authCode[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
	}

	_ = json.NewEncoder(w).Encode(map[string]string{
		"auth":    h.service.GenerateAuthJWT(user),
		"refresh": h.service.GenerateRefreshJWT(user),
	})
}

func MustNewHandlerAuth() *Handler {
	db, _ := dbConnector.GetDBConnection()
	log := loggerService.GetLogger()

	return &Handler{
		log:     log,
		conn:    db,
		service: MustNewAuthService(log, db),
	}
}
