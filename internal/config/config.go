package config

type API struct {
	Postgres Postgres `yaml:"postgres"`
	Server   Server   `yaml:"server"`
}
