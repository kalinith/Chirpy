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
        //User_id uuid.UUID `json:"user_id"`
    }
 	type returnVals struct {
        ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
    }

    tokenString, tokenErr := auth.GetBearerToken(req.Header)
    if tokenErr != nil {
    	returnError(w, 401, "Unauthorized", tokenErr)
    	return
    }

    userID, validErr := auth.ValidateJWT(tokenString, cfg.jwt_Secret)
    if validErr != nil {
    	returnError(w, 401, "Unauthorized", validErr)
    	return
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

   	chirpParams := database.CreateChirpParams{cleanString(params.Body), userID}
   	//log.Printf("Body: %s, user_id: %s", chirpParams.Body, chirpParams.UserID)
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
    	return
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
	returns.IsChirpyRed = dbUser.IsChirpyRed

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

    //Generate Refresh Token
	utcTime := time.Now().UTC()
	expiresIn, _ := time.ParseDuration("60d")
	expiresAt := utcTime.Add(expiresIn)
	refreshToken, _ := auth.MakeRefreshToken()

	//insert into table for client.
	refreshTokenParams := database.CreateRefreshTokenParams{refreshToken, dbUser.ID, expiresAt}
	refreshTokenRecord, qryErr := cfg.dbQueries.CreateRefreshToken(req.Context(), refreshTokenParams)
	if qryErr != nil {
		returnError(w, 500, "failed to save refreshToken", qryErr)
		return
	}

    //Generate Access token from data
	token, tokenErr := auth.MakeJWT(dbUser.ID, cfg.jwt_Secret)
	if tokenErr != nil {
		returnError(w, 500, "failed to generate token", err)
		return
	}

	returns := User{}
 	returns.ID = dbUser.ID
	returns.CreatedAt = dbUser.CreatedAt
	returns.UpdatedAt = dbUser.UpdatedAt
	returns.Email = dbUser.Email
	returns.Token = token
	returns.RefreshToken = refreshTokenRecord.Token
	returns.IsChirpyRed = dbUser.IsChirpyRed

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

func (cfg *apiConfig) refresh(w http.ResponseWriter, req *http.Request) {
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		returnError(w, 401, "Invalid or expired token", err)
		return
	}
	
	RefreshToken, err := cfg.dbQueries.GetToken(req.Context(), bearerToken)
	if err != nil {
		returnError(w, 401, "Invalid or expired token", err)
		return
	}

	if RefreshToken.RevokedAt.Valid {  // Assuming RevokedAt is a sql.NullTime
		returnError(w, 401, "Invalid or expired token", err)
		return
	}
	

	//Generate New Access token from user ID
	returns := struct {
		Token string `json:"token"`
	}{}

	returns.Token, err = auth.MakeJWT(RefreshToken.UserID, cfg.jwt_Secret)
	if err != nil {
		returnError(w, 500, "failed to generate token", err)
		return
	}

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

func (cfg *apiConfig) revoke(w http.ResponseWriter, req *http.Request) {
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		returnError(w, 401, "Invalid or expired token", err)
		return
	}
	
	_ = cfg.dbQueries.RevokeToken(req.Context(), bearerToken)
	w.WriteHeader(204)
}

func (cfg *apiConfig) passwordUpdate(w http.ResponseWriter, req *http.Request) {
	//First check that the auth token in the header is valid
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		returnError(w, 401, "Invalid token", err)
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.jwt_Secret)
	if err != nil {
		returnError(w, 401, "Invalid token", err)
		return
	}

	//pull the parameters from the request body
	params := struct {
 		Password string	`json:"password"` //the new password for the client
  		Email 	 string	`json:"email"` //the new e-mail address for the client
	}{}
    
    decoder := json.NewDecoder(req.Body)
    err = decoder.Decode(&params)
    if err != nil {
		returnError(w, 500, "Invalid json parameters", err)
		return
    }

    //configure parameterd for user update    
    userParam := database.UpdateUserParams{}
    userParam.ID = userID
    userParam.Email = params.Email
    userParam.HashedPassword, err = auth.HashPassword(params.Password)//has the password before saving it
    if err != nil {
    	returnError(w, 401, "Incorrect email address", err)
    	return
    }

    //update username and password
    dbUser, err := cfg.dbQueries.UpdateUser(req.Context(), userParam)
    if err != nil {
    	returnError(w, 401, "not able to update user", err)
    	return
    }

	returns := User{}
 	returns.ID = dbUser.ID
	returns.CreatedAt = dbUser.CreatedAt
	returns.UpdatedAt = dbUser.UpdatedAt
	returns.Email = dbUser.Email
	returns.IsChirpyRed = dbUser.IsChirpyRed

    dat, err := json.Marshal(returns)
	if err != nil {
		returnError(w, 500, "Error marshalling JSON", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
    w.Write(dat)   
}


func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, req *http.Request) {
	//First check that the auth token in the header is valid
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		returnError(w, 401, "Invalid token", err)
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.jwt_Secret)
	if err != nil {
		returnError(w, 401, "Invalid token", err)
		return
	}

	//get ID from path
	id, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		returnError(w, 500, "Invalid ID format", err)
		return
	}

	type chirp struct {
        ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	//execute db query to fetch chirp.
	fetchedChirp, err := cfg.dbQueries.GetChirp(req.Context(), id)
	if err != nil {
		returnError(w, 404, "Chirp does not exist", err)
		return
	}

	if userID != fetchedChirp.UserID {
		returnError(w, 403, "Unauthorized", nil)
		return
	}

	err = cfg.dbQueries.DeleteChirp(req.Context(), id)
	if err != nil {
		returnError(w, 403, "Unauthorized", nil)
		return
	}
	w.WriteHeader(204)
}

func (cfg *apiConfig) polkaWebhooks(w http.ResponseWriter, req *http.Request) {
	//read header to get token
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		returnError(w, 401, "Invalid API Key", nil)
		return
	}

	if apiKey != cfg.polkaKey {
		//on error respond with 404
		returnError(w, 401, "Invalid API Key", nil)
		return
	}
	//read parameters out of the body of the request
	params := struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}{}
    
    decoder := json.NewDecoder(req.Body)
    err = decoder.Decode(&params)
    if err != nil {
		returnError(w, 500, "Invalid json parameters", err)
		return
    }
	//check if event is "user.upgraded" if not respond with 204
	if params.Event == "user.upgraded" {
		//upgrade user to red
		uID, err := uuid.Parse(params.Data.UserID)
		if err != nil {
			returnError(w, 404, "User not found", err)
			return
		}
		_, dberr := cfg.dbQueries.UpdateRed(req.Context(), uID)
		if dberr != nil {
			//on error respond with 404
			returnError(w, 404, "User not found", dberr)
			return
		}
	}
	w.WriteHeader(204)
}

