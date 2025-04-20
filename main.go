package main
import(
	"fmt"
	"net/http"
)

func main() {
	serveMux := http.NewServeMux()
	server := http.Server{
		Addr: ":8080",
		Handler: serveMux,
	}
	result := server.ListenAndServe()
	fmt.Println(result)
	
}
