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
	"strconv"
	"syscall"
)

var serverHolder *server.LbServer
var apiChannel chan int
var apiInstance *api.API

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

	apiChannel = make(chan int)
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

	apiInstance = api.NewAPI(serverHolder, apiChannel)
	apiInstance.Listen(configuration.GeneralConfig.APIAddres())
	go messageHandler(apiChannel, s)
	listenSignal()
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	f, err := os.Create("slb.pid")
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteString(strconv.Itoa(os.Getpid()))
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	log.Println(<-ch)
	log.Println("Prepare to stop server ...")
	serverHolder.Stop()
}

func listenSignal() {
	s := make(chan os.Signal)
	signal.Notify(s, syscall.SIGHUP)
	go func() {
		for {
			<-s
			go func() { apiChannel <- common.RELOAD }()
			log.Println("Receive hot reload signal.")
		}
	}()
}

func HotReload(c *cli.Context) {
	f, err := os.Open("./slb.pid")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var pid int
	_, err = fmt.Fscanf(f, "%d\n", &pid)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Read pid from slb.pid : %d ", pid)
	syscall.Kill(int(pid), syscall.SIGHUP) // reload
	log.Println("Send reload signal to lb server.")

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
				apiInstance.Serer = s
			default:
				log.Println(fmt.Sprintf("Received a unrecognized message: %d", msg))
			}
		}
	}
}
