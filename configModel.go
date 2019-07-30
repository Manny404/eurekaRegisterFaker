package main

type Configuration struct {
	EurekaURL string
	Services  []Service
}

type Service struct {
	Active  bool
	Appname string
	Host    string
	Port    string
}
