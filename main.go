package main

import (
	"flag"
	"os"

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
	db.Setup(os.Getenv("DB_URL"))
	defer db.Cleanup()
	app.ConfigureValidation()
	go rpc.StartServer(grpcPort)
	rest.StartServer(restPort)
}
