package main

import(
	"fmt"
	"net/http"
)

func Health(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
/*
The endpoint should simply return a 200 OK status code indicating that it has
started up successfully and is listening for traffic. The endpoint should return
a Content-Type: text/plain; charset=utf-8 header, and the body will contain a
message that simply says "OK" (the text associated with the 200 status code).
*/


func (cfg *apiConfig) Stats(w http.ResponseWriter, req *http.Request) {
	stat1 := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(stat1))
}

func (cfg *apiConfig) Reset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)
}



