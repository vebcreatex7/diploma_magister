package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/vebcreatex7/diploma_magister/internal/domain/service"
	"github.com/vebcreatex7/diploma_magister/pkg/render"
	"net/http"
)

type admin struct {
	t                   *render.Template
	log                 *logrus.Logger
	clientsService      service.Clients
	equipmentService    service.Equipment
	inventoryService    service.Inventory
	accessGroupsService service.AccessGroup
}

func NewAdmin(
	t *render.Template,
	log *logrus.Logger,
	clientsService service.Clients,
	equipmentService service.Equipment,
	inventoryService service.Inventory,
	accessGroupsService service.AccessGroup,
) admin {
	return admin{
		t:                   t,
		log:                 log,
		clientsService:      clientsService,
		equipmentService:    equipmentService,
		inventoryService:    inventoryService,
		accessGroupsService: accessGroupsService,
	}
}

func (h admin) Routes() chi.Router {

	r := chi.NewRouter()

	//r.Use(pkg.ValidateAdminJWTCookies)

	r.Get("/", h.Home)
	r.Get("/home", h.Home)
	r.Get("/empty", h.Empty)

	r.Get("/users", h.GetUsers)
	r.Get("/users/{uid}", h.GetUserByUID)
	r.Get("/users-edit/{uid}", h.GetUserEditByUID)
	r.Delete("/users/{uid}", h.DeleteUser)
	r.Put("/users/{uid}", h.EditUser)

	r.Get("/equipment", h.GetEquipment)
	r.Get("/equipment/{uid}", h.GetEquipmentByUID)
	r.Get("/equipment-edit/{uid}", h.GetEquipmentEditByUID)
	r.Delete("/equipment/{uid}", h.DeleteEquipment)
	r.Put("/equipment/{uid}", h.EditEquipment)
	r.Get("/equipment-add", h.AddEquipmentPage)
	r.Post("/equipment", h.AddEquipment)

	r.Get("/inventory", h.GetInventory)
	r.Get("/inventory/{uid}", h.GetInventoryByUID)
	r.Get("/inventory-edit/{uid}", h.GetInventoryEditByUID)
	r.Delete("/inventory/{uid}", h.DeleteInventory)
	r.Put("/inventory/{uid}", h.EditInventory)
	r.Get("/inventory-add", h.AddInventoryPage)
	r.Post("/inventory", h.AddInventory)

	r.Get("/access-groups", h.GetAccessGroups)
	r.Get("/access-groups-add", h.AddAccessGroupPage)
	r.Post("/access-groups", h.AddAccessGroup)
	r.Get("/access-groups-edit/{uid}", h.GetAccessGroupEditByUID)
	r.Put("/access-groups/{uid}", h.EditAccessGroup)
	r.Get("/access-groups/{uid}", h.GetAccessGroupByUID)
	r.Delete("/access-groups/{uid}", h.DeleteAccessGroup)

	return r
}

func (h admin) BasePrefix() string {
	return "/admin"
}

func (h admin) Home(w http.ResponseWriter, r *http.Request) {
	p := render.NewPage().
		SetTemplate("admin.gohtml").
		SetPath(r.URL.Path).
		SetCode(200)

	h.t.Render(w, p)
}

func (h admin) Empty(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
