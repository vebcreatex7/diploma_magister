package main

import (
	"embed"
	"fmt"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
	"github.com/vebcreatex7/diploma_magister/cmd/handlers"
	"github.com/vebcreatex7/diploma_magister/cmd/handlers/admin"
	"github.com/vebcreatex7/diploma_magister/cmd/handlers/engineer"
	"github.com/vebcreatex7/diploma_magister/cmd/handlers/laboratorian"
	"github.com/vebcreatex7/diploma_magister/cmd/handlers/scientist"
	"github.com/vebcreatex7/diploma_magister/internal/api/service"
	"github.com/vebcreatex7/diploma_magister/internal/repo/postgres"
	"github.com/vebcreatex7/diploma_magister/pkg/mailer"
	"github.com/vebcreatex7/diploma_magister/pkg/render"
	start2 "github.com/vebcreatex7/diploma_magister/pkg/start"
	"html/template"
	"os"
	"reflect"
)

var (
	//go:embed templates
	templateFS embed.FS
	//parsed templates
	html *template.Template
)

func hasField(v interface{}, name string) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return false
	}
	return rv.FieldByName(name).IsValid()
}
func main() {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	mailPswd, exists := os.LookupEnv("MAIL_RU_PASSWORD")
	fmt.Println(exists)

	m := mailer.NewMailer("kirik.vahramyan@mail.ru", mailPswd, "smtp.mail.ru", 465)

	cfg, err := initConfig("cmd/config/conf.yaml")
	if err != nil {
		log.Fatalf("initing config: %s", err)
	}

	db, _ := start2.Postgres(cfg.Postgres)

	templ, err := start2.Template(templateFS, ".gohtml", true, template.FuncMap{"hasField": hasField})
	if err != nil {
		log.Fatalln(err)
	}

	t := render.NewTemplate(templ, log)

	r := start2.Router()

	clientsRepo := postgres.NewClients(db)
	clientsService := service.NewClients(clientsRepo, db)
	equipmentRepo := postgres.NewEquipment(db)

	inventoryRepo := postgres.NewInventory(db)

	accessGroupRepo := postgres.NewAccessGroups(db)
	accessGroupService := service.NewAccessGroup(
		accessGroupRepo,
		clientsRepo,
		equipmentRepo,
		inventoryRepo,
	)
	equipmentService := service.NewEquipment(equipmentRepo, clientsRepo, accessGroupService, db)
	inventoryService := service.NewInventory(inventoryRepo, accessGroupService, db)

	experimentService := service.NewExperiment(
		db,
		clientsRepo,
		equipmentService,
	)
	maintainceService := service.NewMaintaince(
		db,
		clientsRepo,
		equipmentService,
		m,
	)

	indexHandler := handlers.NewHome(templ, log, t, clientsService)

	adminHandler := admin.NewAdmin(
		t,
		log,
		clientsService,
		equipmentService,
		inventoryService,
		accessGroupService,
		experimentService,
		maintainceService,
	)

	scientistHandler := scientist.NewScientist(
		t,
		log,
		clientsService,
		equipmentService,
		inventoryService,
		accessGroupService,
		experimentService,
	)

	engineerHandler := engineer.NewEngineer(
		t,
		log,
		clientsService,
		equipmentService,
		inventoryService,
		accessGroupService,
		experimentService,
		maintainceService,
	)

	laboratorianHandler := laboratorian.NewLaboratorian(
		t,
		log,
		clientsService,
		equipmentService,
		inventoryService,
		accessGroupService,
		experimentService,
	)

	r.Mount(indexHandler.BasePrefix(), indexHandler.Routes())
	r.Mount(adminHandler.BasePrefix(), adminHandler.Routes())
	r.Mount(scientistHandler.BasePrefix(), scientistHandler.Routes())
	r.Mount(engineerHandler.BasePrefix(), engineerHandler.Routes())
	r.Mount(laboratorianHandler.BasePrefix(), laboratorianHandler.Routes())

	s := start2.Server(cfg.Server, r)

	if err := s.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
