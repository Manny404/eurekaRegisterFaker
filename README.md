# eurekaRegisterFaker

Helper to register some Server to your local registry, so not all servies has to be running on your local maschine


## Build
```
go build
```

## Useage

Build your executable with go build and add all your service values to the config.json. Than you type:

```
./eurekaRegisterFaker
```

all your servies should now be visible in Eureka.

In case you use the Jhipster registry, you can query all services with this call.
´´´
http://admin:admin@localhost:8761/eureka/apps/

´´´