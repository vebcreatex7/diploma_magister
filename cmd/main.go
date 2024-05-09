package main

import (
	"embed"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/vebcreatex7/diploma_magister/cmd/handlers"
	"github.com/vebcreatex7/diploma_magister/internal/api/service"
	"github.com/vebcreatex7/diploma_magister/internal/repo/postgres"
	start2 "github.com/vebcreatex7/diploma_magister/pkg/start"
	"html/template"
	"log"
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
	cfg, err := initConfig("cmd/config/conf.yaml")
	if err != nil {
		log.Fatalf("initing config: %s", err)
	}

	db, _ := start2.Postgres(cfg.Postgres)

	t, err := start2.Template(templateFS, ".gohtml", true, template.FuncMap{"hasField": hasField})
	if err != nil {
		log.Fatalln(err)
	}

	r := start2.Router()

	clientsRepo := postgres.NewClients(db)
	clientsService := service.NewClients(clientsRepo)

	equipmentRepo := postgres.NewEquipment(db)
	equipmentServce := service.NewEquipment(equipmentRepo)

	indexHandler := handlers.NewHome(t, clientsService)

	adminHandler := handlers.NewAdmin(t, clientsService, equipmentServce)

	r.Mount(indexHandler.BasePrefix(), indexHandler.Routes())
	r.Mount(adminHandler.BasePrefix(), adminHandler.Routes())

	s := start2.Server(cfg.Server, r)

	if err := s.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
