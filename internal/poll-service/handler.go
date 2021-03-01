package poll_service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Miroshinsv/disko_go/internal/poll-service/models"
	userService "github.com/Miroshinsv/disko_go/internal/user-service"
	dbConnector "github.com/Miroshinsv/disko_go/pkg/db-connector"
	loggerService "github.com/Miroshinsv/disko_go/pkg/logger-service"
)

type Handler struct {
	log     loggerService.ILogger
	conn    dbConnector.IConnector
	service *Service
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	var income models.Income
	err = json.Unmarshal(body, &income)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	err = income.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)

		return
	}

	poll, err := h.service.Create(income)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	_ = json.NewEncoder(w).Encode(poll)
}

func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	var poll = &models.Poll{}
	err = json.Unmarshal(body, &poll)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid id")

		return
	}

	var existing = &models.Poll{}
	h.conn.GetConnection().Where(fmt.Sprintf("id=%d", i)).Find(&existing)
	if existing.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Unknown id")

		return
	}

	db := h.conn.GetConnection().Model(&existing).Updates(poll)
	if db.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(db.Error.Error()))

		return
	}

	_ = json.NewEncoder(w).Encode(poll)
}

func (h Handler) Vote(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid id")
		return
	}

	var poll = &models.Poll{}
	db := h.conn.GetConnection().Where(fmt.Sprintf("id=%d", i)).Find(poll)
	if db.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("invalid poll ID")
		return
	}
	err = h.service.Vote(poll, r.Context().Value("user").(*userService.Users))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}
	_ = json.NewEncoder(w).Encode("ok")
}

func (h Handler) View(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid id")

		return
	}

	var poll = &models.Poll{}
	db := h.conn.GetConnection().Where(fmt.Sprintf("id=%d", i)).Find(poll)
	if db.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("invalid poll ID")

		return
	}

	res, err := h.service.ShowResults(poll, r.Context().Value("user").(*userService.Users))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)

		return
	}

	_ = json.NewEncoder(w).Encode(res)
}

func (h Handler) ShowCount(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid id")

		return
	}

	var poll = &models.Poll{}
	db := h.conn.GetConnection().Where(fmt.Sprintf("id=%d", i)).Find(poll)
	if db.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("invalid poll ID")

		return
	}

	res, err := h.service.ShowVotesCount(poll)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode("internal error")

		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"poll":  poll,
		"votes": res,
	})
}

func MustNewHandlerPoll() *Handler {
	db, _ := dbConnector.GetDBConnection()
	log := loggerService.GetLogger()

	return &Handler{
		log:     log,
		conn:    db,
		service: MustNewPollService(log, db),
	}
}
