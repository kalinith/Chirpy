package main

import(
	//"fmt"
	"log"
	"net/http"
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

func main() {
	port := "8080"
	webpath := "./content" //root directory for the website. the lesson uses the programs dir for this.

	rootpath := http.Dir(webpath)
	rootHandler := http.FileServer(rootpath)

	serveMux := http.NewServeMux()
	serveMux.Handle("/app/",http.StripPrefix("/app",rootHandler)) //Static file content
	serveMux.HandleFunc("/healthz", health) //health check to see if site is ready to receive.

	server := &http.Server{
		Addr: ":" + port, //they used a constant for the port, this may be required at some point.
		Handler: serveMux,
	}
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}


/*
WIP examples go here

*/