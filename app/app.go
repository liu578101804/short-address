package app

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"encoding/json"
	"gopkg.in/validator.v2"
	"fmt"
	"github.com/justinas/alice"
)

type App struct {
	Router *mux.Router
	Middleware *Middleware
	Config *Env
}

func (a *App) Initialize(e *Env) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	a.Config = e
	a.Router = mux.NewRouter()
	a.Middleware = &Middleware{}
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	//初始化中间件
	m := alice.New(a.Middleware.LoggingHandler, a.Middleware.RecoverHandler)

	a.Router.Handle("/api/create", m.ThenFunc(a.createShortLink)).Methods("POST")
	a.Router.Handle("/api/info", m.ThenFunc(a.getShortLink)).Methods("GET")
	a.Router.Handle("/{short_link:[a-zA-Z0-9]{1,11}}", m.ThenFunc(a.redirect)).Methods("GET")
}

func (a *App) createShortLink(w http.ResponseWriter, r *http.Request) {
	var req ShortLinkReq
	if err := json.NewDecoder(r.Body).Decode(&req);err != nil{
		respondWithError(w, StatusError{http.StatusBadRequest, fmt.Errorf("parse parameters failed %v", r.Body)})
		return
	}
	if err := validator.Validate(req);err != nil{
		respondWithError(w, StatusError{http.StatusBadRequest, fmt.Errorf("validate parameters failed %v", r.Body)})
		return
	}
	defer r.Body.Close()

	//写入redis
	s,err := a.Config.S.Shorten(req.Url,req.ExpirationInMinute)
	if err != nil {
		respondWithError(w,err)
	}else{
		respondWithJson(w,http.StatusCreated,ShortLinkResp{
			ShortLink: s,
		})
	}
}

func (a *App) getShortLink(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	s := vals.Get("short_link")

	d, err := a.Config.S.ShortLinkInfo(s)
	if err != nil {
		respondWithError(w,err)
	}else{
		respondWithJson(w,http.StatusOK, d)
	}
}

func (a *App) redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	url,err := a.Config.S.UnShorten(vars["short_link"])
	if err != nil {
		respondWithError(w,err)
	}else{
		http.Redirect(w,r, url,http.StatusTemporaryRedirect)
	}
}

func (a *App) Run(addr string)  {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

//返回错误
func respondWithError(w http.ResponseWriter, err error)  {
	switch e := err.(type) {
	case Error:
		log.Printf("HTTP %d - %s",e.Status(),e)
		respondWithJson(w, e.Status(), e.Error())
	default:
		respondWithJson(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{})  {
	resp, _ := json.Marshal(Response{
		Code: code,
		Message: http.StatusText(code),
		Content: payload,
	})

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(code)
	w.Write(resp)
}

type Response struct {
	Code 		int 		`json:"code"`
	Message 	string 		`json:"message"`
	Content 	interface{} `json:"content"`
}