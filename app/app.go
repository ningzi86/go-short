package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go-short/serror"
	"github.com/justinas/alice"
	 "go-short/middleware"

	validator2 "gopkg.in/validator.v2"
	"go-short/env"
)

type App struct {
	Router       *mux.Router
	Middlerwares *middleware.Middlerware
	Config       *env.Env
}

type shortenRequ struct {
	URL     string `json:"url" validate:"nonzero"`
	Expired int64  `json:"expired" validate:"min=0"`
}

type shortResp struct {
	ShortUrl string `json:"short_url"`
}

func (app *App) Init(e *env.Env) {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app.Router = mux.NewRouter()
	app.Middlerwares = &middleware.Middlerware{}
	app.Config = e
	app.initializeRoutes()

}

func (app *App) initializeRoutes() {

	m := alice.New(app.Middlerwares.LogginHandler, app.Middlerwares.RecoverHandler)

	app.Router.Handle("/api/shorten", m.ThenFunc(app.createShortUrl)).Methods("POST")
	app.Router.Handle("/api/info", m.ThenFunc(app.getShortUrl)).Methods("GET")
	app.Router.Handle("/{shortUrl:[0-9a-zA-Z]{1,11}}", m.ThenFunc(app.redirect)).Methods("GET")

}

func (app *App) createShortUrl(w http.ResponseWriter, r *http.Request) {

	var requ shortenRequ
	if err := json.NewDecoder(r.Body).Decode(&requ); err != nil {
		responseWithError(w, serror.StatusError{http.StatusBadRequest, fmt.Errorf("parse parameters failed %v", r.Body)})
		return
	}
	if err := validator2.Validate(requ); err != nil {
		responseWithError(w, serror.StatusError{http.StatusBadRequest, fmt.Errorf("validate parameters failed %v", requ)})
		return
	}
	defer r.Body.Close()

	s, err := app.Config.S.Shorten(requ.URL, requ.Expired)

	if err != nil {
		responseWithError(w, err)
	} else {
		responseWithJSON(w, http.StatusCreated, &shortResp{ShortUrl: s})
	}

}

func responseWithError(w http.ResponseWriter, err error) {

	switch e := err.(type) {
	case serror.Error:
		log.Printf("HTTP ^%d - %s", e.Status(), e)
		responseWithJSON(w, e.Status(), e.Error())

	default:
		responseWithJSON(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

}

func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	resp, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)

}

func (app *App) getShortUrl(w http.ResponseWriter, r *http.Request) {

	vars := r.URL.Query()
	s := vars.Get("shortUrl")

	d, err := app.Config.S.ShortenInfo(s)
	if err != nil {
		responseWithError(w, err)
	} else {
		responseWithJSON(w, http.StatusOK, d)
	}

}

func (app *App) redirect(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	u, err := app.Config.S.Unshorten(vars["shortUrl"])
	if err != nil {
		responseWithError(w, err)
	} else {
		http.Redirect(w, r, u, http.StatusTemporaryRedirect)
	}

}

func (app *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, app.Router))
}
