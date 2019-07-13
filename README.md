# eurekaRegisterFaker

Helper to register some Server to your local registry, so not all servies has to be running on your local maschine.


## Build
```
go build
```

## Useage

Build your executable with go build and add all your service values to the config.json. Than you can type:

```
./eurekaRegisterFaker
```

All your servies from the json file should now be visible in Eureka.

In case you use the Jhipster registry, you can query all services with this call.
```
http://admin:admin@localhost:8761/eureka/apps/
```
For example, if you exit the program with ctrg + c, all services will be deregistered from Eureka.
