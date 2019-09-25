package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/ArthurHlt/go-eureka-client/eureka"
)

// RegisterEureka registriert den Service bei Eureka
func (a *App) RegisterEureka() {

	if a.Conf.EurekaURL == "" {
		fmt.Println("No EurekaURL")
		return
	}

	var closeAll []chan int

	var waitgroup sync.WaitGroup

	// Register deregister Call
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		for _, closeChannel := range closeAll {

			closeChannel <- 1
		}
	}()

	for _, service := range a.Conf.Services {

		if !service.Active {
			continue
		}

		closeChannel := make(chan int)
		closeAll = append(closeAll, closeChannel)
		waitgroup.Add(1)
		go a.registerOneService(service, closeChannel, &waitgroup)
	}

	waitgroup.Wait()
	fmt.Println("All Closed")
}

func (a *App) registerOneService(service Service, close chan int, waitgroup *sync.WaitGroup) {

	defer waitgroup.Done()

	client := eureka.NewClient([]string{
		a.Conf.EurekaURL, //From a spring boot based eureka server
		// add others servers here
	})

	fmt.Println("Register " + service.Appname)

	port, err := strconv.Atoi(service.Port)
	if err != nil {
		panic("Port not valid int " + service.Port)
	}

	instance := eureka.NewInstanceInfo(service.Host, service.Appname, service.Host, port, 30, false) //Create a new instance to register
	instance.Metadata = &eureka.MetaData{
		Map: make(map[string]string),
	}
	instance.VipAddress = service.Appname
	instance.SecureVipAddress = service.Appname
	instance.HomePageUrl = service.Host + ":" + service.Port + "/health"
	instance.HealthCheckUrl = service.Host + ":" + service.Port + "/health"
	instance.StatusPageUrl = service.Host + ":" + service.Port + "/info"
	instance.Metadata.Map["management.port"] = service.Port
	instance.Metadata.Map["name"] = service.Appname //add metadata for example
	instance.Metadata.Map["profile"] = "dev"
	instance.Metadata.Map["zone"] = "primary"
	instance.Metadata.Map["version"] = "v1"

	instance.InstanceID = service.Appname + ":" + instance.HostName + "" + strconv.FormatInt(rand.Int63(), 10)

	client.RegisterInstance(service.Appname, instance) // Register new instance in your eureka(s)
	//applications, _ := client.GetApplications()           // Retrieves all applications from eureka server(s)
	//client.GetApplication(instance.App)                   // retrieve the application "test"
	//client.GetInstance(instance.App, instance.HostName)   // retrieve the instance from "test.com" inside "test"" app
	// say to eureka that your app is alive (here you must send heartbeat before 30 sec)

	// HeartBeat
Loop:
	for {
		select {

		case <-close:
			// Register deregister Call
			fmt.Println("Deregister Service " + instance.InstanceID)
			client.UnregisterInstance(instance.App, instance.InstanceID)
			break Loop
		case <-time.After(30 * time.Second):

			err := client.SendHeartbeat(instance.App, instance.InstanceID)
			if err != nil {
				waitgroup.Add(1)
				go a.registerOneService(service, close, waitgroup)
				break Loop
			}
		}

	}
}
