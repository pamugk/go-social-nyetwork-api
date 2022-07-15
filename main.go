package main

import (
	"flag"

	"github.com/pamugk/social-nyetwork-server/internal/app"
	"github.com/pamugk/social-nyetwork-server/internal/db"
	"github.com/pamugk/social-nyetwork-server/internal/infrastructure/rest"
	"github.com/pamugk/social-nyetwork-server/internal/infrastructure/rpc"
)

var (
	restPort = flag.Int("restPort", 8080, "The REST server port")
	grpcPort = flag.Int("grpcPort", 8081, "The gRPC server port")
)

func main() {
	db.Setup("postgres://postgres:postgres@localhost:5432/social_nyetwork")
	defer db.Cleanup()
	app.ConfigureValidation()
	go rpc.StartServer(grpcPort)
	rest.StartServer(restPort)
}
