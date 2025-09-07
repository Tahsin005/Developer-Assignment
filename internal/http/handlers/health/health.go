package health

import (
	"fmt"
	"net/http"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}
