# User Manager System.

[![GoReference](https://pkg.go.dev/badge/golang.org/x/tools)](https://pkg.go.dev/golang.org/x/tools)

# Contents
- [User Manager System.](#UserManagerSystem)
- [Features.](#Features)
- [Project Structure.](#ProjectStructure)
- [API Documentation.](#APIDocumentation)
- [Usage.](#Usage)
  - [Testing.](#Testing)
  - [Install.](#Install)

# User Manager System. <a name="UserManagerSystem"></a>
User Manager System is Http Api designed to manage the users of
the application.
The App is built using Gin Web Framework and can be
deployed as microservices in a single Docker network.

Though the data models of the api aren´t compared to real world
examples, the objective is to demostrate different programing
concepts, such as CQRS, Event Driven Architecture, DDD, etc.

On the other hand, I'm open to working together with any front
end developer who wants to build client side to add to their
portfolio.

# Features. <a name="Features"></a>
* Register, update and delete users accounts.
* Json Web Tokens to handle Authentications.
* CQRS Pattern.
* Email sender based on domain events.
* Swagger documentation (Open API 3.0).
* Custom response codes.

# Project Structure. <a name="ProjectStructure"></a>
The __pkg__ directory contains a high level modules:
- kernel: The heart of the system. All dependencies
are handling in this pkg.
- server: Depends on kernel pkg and have the
responsibility of start the Http Serve.
- app: Contains all business logic. Is the heart
of the Application.

In the other hand, the __test__ directory contains:
- Unit: Unit testing. Without dependencies.
- Integration: This pkg is for end-to-end testing.
Up the dependencies using Docker test and
Mock Gin Http Server. 

# API Documentation. <a name="APIDocumentation"></a>
All the Api documentation are built using open api 3.0
specification. Please see [swagger yaml](swagger.yaml).

# Usage. <a name=""></a>
## Dependencies.
- Docker.
- Docker-compose.
- Bash script.

## Testing. <a name=""></a>

Integration tests run on docker containers
created automatically in every main_test.go
Containers created for testing are not always
automatically destroyed.
That's why the script takes care of destroying
existing test containers.
To run test.sh just run:
```shell
make run_tests
```
If you not have bin bash, use:
```shell
make go_test
```
and use:
```shell
make stop_test_services
```
to stop the test containers.

## Install. <a name=""></a>
Deploy and run up in background mode:
```shell
make run
```

Deploy and run up online:
```shell
make run_online
```

If some redis port or mysql port are allocated,
try with:
```shell
make stop_services
```