package main

import(
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	//"github.com/google/uuid"
	//"Chirpy/internal/database"
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
	if cfg.platform != "dev" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8") // normal header
		w.WriteHeader(403)
		w.Write([]byte("Forbidden"))
		log.Printf("Non Dev call to reset system")
		return
	}
	err := cfg.dbQueries.DeleteUsers(req.Context())
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8") // normal header
		w.WriteHeader(500)
		w.Write([]byte("error: Unable to reset user table"))
		log.Printf("Unable to remove users error: %s", err)
		return
	}
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)

	//Update the POST /admin/reset endpoint to delete all users in the database
	// (but don't mess with the schema). You'll need a new SQLC query for this.
	// Add a new value to your .env file called PLATFORM and set it equal to "dev".
	// Read it into your apiConfig. If PLATFORM is not equal to "dev", this endpoint 
	// should return a 403 Forbidden. This ensures that this extremely dangerous endpoint
	// can only be accessed in a local development environment.

}

func (cfg *apiConfig) validate_Chirp(w http.ResponseWriter, req *http.Request) {
    type parameters struct {
        // these tags indicate how the keys in the JSON should be mapped to the struct fields
        // the struct fields must be exported (start with a capital letter) if you want them parsed
        Body string `json:"body"`
    }
 	type returnVals struct {
        // the key will be the name of struct field unless you give it an explicit JSON tag
        Err string `json:"error"`
        Valid bool `json:"valid"`
        Cleaned_body string `json:"cleaned_body"`
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
		returns.Cleaned_body = ""
		header = 500
    }
    if err == nil {
    	if len(params.Body) > 140 {
    		header = 400
    		returns.Err = "Chirp is too long"
    		returns.Valid = false
    		returns.Cleaned_body = ""
    	} else {
    		header = 200
    		returns.Valid = true
    		returns.Cleaned_body = cleanString(params.Body)
    	}
    }

    dat, err := json.Marshal(returns)
	if err != nil {
			dat, _ = json.Marshal(returnVals{fmt.Sprintf("Error marshalling JSON: %s", err), false, ""})
			header = 500 
	}
	w.WriteHeader(header)
    w.Write(dat)

}

func (cfg *apiConfig) addUser(w http.ResponseWriter, req *http.Request) {
    type parameters struct {
        // these tags indicate how the keys in the JSON should be mapped to the struct fields
        // the struct fields must be exported (start with a capital letter) if you want them parsed
		Email string `json:"email"`
	}

    params := parameters{}
    returns := User{}

    decoder := json.NewDecoder(req.Body)
    err := decoder.Decode(&params)
    if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8") // normal header
		w.WriteHeader(500)
		w.Write([]byte("error decoding Add User Request Parameters"))
		log.Printf("error decoding Add User Request: %s", err)
		return
    }

    //create user here using params.Email
    dbUser, err := cfg.dbQueries.CreateUser(req.Context(), params.Email)
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8") // normal header
		w.WriteHeader(500)
		w.Write([]byte("error: Unable to add user"))
		log.Printf("Unable to add user error: %s", err)
		return
	}

	returns.ID = dbUser.ID
	returns.CreatedAt = dbUser.CreatedAt
	returns.UpdatedAt = dbUser.UpdatedAt
	returns.Email = dbUser.Email

    //formulate reply here
    dat, err := json.Marshal(returns)
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8") // normal header
		w.WriteHeader(500)
		w.Write([]byte("Error marshalling JSON"))
		log.Printf("Error marshalling JSON: %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
    w.Write(dat)



}