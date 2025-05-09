package main

import(
	"fmt"
	"log"
	"time"
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"Chirpy/internal/database"
	"Chirpy/internal/auth"
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

func returnError(w http.ResponseWriter, errHeader int, errMessage string, err error) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8") // normal header
	w.WriteHeader(errHeader)
	w.Write([]byte(errMessage))
	log.Printf("%s: %s", errMessage, err)
	return
}

func (cfg *apiConfig) metrics(w http.ResponseWriter, req *http.Request) {
	//stat1 := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/html; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf("<html>\n  <body>\n    <h1>Welcome, Chirpy Admin</h1>\n    <p>Chirpy has been visited %d times!</p>\n  </body>\n</html>", cfg.fileserverHits.Load())

	w.Write([]byte(body))
}

func (cfg *apiConfig) reset(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		returnError(w, 403, "Forbidden", fmt.Errorf("Non Dev call to reset system"))
		return
	}
	err := cfg.dbQueries.DeleteUsers(req.Context())
	if err != nil {
		returnError(w, 500, "error: Unable to reset user table", err)
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

func (cfg *apiConfig) addChirp(w http.ResponseWriter, req *http.Request) {
    type parameters struct {
        Body string `json:"body"`
        User_id uuid.UUID `json:"user_id"`
    }
 	type returnVals struct {
        ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
    }

    w.Header().Set("Content-Type", "application/json")
    params := parameters{}
    returns := returnVals{}

    decoder := json.NewDecoder(req.Body)
    err := decoder.Decode(&params)
    if err != nil {
   		returnError(w, 500, "unable to decode parameters", err)
   		return
    }

   	if len(params.Body) > 140 {
   		returnError(w, 400, "Chirp is too long", nil)
   		return
   	}

   	chirpParams := database.CreateChirpParams{cleanString(params.Body), params.User_id}
   	log.Printf("Body: %s, user_id: %s", chirpParams.Body, chirpParams.UserID)
   	chirp, err := cfg.dbQueries.CreateChirp(req.Context(), chirpParams)
   	if err != nil {
   		returnError(w, 500, "Chirp not saved", err)
   		return
   	}

   	returns.ID = chirp.ID
   	returns.CreatedAt = chirp.CreatedAt
   	returns.UpdatedAt = chirp.UpdatedAt
   	returns.Body = chirp.Body
   	returns.UserID = chirp.UserID

    dat, err := json.Marshal(returns)
	if err != nil {
			returnError(w, 500, "Error marshalling JSON", err)
			return
	}
	w.WriteHeader(201)
    w.Write(dat)

}

func (cfg *apiConfig) addUser(w http.ResponseWriter, req *http.Request) {
    type parameters struct {
        // these tags indicate how the keys in the JSON should be mapped to the struct fields
        // the struct fields must be exported (start with a capital letter) if you want them parsed
		Email string `json:"email"`
		Password string `json:"password"`
	}

    params := parameters{}
    returns := User{}

    decoder := json.NewDecoder(req.Body)
    err := decoder.Decode(&params)
    if err != nil {
		returnError(w, 500, "error decoding Add User Request Parameters", err)
		return
    }
    hashpassword, err := auth.HashPassword(params.Password)
    if err != nil {
    	returnError(w, 500, "error encrypting password", err)
    }

    dbparam := database.CreateUserParams{}
    dbparam.Email = params.Email
    dbparam.HashedPassword = hashpassword
    dbUser, err := cfg.dbQueries.CreateUser(req.Context(), dbparam)
	if err != nil {
		returnError(w, 500, "error: Unable to add user", err)
		return
	}

	returns.ID = dbUser.ID
	returns.CreatedAt = dbUser.CreatedAt
	returns.UpdatedAt = dbUser.UpdatedAt
	returns.Email = dbUser.Email

    //formulate reply here
    dat, err := json.Marshal(returns)
	if err != nil {
		returnError(w, 500, "Error marshalling JSON", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
    w.Write(dat)
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, req *http.Request) {
	type chirp struct {
        ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	fetchedChirps, err := cfg.dbQueries.GetChirps(req.Context())
	if err != nil {
		returnError(w, 500, "unable to fetch Chirps", err)
		return
	}
	returnChirps := []chirp{}

	for _, fetchedChirp := range fetchedChirps {
		newChirp := chirp{}
		newChirp.ID = fetchedChirp.ID
   		newChirp.CreatedAt = fetchedChirp.CreatedAt
		newChirp.UpdatedAt = fetchedChirp.UpdatedAt
		newChirp.Body = fetchedChirp.Body
		newChirp.UserID = fetchedChirp.UserID
		returnChirps = append(returnChirps, newChirp)
	}

    dat, err := json.Marshal(returnChirps)
	if err != nil {
		returnError(w, 500, "Error marshalling JSON", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
    w.Write(dat)
}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, req *http.Request) {
	type chirp struct {
        ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	//get ID from path
	id, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		returnError(w, 500, "Invalid ID format", err)
		return
	}

	//execute db query to fetch chirp.
	fetchedChirp, err := cfg.dbQueries.GetChirp(req.Context(), id)
	if err != nil {
		returnError(w, 404, "Chirp does not exist", err)
		return
	}

	newChirp := chirp{}
	newChirp.ID = fetchedChirp.ID
	newChirp.CreatedAt = fetchedChirp.CreatedAt
	newChirp.UpdatedAt = fetchedChirp.UpdatedAt
	newChirp.Body = fetchedChirp.Body
	newChirp.UserID = fetchedChirp.UserID

    dat, err := json.Marshal(newChirp)
	if err != nil {
		returnError(w, 500, "Error marshalling JSON", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
    w.Write(dat)
}

func (cfg *apiConfig) login(w http.ResponseWriter, req *http.Request) {
	params := struct {
 		Password string	`json:"password"`
  		Email 	 string	`json:"email"`
	}{}
    
    decoder := json.NewDecoder(req.Body)
    err := decoder.Decode(&params)
    if err != nil {
		returnError(w, 500, "Invalid json parameters", err)
		return
    }

    //fetch the users hash
    dbUser := database.User{}
    dbUser, err = cfg.dbQueries.GetUser(req.Context(), params.Email)
    if err != nil {
    	returnError(w, 401, "Incorrect email or password", err)
    }
    //check password against hash
    err = auth.CheckPasswordHash(dbUser.HashedPassword, params.Password)
    if err != nil {
    	returnError(w, 401, "Incorrect email or password", err)
    }

    //in a future lesson assign token here
	returns := User{}
 	returns.ID = dbUser.ID
	returns.CreatedAt = dbUser.CreatedAt
	returns.UpdatedAt = dbUser.UpdatedAt
	returns.Email = dbUser.Email

    //formulate reply here
    dat, err := json.Marshal(returns)
	if err != nil {
		returnError(w, 500, "Error marshalling JSON", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
    w.Write(dat)
}
