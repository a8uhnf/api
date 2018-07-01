package controller

import (
	"encoding/json"
	"github.com/HackIllinois/api-checkin/models"
	"github.com/HackIllinois/api-checkin/service"
	"github.com/HackIllinois/api-commons/errors"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
)

func SetupController(route *mux.Route) {
	router := route.Subrouter()

	router.Handle("/{id}/", alice.New().ThenFunc(GetUserCheckin)).Methods("GET")
	router.Handle("/", alice.New().ThenFunc(CreateUserCheckin)).Methods("POST")
	router.Handle("/", alice.New().ThenFunc(UpdateUserCheckin)).Methods("PUT")
	router.Handle("/", alice.New().ThenFunc(GetCurrentUserCheckin)).Methods("GET")
	router.Handle("/qr/", alice.New().ThenFunc(GetQrCodeString)).Methods("GET")
}

/*
	Endpoint to get a specified user's checkin
*/
func GetUserCheckin(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	user_checkin, err := service.GetUserCheckin(id)

	if err != nil {
		panic(errors.UnprocessableError(err.Error()))
	}

	json.NewEncoder(w).Encode(user_checkin)
}

/*
	Endpoint to get the current user's checkin
*/
func GetCurrentUserCheckin(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("HackIllinois-Identity")

	user_checkin, err := service.GetUserCheckin(id)

	if err != nil {
		panic(errors.UnprocessableError(err.Error()))
	}

	json.NewEncoder(w).Encode(user_checkin)
}

/*
	Endpoint to set the specified user's checkin
*/
func CreateUserCheckin(w http.ResponseWriter, r *http.Request) {
	var user_checkin models.UserCheckin
	json.NewDecoder(r.Body).Decode(&user_checkin)

	isRsvped, err := service.IsAttendeeRsvped(user_checkin.ID)

	if err != nil {
		panic(errors.UnprocessableError(err.Error()))
	}

	if !isRsvped {
		panic(errors.UnprocessableError("Attendee must be rsvped to checkin"))
	}

	err = service.CreateUserCheckin(user_checkin.ID, user_checkin)

	if err != nil {
		panic(errors.UnprocessableError(err.Error()))
	}

	updated_checkin, err := service.GetUserCheckin(user_checkin.ID)

	if err != nil {
		panic(errors.UnprocessableError(err.Error()))
	}

	json.NewEncoder(w).Encode(updated_checkin)
}

/*
	Endpoint to update the specified user's checkin
*/
func UpdateUserCheckin(w http.ResponseWriter, r *http.Request) {
	var user_checkin models.UserCheckin
	json.NewDecoder(r.Body).Decode(&user_checkin)

	err := service.UpdateUserCheckin(user_checkin.ID, user_checkin)

	if err != nil {
		panic(errors.UnprocessableError(err.Error()))
	}

	updated_checkin, err := service.GetUserCheckin(user_checkin.ID)

	if err != nil {
		panic(errors.UnprocessableError(err.Error()))
	}

	json.NewEncoder(w).Encode(updated_checkin)
}

/*
	Endpoint to get the string to be embedded into the current user's QR code
*/
func GetQrCodeString(w http.ResponseWriter, r *http.Request) {	
	
	id := r.Header.Get("HackIllinois-Identity")
	
	uriString, err := service.GetQrString()

	if err != nil {
		panic(errors.UnprocessableError(err.Error()))
	}

	qr_string_container := models.QrStringContainer{
		ID: id,
		QrString: uriString,
	}

	json.NewEncoder(w).Encode(qr_string_container)	
}
