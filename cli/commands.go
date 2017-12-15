package cli

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/linsheng9731/slb/api"
	"github.com/linsheng9731/slb/common"
	"github.com/linsheng9731/slb/config"
	"github.com/linsheng9731/slb/logger"
	"github.com/linsheng9731/slb/modules"
	"github.com/linsheng9731/slb/server"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var serverHolder *server.LbServer
var apiChannel chan int
var apiInstance *api.API
var lg *logger.Logger

func handlePanic() {
	if err := recover(); err != nil {
		lg.Error(err)
		lg.Fatal("LbServer start catch panic,exit.")
	}
}

func RunServer(c *cli.Context) {
	var s *server.LbServer

	apiChannel = make(chan int)
	defer handlePanic()

	lg = logger.NewLogger("./server.log", config.GlobalConfig)

	lg.Info("Start SLB (LbServer) ")
	m := modules.NewMetrics()
	m.IntervalTask()
	s = server.NewServer(config.GlobalConfig, m)
	serverHolder = s
	lg.Info("Prepare to run server ...")
	s.Run()

	apiInstance = api.NewAPI(serverHolder, apiChannel)
	apiInstance.Listen(config.GlobalConfig.Addres())
	//go messageHandler(apiChannel, s)
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
	lg.Info("Prepare to stop server ...")
	serverHolder.Stop()
}

func listenSignal() {
	s := make(chan os.Signal)
	signal.Notify(s, syscall.SIGHUP)
	go func() {
		for {
			<-s
			go func() { apiChannel <- common.RELOAD }()
			lg.Info("Receive hot reload signal.")
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
	lg.Info("Read pid from slb.pid : %d ", pid)
	syscall.Kill(int(pid), syscall.SIGHUP) // reload
	lg.Info("Send reload signal to lb server.")

}

func StopCommand(c *cli.Context) {
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
	lg.Info("Read pid from slb.pid : %d ", pid)
	syscall.Kill(int(pid), syscall.SIGINT) // interrupt
	lg.Info("Send interrupt signal to lb server.")
}

func messageHandler(apiChannel chan int, s *server.LbServer) {
	//for {
	//	select {
	//	case msg := <-apiChannel:
	//		switch msg {
	//		case common.RELOAD:
	//			lg.Info("Received reload message.")
	//			configuration := config.Setup(common.CONFIG_FILENAME)
	//			s.Stop()
	//			s = server.NewServer(configuration)
	//			lg.Info("Prepare to run server ...")
	//			s.Run()
	//			serverHolder = s
	//			apiInstance.Serer = s
	//		default:
	//			lg.Info(fmt.Sprintf("Received a unrecognized message: %d", msg))
	//		}
	//	}
	//}
}
