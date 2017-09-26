package cli

import (
	"fmt"
	//"io/ioutil"
	"github.com/codegangsta/cli"
	"github.com/linsheng9731/SLB/api"
	"github.com/linsheng9731/SLB/common"
	"github.com/linsheng9731/SLB/config"
	"github.com/linsheng9731/SLB/modules"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	CONFIG_FILENAME = "config.json"
)

func handlePanic() {
	if err := recover(); err != nil {
		log.Fatal("Server start catch panic,exit.")
	}
}

func RunServer(c *cli.Context) {
	var server *modules.Server

	if c.Bool("silence") {
		log.SetOutput(ioutil.Discard)
	}

	apiChannel := make(chan int)
	defer handlePanic()

	filename := CONFIG_FILENAME
	if c.String("filename") != "" {
		filename = c.String("filename")
	}
	log.Println("Start SSLB (Server) ")
	configuration := config.Setup(filename)

	server = modules.NewServer(configuration)
	log.Println("Prepare to run server ...")
	server.Run()

	apiInstance := api.NewAPI(apiChannel)
	apiInstance.Listen(configuration.GeneralConfig.APIAddres())
	go messageHandler(apiChannel, server)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	server.Stop()
}

func messageHandler(apiChannel chan int, server *modules.Server) {
	for {
		select {
		case msg := <-apiChannel:
			switch msg {
			case common.RELOAD:
				log.Println("Received reload message.")
				configuration := config.Setup(CONFIG_FILENAME)
				server.Stop()
				server = modules.NewServer(configuration)
				log.Println("Prepare to run server ...")
				server.Run()
			default:
				log.Println(fmt.Sprintf("Received a unrecognized message: %d", msg))
			}
		}
	}
}
