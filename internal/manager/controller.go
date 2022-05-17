package manager

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Controller handles HTTP requests as well as setting up any required
// middleware across all endpoints.
type Controller struct {
	router chi.Router
}

// NewController sets up a new controller and the required middleware.
func NewController() *Controller {
	router := chi.NewRouter()

	router.Use(middleware.RealIP)
	// TODO: switched to shared logger
	// https://gist.github.com/ndrewnee/6187a01427b9203b9f11ca5864b8a60d
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Timeout(time.Minute))
	router.Use(middleware.AllowContentType("application/json"))
	router.Use(middleware.ContentCharset("UTF-8"))
	router.Use(middleware.CleanPath)
	router.Use(middleware.Heartbeat("/healthcheck"))

	return &Controller{
		router: router,
	}
}

func (controller *Controller) ServeHTTP(
	response http.ResponseWriter,
	request *http.Request,
) {
}

var _ http.Handler = (*Controller)(nil)
