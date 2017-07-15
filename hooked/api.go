package hooked

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Note: We eventually want to make this part of the context instead of a
// global variable, but this is OK for demonstration purposes...
var db *sql.DB

const (
	Success         = `{"msg": "OK"}`
	JSONContentType = "application/json; charset=UTF-8"
)

func getNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Write([]byte("Tried to get notification for id = " + vars["id"]))
}

func sendSuccess(w http.ResponseWriter) {
	w.Header().Set("Content-Type", JSONContentType)
	fmt.Fprint(w, Success)
}

func postActivityHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var a Activity
	err := decoder.Decode(&a)
	if err != nil {
		log.Printf("[ERROR] event=json-decode err=%q", err)
		http.Error(w, "Bad JSON body", http.StatusBadRequest)
		return
	}
	err = a.Validate(db)
	if err != nil {
		log.Printf("[ERROR] event=validation-error err=%q", err)
		http.Error(w, "Bad activity: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[DEBUG] Got new activity: %v", a)

	err = a.Save(db)
	if err != nil {
		log.Printf("[ERROR] event=saving-activity err=%q", err)
		http.Error(w, "Could not save activity: "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	// Send push notifications
	err = a.PushNotify()
	if err != nil {
		log.Printf("[ERROR] event=push-notification err=%q", err)
		http.Error(w, "Push notification error: "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	sendSuccess(w)
}

func Serve(d *sql.DB, port string) {
	db = d
	r := mux.NewRouter()
	r.HandleFunc("/user/{id}/notifications",
		getNotificationsHandler).Methods("GET")
	r.HandleFunc("/activity", postActivityHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8086", r))
}
