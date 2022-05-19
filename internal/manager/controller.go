package manager

import (
	"net/http"
	"time"

	"github.com/durandj/ley/internal/manager/network"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Controller handles HTTP requests as well as setting up any required
// middleware across all endpoints.
type Controller struct {
	router            chi.Router
	networkController *network.Controller
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
	router.Use(middleware.CleanPath)
	router.Use(middleware.Heartbeat("/healthcheck"))

	networkController := &network.Controller{
		NetworkService: network.NewService(),
	}
	router.Route("/network", func(subrouter chi.Router) {
		networkController.RegisterRoutes(subrouter)
	})

	return &Controller{
		router:            router,
		networkController: networkController,
	}
}

func (controller *Controller) ServeHTTP(
	response http.ResponseWriter,
	request *http.Request,
) {
	controller.router.ServeHTTP(response, request)
}

var _ http.Handler = (*Controller)(nil)
