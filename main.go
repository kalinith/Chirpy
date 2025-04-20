package main

import(
	//"fmt"
	"log"
	"net/http"
)

func main() {
	port := "8080"
	webpath := "./content"

	rootpath := http.Dir(webpath)
	rootHandler := http.FileServer(rootpath)

	serveMux := http.NewServeMux()
	serveMux.Handle("/",rootHandler)

	server := &http.Server{
		Addr: ":" + port, //they used a constant for the port, this may be required at some point.
		Handler: serveMux,
	}
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}


/*
WIP examples go here

http.Handle("/", http.FileServer(http.Dir("/tmp")))

*/