package routes

import (
	"log"
	"net/http"
)

func Setup(
	mux *http.ServeMux,
	logger *log.Logger,
) {
	mux.Handle("/api/test", TestHandle(logger))
}

func TestHandle(logger *log.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			logger.Println("test endpoint hit")
			w.Write([]byte("testing, testing, 1.. 2.. 3!"))
			// TODO Write returns int/error, how do I return err inside a handler?
		})
}
