package main

import (
	"embed"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/vebcreatex7/diploma_magister/cmd/handlers"
	"github.com/vebcreatex7/diploma_magister/internal/api/service"
	"github.com/vebcreatex7/diploma_magister/internal/pkg/start"
	"github.com/vebcreatex7/diploma_magister/internal/repo/postgres"
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

	db, _ := start.Postgres(cfg.Postgres)

	t, err := start.Template(templateFS, ".gohtml", true, template.FuncMap{"hasField": hasField})
	if err != nil {
		log.Fatalln(err)
	}

	r := start.Router()

	clientsRepo := postgres.NewClients(db)
	clientsService := service.NewClients(clientsRepo)
	clientsHandler := handlers.NewClients(t, clientsService)

	indexHandler := handlers.NewHome(t, clientsService)

	adminHandler := handlers.NewAdmin(t, clientsService)

	r.Mount(indexHandler.BasePrefix(), indexHandler.Routes())
	r.Mount(clientsHandler.BasePrefix(), clientsHandler.Routes())
	r.Mount(adminHandler.BasePrefix(), adminHandler.Routes())

	s := start.Server(cfg.Server, r)

	if err := s.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
