package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type healthchecker struct {
	callback chan func()
	statuses map[string]int
}

func newHealthchecker() *healthchecker {
	hc := &healthchecker{
		callback: make(chan func()),
		statuses: make(map[string]int),
	}

	go func() {
		go func() {
			for {
				(<-hc.callback)()
			}
		}()
	}()

	return hc
}

/* Handler method ServeHTTP for healthchecker . */
func (hc *healthchecker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sender := r.RemoteAddr

	log.Printf("[%s] GET /healthcheck: getting healthcheck", sender)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ch := make(chan []byte)

	hc.callback <- func() {
		snapshot, err := json.Marshal(hc.statuses)
		if err != nil {
			log.Printf("[%s] GET /healthcheck: unable to statuses users %q", sender, err.Error())

			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		log.Printf("[%s] GET /healthcheck: got healthcheck", sender)
		ch <- snapshot
	}

	w.Write(<-ch)
}

func (hc *healthchecker) increment(status int) {
	hc.callback <- func() {
		hc.statuses[fmt.Sprint(status)] += 1
	}
}
