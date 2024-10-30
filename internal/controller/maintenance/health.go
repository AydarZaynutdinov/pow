package maintenance

import (
	"github.com/coreos/pkg/httputil"
	"net/http"
)

type Response struct {
	Status string `json:"status"`
}

func ReadinessHandler(w http.ResponseWriter, _ *http.Request) {
	_ = httputil.WriteJSONResponse(w, http.StatusOK, Response{
		Status: "ok",
	})
}

func LivenessHandler(w http.ResponseWriter, _ *http.Request) {
	_ = httputil.WriteJSONResponse(w, http.StatusOK, Response{
		Status: "ok",
	})
}
