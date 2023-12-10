package save

import (
	"fmt"
	"net/http"
)

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "TODO: save")
}
