package main

type Configuration struct {
	EurekaURL string
	Services  []Service
}

type Service struct {
	Appname string
	Host    string
	Port    string
}
