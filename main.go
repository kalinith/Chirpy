package main

import _ "github.com/lib/pq"
import(
	"os"
	"log"
	"net/http"
	"database/sql"
	"github.com/joho/godotenv"
	"Chirpy/internal/database"
)

func main() {
	godotenv.Load()
	port := "8080"
	webpath := "./content" //root directory for the website. the lesson uses the programs dir for this.
	apiCfg := &apiConfig{}
	apiCfg.fileserverHits.Store(int32(0))
	dbURL := os.Getenv("DB_URL")
	apiCfg.platform = os.Getenv("PLATFORM")
	apiCfg.jwt_Secret = os.Getenv("JWT_SECRET")
	
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("unable to open DB")
	}
	apiCfg.dbQueries = database.New(db)

	rootpath := http.Dir(webpath)
	rootHandler := http.FileServer(rootpath)

	serveMux := http.NewServeMux()
	serveMux.Handle("/app/",apiCfg.middlewareMetricsInc(http.StripPrefix("/app",rootHandler))) //Static file content
	serveMux.HandleFunc("GET /api/healthz", health) //health check to see if site is ready to receive.
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.metrics) //show the server statistics
	serveMux.HandleFunc("POST /admin/reset", apiCfg.reset) //reset metrics
	serveMux.HandleFunc("POST /api/users", apiCfg.addUser) //add a Chirp user based on e-mail
	serveMux.HandleFunc("POST /api/chirps", apiCfg.addChirp) //add a Chirp
	serveMux.HandleFunc("GET /api/chirps", apiCfg.getChirps) //fetch all chirps
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirp) //fetch one chirp by ID
	serveMux.HandleFunc("POST /api/login", apiCfg.login)//login user will eventually return a token.
	serveMux.HandleFunc("POST /api/refresh", apiCfg.refresh)//Refresh the access token.
	serveMux.HandleFunc("POST /api/revoke", apiCfg.revoke)//revoke refresh token
	serveMux.HandleFunc("PUT /api/users", apiCfg.passwordUpdate)//revoke refresh token
	serveMux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.deleteChirp)//Allow a user to delete a Chirp he owns
	serveMux.HandleFunc("POST /api/polka/webhooks", apiCfg.polkaWebhooks)//revoke refresh token

	server := &http.Server{
		Addr: ":" + port, //they used a constant for the port, this may be required at some point.
		Handler: serveMux,
	}
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}


/*
WIP examples go here

 mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
*/