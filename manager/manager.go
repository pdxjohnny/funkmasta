package manager

import (
	"fmt"
	"net/http"

	"github.com/pdxjohnny/funkmasta/api"
	// "github.com/pdxjohnny/funkmasta/backend"
)

func createHandler(w http.ResponseWriter, r *http.Request) {
	_, err := api.ParseCreate(r.Body)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}

	// backend.Create(s)
}

func main() {
	http.HandleFunc("/create/", createHandler)

	http.ListenAndServe(":8080", nil)
}
