package main

import(
	"fmt"
	"net/http"
	"encoding/json"
)

func health(w http.ResponseWriter, req *http.Request) {
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


func (cfg *apiConfig) metrics(w http.ResponseWriter, req *http.Request) {
	//stat1 := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/html; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf("<html>\n  <body>\n    <h1>Welcome, Chirpy Admin</h1>\n    <p>Chirpy has been visited %d times!</p>\n  </body>\n</html>", cfg.fileserverHits.Load())

	w.Write([]byte(body))
}

func (cfg *apiConfig) reset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) validate(w http.ResponseWriter, req *http.Request) {
    type parameters struct {
        // these tags indicate how the keys in the JSON should be mapped to the struct fields
        // the struct fields must be exported (start with a capital letter) if you want them parsed
        Body string `json:"body"`
    }
 	type returnVals struct {
        // the key will be the name of struct field unless you give it an explicit JSON tag
        Err string `json:"error"`
        Valid bool `json:"valid"`
    }

    w.Header().Set("Content-Type", "application/json")
    header := 404
    params := parameters{}
    returns := returnVals{}

    decoder := json.NewDecoder(req.Body)
    err := decoder.Decode(&params)
    if err != nil {
        // an error will be thrown if the JSON is invalid or has the wrong types
        // any missing fields will simply have their values in the struct set to their zero value
		returns.Err = "Something went wrong"
		returns.Valid = false
		header = 500
    }
    if err == nil {
    	if len(params.Body) > 140 {
    		header = 400
    		returns.Err = "Chirp is too long"
    		returns.Valid = false
    	} else {
    		header = 200
    		returns.Valid = true
    	}
    }

    dat, err := json.Marshal(returns)
	if err != nil {
			dat, _ = json.Marshal(returnVals{fmt.Sprintf("Error marshalling JSON: %s", err), false})
			header = 500 
	}
	w.WriteHeader(header)
    w.Write(dat)

}
