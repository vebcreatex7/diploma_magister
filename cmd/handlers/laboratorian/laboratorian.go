package laboratorian

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/vebcreatex7/diploma_magister/internal/domain/service"
	"github.com/vebcreatex7/diploma_magister/pkg"
	"github.com/vebcreatex7/diploma_magister/pkg/render"
	"net/http"
)

type laboratorian struct {
	db                  *goqu.Database
	t                   *render.Template
	log                 *logrus.Logger
	clientsService      service.Clients
	equipmentService    service.Equipment
	inventoryService    service.Inventory
	accessGroupsService service.AccessGroup
	experimentService   service.Experiment
}

func NewLaboratorian(
	t *render.Template,
	log *logrus.Logger,
	clientsService service.Clients,
	equipmentService service.Equipment,
	inventoryService service.Inventory,
	accessGroupsService service.AccessGroup,
	experimentService service.Experiment,
) laboratorian {
	return laboratorian{
		t:                   t,
		log:                 log,
		clientsService:      clientsService,
		equipmentService:    equipmentService,
		inventoryService:    inventoryService,
		accessGroupsService: accessGroupsService,
		experimentService:   experimentService,
	}
}

func (h laboratorian) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(pkg.ValidateLaboratorianJWTCookies)

	r.Get("/", h.Home)
	r.Get("/home", h.Home)
	r.Get("/empty", h.Empty)

	r.Get("/experiments", h.GetExperiments)
	r.Get("/experiments/finish/{uid}", h.GetExperimentFinishByUID)
	r.Get("/experiments/{uid}", h.GetExperimentByUID)
	r.Delete("/experiments/finish/{uid}", h.FinishExperimentByUID)

	return r
}

func (h laboratorian) BasePrefix() string {
	return "/laboratorians"
}

func (h laboratorian) Home(w http.ResponseWriter, r *http.Request) {
	p := render.NewPage().
		SetTemplate("laboratorian.gohtml").
		SetPath(r.URL.Path).
		SetCode(200)

	h.t.Render(w, p)
}

func (h laboratorian) Empty(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
