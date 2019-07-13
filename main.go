package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type App struct {
	Conf *Configuration
}

func main() {

	fmt.Println("Starting registerFaker")
	fmt.Println("Reading Config")

	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error Read Config:", err)
	} else {

		a := App{}
		a.Conf = &configuration

		go sayRunning()

		a.RegisterEureka()
	}
}

func sayRunning() {

	fmt.Println("Registryfaker running :)")
}
