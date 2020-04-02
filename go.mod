module github.com/snowzach/gogrpcapi

require (
	github.com/blendle/zapdriver v1.3.1
	github.com/docker/docker v1.13.1 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-chi/chi v4.1.0+incompatible
	github.com/go-chi/cors v1.0.1
	github.com/go-chi/render v1.0.1
	github.com/golang-migrate/migrate/v4 v4.10.0
	github.com/golang/protobuf v1.3.5
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.14.3
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.3.0
	github.com/mitchellh/mapstructure v1.2.2 // indirect
	github.com/pelletier/go-toml v1.6.0 // indirect
	github.com/rs/xid v1.2.1
	github.com/snowzach/certtools v1.0.2
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v0.0.7
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.4.0
	go.uber.org/zap v1.14.1
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
	golang.org/x/sys v0.0.0-20200331124033-c3d80250170d // indirect
	google.golang.org/genproto v0.0.0-20200401122417-09ab7b7031d2
	google.golang.org/grpc v1.28.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/ini.v1 v1.55.0 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
	gotest.tools v2.2.0+incompatible // indirect
)

go 1.13

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20190717161051-705d9623b7c1
