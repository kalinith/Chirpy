package main

import(
	"log"
	"net/http"
	"sync/atomic"
	"strings"
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

func cleanString(s string) string {
	if s == "" {
		return s
	}
	cussWords := [3]string{"kerfuffle","sharbert", "fornax"}

	words := strings.Split(s, " ")
	for i := 0; i < len(words); i++ {
		for _, cuss := range cussWords {
			if strings.ToLower(words[i]) == cuss {
				words[i] = "****"
				break
			}
		}
	}

	return strings.Join(words, " ")
}