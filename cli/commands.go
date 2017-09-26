package cli

import (
	"fmt"
	//"io/ioutil"
	"github.com/codegangsta/cli"
	"github.com/linsheng9731/SLB/api"
	"github.com/linsheng9731/SLB/common"
	"github.com/linsheng9731/SLB/config"
	"github.com/linsheng9731/SLB/server"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var serverHolder *server.LbServer

func handlePanic() {
	if err := recover(); err != nil {
		log.Fatal("LbServer start catch panic,exit.")
	}
}

func RunServer(c *cli.Context) {
	var s *server.LbServer

	if c.Bool("silence") {
		log.SetOutput(ioutil.Discard)
	}

	apiChannel := make(chan int)
	defer handlePanic()

	filename := common.CONFIG_FILENAME
	if c.String("filename") != "" {
		filename = c.String("filename")
	}
	log.Println("Start SLB (LbServer) ")
	configuration := config.Setup(filename)

	s = server.NewServer(configuration)
	serverHolder = s
	log.Println("Prepare to run server ...")
	s.Start()

	apiInstance := api.NewAPI(apiChannel)
	apiInstance.Listen(configuration.GeneralConfig.APIAddres())
	go messageHandler(apiChannel, s)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	log.Println("Prepare to stop server ...")
	serverHolder.Stop()
}

func messageHandler(apiChannel chan int, s *server.LbServer) {
	for {
		select {
		case msg := <-apiChannel:
			switch msg {
			case common.RELOAD:
				log.Println("Received reload message.")
				configuration := config.Setup(common.CONFIG_FILENAME)
				s.Stop()
				s = server.NewServer(configuration)
				log.Println("Prepare to run server ...")
				s.Start()
				serverHolder = s
			default:
				log.Println(fmt.Sprintf("Received a unrecognized message: %d", msg))
			}
		}
	}
}
