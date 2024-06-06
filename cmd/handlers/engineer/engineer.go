package engineer

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/vebcreatex7/diploma_magister/internal/domain/service"
	"github.com/vebcreatex7/diploma_magister/pkg"
	"github.com/vebcreatex7/diploma_magister/pkg/render"
	"net/http"
)

type engineer struct {
	db                  *goqu.Database
	t                   *render.Template
	log                 *logrus.Logger
	clientsService      service.Clients
	equipmentService    service.Equipment
	inventoryService    service.Inventory
	accessGroupsService service.AccessGroup
	experimentService   service.Experiment
	maintainceService   service.Maintaince
}

func NewEngineer(
	t *render.Template,
	log *logrus.Logger,
	clientsService service.Clients,
	equipmentService service.Equipment,
	inventoryService service.Inventory,
	accessGroupsService service.AccessGroup,
	experimentService service.Experiment,
	maintainceService service.Maintaince,
) engineer {
	return engineer{
		t:                   t,
		log:                 log,
		clientsService:      clientsService,
		equipmentService:    equipmentService,
		inventoryService:    inventoryService,
		accessGroupsService: accessGroupsService,
		experimentService:   experimentService,
		maintainceService:   maintainceService,
	}
}

func (h engineer) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(pkg.ValidateEngineerJWTCookies)

	r.Get("/", h.Home)
	r.Get("/home", h.Home)
	r.Get("/empty", h.Empty)

	r.Get("/equipment", h.GetAllEquipment)
	r.Get("/equipment/schedule", h.GetEquipmentSchedule)
	r.Get("/equipment/schedule/empty", h.GetEquipmentScheduleEmpty)

	r.Get("/inventory", h.GetAllInventory)
	r.Get("/inventory/{uid}", h.GetInventoryByUID)
	r.Get("/inventory-edit/{uid}", h.GetInventoryEditByUID)
	r.Delete("/inventory/{uid}", h.DeleteInventoryByUID)
	r.Put("/inventory/{uid}", h.EditInventoryByUID)
	r.Get("/inventory-add", h.AddInventoryPage)
	r.Post("/inventory", h.AddInventory)

	r.Get("/maintaince/add", h.AddMaintaincePage)
	r.Get("/maintaince/equipment/add", h.MaintainceEquipmentFormAdd)
	r.Post("/maintaince/add", h.AddMaintaince)
	r.Get("/maintaince", h.GetMyMaintaince)
	r.Delete("/maintaince/{uid}", h.DeleteMaintainceByUID)

	return r
}

func (h engineer) BasePrefix() string {
	return "/engineers"
}

func (h engineer) Home(w http.ResponseWriter, r *http.Request) {
	p := render.NewPage().
		SetTemplate("engineer.gohtml").
		SetPath(r.URL.Path).
		SetCode(200)

	h.t.Render(w, p)
}

func (h engineer) Empty(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
