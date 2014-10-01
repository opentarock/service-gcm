package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"

	"github.com/opentarock/service-api/go/proto_gcm"
	nservice "github.com/opentarock/service-api/go/service"
	"github.com/opentarock/service-gcm/gcm"
	"github.com/opentarock/service-gcm/service"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

var gcmapikey = flag.String("gcmapikey", "", "Google Cloud Messaging api key")

var dryrun = flag.Bool("dryrun", false, "Test requests without actually sending the messages")

func main() {
	flag.Parse()
	// profiliing related flag
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if *gcmapikey == "" {
		log.Fatalf("Google Cloud Messaging api key is required")
	}

	log.SetFlags(log.Ldate | log.Lmicroseconds)

	log.Println("Starting gcm service ...")

	gcmService := nservice.NewRepService(nservice.MakeServiceBindAddress(nservice.GcmServiceDefaultPort))

	gcmSender := gcm.NewRetrySender(*gcmapikey)
	gcmSender.DryRun = *dryrun
	if gcmSender.DryRun {
		log.Println("Dry run mode enabled")
	}

	handlers := service.NewGcmServiceHandlers(gcmSender)

	gcmService.AddHandler(proto_gcm.SendMessageRequestMessage, handlers.SendMessageHandler())

	err := gcmService.Start()
	if err != nil {
		log.Fatalf("Error starting gcm service: %s", err)
	}
	defer gcmService.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	sig := <-c
	log.Printf("Interrupted by %s", sig)
}
