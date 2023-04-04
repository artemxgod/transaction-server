package tnserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/artemxgod/transaction-server/internal/app/model"
	"github.com/artemxgod/transaction-server/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/segmentio/kafka-go"
)

var (
	errNotAuthenticated = errors.New("not authenticated")
	errNotEnoughFunds   = errors.New("not enough funds")
)

const (
	sessionName        = "TransactionSession"
	ctxKeyUser  ctxKey = iota
)

type ctxKey int8

type server struct {
	router       *mux.Router
	store        store.Store
	sessionStore sessions.Store
	UserMap	map[string]*model.User
}


func newServer(p_store store.Store, p_sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		store:        p_store,
		sessionStore: p_sessionStore,
		UserMap: make(map[string]*model.User),
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}



func (s *server) configureRouter() {
	s.router.HandleFunc("/users", s.handleUserCreate()).Methods("POST")
	s.router.HandleFunc("/login", s.handleLogin()).Methods("POST")

	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)

	private.HandleFunc("/transaction", s.handleTransaction()).Methods("POST")
	private.HandleFunc("/addfunds", s.handleAddFunds()).Methods("PATCH")
	private.HandleFunc("/removefunds", s.handleRemoveFunds()).Methods("PATCH")
	private.HandleFunc("/transaction/info", s.handleTransactionInfo()).Methods("GET")
}


func (s *server) handleLogin() http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByName(req.Name)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
		}

		if code, err := s.CreateSession(u, w, r); err != nil {
			s.error(w, r, code, err)
			return
		}

		s.respond(w, r, http.StatusOK, u)

	}
}

func (s *server) handleTransaction() http.HandlerFunc {
	type request struct {
		Funds float64 `json:"funds"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := r.Context().Value(ctxKeyUser).(*model.User)
		if u.Balance < req.Funds {
			s.error(w, r, http.StatusBadRequest, errNotEnoughFunds)
			return
		}

		rd := rand.New(rand.NewSource(time.Now().UnixNano()))

		trID := rd.Intn(1e6) + 1e6

		msg := fmt.Sprintf("Your order № %d for transaction was completed\n", trID)

		s.store.User().ChangeBalance(u.ID, -req.Funds)

		u.Writechan <- msg

		// we are expecting user to be logged in, so we check contect for him
		s.respond(w, r, http.StatusOK, map[string]int{"transactionID": trID})
	})
}

func (s *server) handleTransactionInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*model.User)

		var txt string
		tick := time.NewTicker(time.Millisecond * 1500)
		br := false
		for !br {
			select {
			case <-tick.C:
				fmt.Println("ticked")
				br = true
			case msg := <-u.Readchan:
				fmt.Println(msg)
				txt += msg
			}
		}

		fmt.Println(txt)

		s.respond(w, r, http.StatusOK, map[string][]string{"messages:": strings.Split(txt, "\n")})
	}
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// getting cashed session by name
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		// checking if user session exists
		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		// check if user is in database
		u, err := s.store.User().Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}



		if _, ok := s.UserMap[u.Name]; !ok {
			fmt.Println("first")
			u.Reader = kafka.NewReader(kafka.ReaderConfig{
				Brokers:  []string{"localhost:9092"},
				Topic:    u.Name,
				GroupID:  "group1",
				MinBytes: 5,
				MaxBytes: 10e6,
			})
	
			u.Writer = kafka.NewWriter(kafka.WriterConfig{
				Brokers: []string{"localhost:9092"},
				Topic:   u.Name,
			})
			u.Writechan = make(chan string, 100)
			u.Readchan = make(chan string, 100)

			s.UserMap[u.Name] = u

			go u.Write()
			go u.Read()
		} else {
			u = s.UserMap[u.Name]
		}
			// fmt.Println("first time")
			// u.Writechan = session.Values["writechan"].(chan string)
			// u.Readchan = rc

			// go u.Write()
			// go u.Read()
		

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *server) handleAddFunds() http.HandlerFunc {
	type request struct {
		Funds float64 `json:"funds"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := r.Context().Value(ctxKeyUser).(*model.User)
		edited_user, err := s.store.User().ChangeBalance(u.ID, req.Funds)
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, edited_user)
	}

}

func (s *server) handleRemoveFunds() http.HandlerFunc {
	type request struct {
		Funds float64 `json:"funds"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := r.Context().Value(ctxKeyUser).(*model.User)
		edited_user, err := s.store.User().ChangeBalance(u.ID, -req.Funds)
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, edited_user)
	}

}

func (s *server) handleUserCreate() http.HandlerFunc {
	type request struct {
		Name    string  `json:"name"`
		Balance float64 `json:"balance"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Name:    req.Name,
			Balance: req.Balance,
		}

		if err := s.store.User().CreateRecord(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		if code, err := s.CreateSession(u, w, r); err != nil {
			s.error(w, r, code, err)
			return
		}



		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) CreateSession(u *model.User, w http.ResponseWriter, r *http.Request) (int, error) {
	session, err := s.sessionStore.Get(r, sessionName)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	session.Values["user_id"] = u.ID
	if err := s.sessionStore.Save(r, w, session); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

// respond to http request
func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			fmt.Println(err)
		}
	}
}
