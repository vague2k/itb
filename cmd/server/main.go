package main

import (
	"log"
	"os"

	"itb.ihatedoing.work/config"
	"itb.ihatedoing.work/server"
	"itb.ihatedoing.work/server/database"
	"itb.ihatedoing.work/server/handlers"
)

func main() {
	command := "serve"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	cfg := config.Init()

	switch command {
	case "serve":
		h := handlers.NewHandler(cfg)

		s := server.NewServer(cfg)
		s.SetupAssetsRoutes()
		s.RegisterRoutes(h)

		if err := s.Run(); err != nil {
			log.Fatalf("server stopped: %v", err)
		}
	case "migrate":
		sub := "up"
		if len(os.Args) > 2 {
			sub = os.Args[2]
		}

		var err error
		switch sub {
		case "up":
			err = database.MigrateUp(cfg.UnderlyingDB)
		case "down":
			err = database.MigrateDown(cfg.UnderlyingDB)
		case "status":
			err = database.MigrateStatus(cfg.UnderlyingDB)
		default:
			log.Fatalf("unknown migrate command %q (expected up, down, or status)", sub)
		}
		if err != nil {
			log.Fatalf("migrate %s: %v", sub, err)
		}
	default:
		log.Fatalf("unknown command %q (expected serve or migrate)", command)
	}
}
