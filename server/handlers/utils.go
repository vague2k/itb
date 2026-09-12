package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"itb.ihatedoing.work/server/services"
	"itb.ihatedoing.work/views/shadcn/toast"
)

func handleError(log *services.LoggingService, w http.ResponseWriter, r *http.Request, err error) {
	// Redirect the response into the toast viewport and append an SSR stub so
	// the bundled toast script adopts it. The caller's swap target is untouched.
	w.Header().Set("HX-Retarget", "[data-tui-toaster]")
	w.Header().Set("HX-Reswap", "beforeend")

	var ue services.UserError
	var ie services.InternalError

	switch {
	case errors.As(err, &ue):
		renderErrorToast(w, r, ue.Error())
	case errors.As(err, &ie):
		log.Error(r, ie.Error(), ie.Unwrap())
		renderErrorToast(w, r, "Internal error. Please try again later.")
	default:
		log.Error(r, "request failed", err)
		renderErrorToast(w, r, "Request failed")
	}
}

func renderErrorToast(w http.ResponseWriter, r *http.Request, msg string) {
	toast.Toast(toast.Props{Title: msg, Type: toast.TypeError}).Render(r.Context(), w)
}

func categoryID(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
}
