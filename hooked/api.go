package hooked

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	JSONContentType = "application/json; charset=UTF-8"
)

func getNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Write([]byte("Tried to get notification for id = " + vars["id"]))
}

func postActivityHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] event=read-body err=%q", err)
		http.Error(w, "Bad body", http.StatusBadRequest)
		return
	}
	w.Write([]byte("Tried to post activity " + string(body)))
}

func Serve(port string) {
	r := mux.NewRouter()
	r.HandleFunc("/user/{id}/notifications",
		getNotificationsHandler).Methods("GET")
	r.HandleFunc("/activity", postActivityHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8086", r))
}
