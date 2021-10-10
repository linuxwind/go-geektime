package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)


type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}


func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         200,
		}
		handler.ServeHTTP(recorder, r)
		log.Printf("ClientIP: %s, HttpStatus: %d, URL: %s", r.Host, recorder.Status, r.URL)
	})
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404 NoFound")
	}
}

func IndexServer(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	w.Header().Add("Test", "TestAddHeader")


	var OsVersion string
	OsVersion = os.Getenv("VERSION")
	w.Header().Set("VERSION", OsVersion )

	fmt.Fprint(w, "Welcome to Home Page!")
}

func HealthZServer(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/healthz" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := "Status OK\n"

	fmt.Fprint(w, resp)
	fmt.Fprint(w,"Welcome to Healthz Page!")
}


func main() {

	http.HandleFunc("/", IndexServer)
	http.HandleFunc("/healthz", HealthZServer)

	err := http.ListenAndServe(":80", Log(http.DefaultServeMux))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
