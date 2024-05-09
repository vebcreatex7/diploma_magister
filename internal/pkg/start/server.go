package start

import (
	"fmt"
	"github.com/vebcreatex7/diploma_magister/internal/config"
	"net/http"
)

func Server(cfg config.Server, h http.Handler) *http.Server {
	s := http.Server{Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), Handler: h}

	return &s
}
