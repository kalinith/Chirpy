package main

import(
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("nr hits now: %d\n",cfg.fileserverHits.Add(1))
		next.ServeHTTP(w, r)
	})
}

func main() {
	port := "8080"
	webpath := "./content" //root directory for the website. the lesson uses the programs dir for this.
	apiCfg := &apiConfig{}
	apiCfg.fileserverHits.Store(int32(0))

	rootpath := http.Dir(webpath)
	rootHandler := http.FileServer(rootpath)

	serveMux := http.NewServeMux()
	serveMux.Handle("/app/",apiCfg.middlewareMetricsInc(http.StripPrefix("/app",rootHandler))) //Static file content
	serveMux.HandleFunc("GET /healthz", Health) //health check to see if site is ready to receive.
	serveMux.HandleFunc("GET /metrics", apiCfg.Stats) //show the server statistics
	serveMux.HandleFunc("POST /reset", apiCfg.Reset) 

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