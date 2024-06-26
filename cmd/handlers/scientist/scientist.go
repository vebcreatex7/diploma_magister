package scientist

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/vebcreatex7/diploma_magister/internal/domain/service"
	"github.com/vebcreatex7/diploma_magister/pkg"
	"github.com/vebcreatex7/diploma_magister/pkg/render"
	"net/http"
)

type scientist struct {
	db                  *goqu.Database
	t                   *render.Template
	log                 *logrus.Logger
	clientsService      service.Clients
	equipmentService    service.Equipment
	inventoryService    service.Inventory
	accessGroupsService service.AccessGroup
	experimentService   service.Experiment
}

func NewScientist(
	t *render.Template,
	log *logrus.Logger,
	clientsService service.Clients,
	equipmentService service.Equipment,
	inventoryService service.Inventory,
	accessGroupsService service.AccessGroup,
	experimentService service.Experiment,
) scientist {
	return scientist{
		t:                   t,
		log:                 log,
		clientsService:      clientsService,
		equipmentService:    equipmentService,
		inventoryService:    inventoryService,
		accessGroupsService: accessGroupsService,
		experimentService:   experimentService,
	}
}

func (h scientist) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(pkg.ValidateScientistJWTCookies)

	r.Get("/", h.Home)
	r.Get("/home", h.Home)
	r.Get("/empty", h.Empty)

	r.Get("/access-groups", h.GetMyAccessGroups)

	r.Get("/equipment", h.GetAllMyEquipment)
	r.Get("/equipment/schedule", h.GetEquipmentSchedule)
	r.Get("/equipment/schedule/empty", h.GetEquipmentScheduleEmpty)

	r.Get("/inventory", h.GetAllMyInventory)

	r.Get("/experiments/add", h.AddExperimentPage)
	r.Get("/experiments/equipment/add", h.ExperimentEquipmentFormAdd)
	r.Get("/experiments/inventory/add", h.ExperimentInventoryFormAdd)
	r.Post("/experiments/add", h.AddExperiment)
	r.Get("/experiments", h.GetMyExperiments)
	r.Delete("/experiments/{uid}", h.DeleteExperimentByUID)

	return r
}

func (h scientist) BasePrefix() string {
	return "/scientists"
}

func (h scientist) Home(w http.ResponseWriter, r *http.Request) {
	p := render.NewPage().
		SetTemplate("scientist.gohtml").
		SetPath(r.URL.Path).
		SetCode(200)

	h.t.Render(w, p)
}

func (h scientist) Empty(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
