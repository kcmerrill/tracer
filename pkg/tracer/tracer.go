package tracer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Tracer is the main
type Tracer struct {
	auth   string
	port   string
	panic  string
	lock   *sync.Mutex
	Checks map[string]*check
}

// Start a new Tracer server
func Start(auth, port, panic string) {
	t := &Tracer{
		auth:   auth,
		port:   port,
		panic:  panic,
		Checks: make(map[string]*check),
		lock:   &sync.Mutex{},
	}
	t.serve()
}

func (t *Tracer) clearCheck(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	if c, exists := t.Checks[vars["check"]]; exists {
		// woot
		c.ok()
		cJSON, _ := json.Marshal(c)
		resp.WriteHeader(http.StatusOK)
		fmt.Fprint(resp, string(cJSON))
		t.lock.Lock()
		defer t.lock.Unlock()
		// bye felicia
		delete(t.Checks, vars["check"])
	} else {
		// not found
		http.NotFound(resp, req)
	}
}

func (t *Tracer) createCheck(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	t.lock.Lock()
	defer t.lock.Unlock()

	// does it exist? lets cancel it .... <-- controversial don't @ me
	if c, exists := t.Checks[vars["check"]]; exists {
		// cancel it ... so we can create the new one
		c.cancel()
	}

	panic, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(resp, fmt.Sprintf(`{"error": "Cannot read request body","check":"%s"}"`, vars["check"]))
		return
	}

	// create a new check, can we use the body?
	var c *check
	if len(panic) == 0 {
		c = initCheck(vars["check"], vars["duration"], t.panic)
	} else {
		c = initCheck(vars["check"], vars["duration"], string(panic))
	}

	t.Checks[vars["check"]] = c
	cJSON, _ := json.Marshal(c)
	resp.WriteHeader(http.StatusOK)
	fmt.Fprint(resp, string(cJSON))
	return
}

func (t *Tracer) serve() {
	addr := []string{"0.0.0.0", t.port}
	if strings.Contains(t.port, ":") {
		addr = strings.Split(t.port, ":")
	}

	r := mux.NewRouter()
	r.HandleFunc(`/{check}`, t.httpAuth(t.auth, t.clearCheck)).Methods("GET")
	r.HandleFunc(`/{check}/in/{duration}`, t.httpAuth(t.auth, t.createCheck)).Methods("GET", "POST", "PUT")
	r.HandleFunc(`/{check}/{duration}`, t.httpAuth(t.auth, t.createCheck)).Methods("GET", "POST", "PUT")
	srv := &http.Server{
		Handler:      r,
		Addr:         strings.Join(addr, ":"),
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	log.WithFields(log.Fields{
		"binding": t.port,
	}).Info("Starting Tracer ...")
	if err := srv.ListenAndServe(); err != nil {
		log.WithFields(log.Fields{
			"binding": t.port,
		}).Error(err.Error())
	}
}

func (t *Tracer) httpAuth(password string, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, _, _ := r.BasicAuth()
		if password != "" && password != token {
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}
		fn(w, r)
	}
}
