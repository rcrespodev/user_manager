# User Manager System.

[![GoReference](https://pkg.go.dev/badge/golang.org/x/tools)](https://pkg.go.dev/golang.org/x/tools)

# Contents

- [User Manager System.](#UserManagerSystem)
    - [System Design.](#system_design)
- [Features.](#Features)
    - [API Features.](#api_features)
    - [Authorization and Authentication.](#Authorization)
    - [CQRS Pattern.](#cqrs_pattern)
        - [Command.](#Command)
        - [Query.](#Query)
        - [Enqueues and Events.](#Enqueues_and_Events)
    - [Email sender.](#Email_sender)
    - [Data Layer.](#data_layer)
    - [Application Logs.](#app_logs)
- [Project Structure.](#ProjectStructure)
- [API Documentation.](#APIDocumentation)
- [Usage.](#Usage)
    - [Dependencies](#Dependencies)
    - [Testing.](#Testing)
    - [Install.](#Install)

# User Manager System. <a name="UserManagerSystem"></a>

User Manager System is Http Api designed to manage the users of
the application.
The App is built using Gin Web Framework and can be
deployed as backend in a single Docker network.
Also, the app have units and integrations tests.
And a demo for simulate correct interaction between Client and Server.

Though the data models of the api aren´t compared to real world
examples, the objective is to demonstrate different programing
concepts, such as CQRS, Event Driven Architecture, DDD, etc.

On the other hand, I'm open to working together with any front
end developer who wants to build client side to add to their
portfolio.

## System Design. <a name="system_design"></a>

Below is the high level design of the system:
![system_design_pdf](system_design.png)

# Features. <a name="Features"></a>

## API Features. <a name="api_features"></a>

The public HTTP API have the next uses cases:

- Check health status.
- Register an new user.
- Get user schema of user.
- Delete an registered user.
- Login an registered user.
- Logout user with active session.

However, all API documentation are built using open api 3.0
specification. Please see [swagger yaml](swagger.yaml) for more information.

## Authorization and Authentication. <a name="Authorization"></a>

The application core not have any Authentication. Only have the user password
for users authorization.
User password are stored in DB using hash format(bcrypt hash).
For more information, please see [user password](pkg/app/user/domain/userPassword.go)

The http layer use JWT to Authenticate the Clients request.
In every http response, the application generates a new jwt and return it
in authorization header.
For this prosperous exists the Response Gin Middleware [gin response](api/v1/handlers/ginResponse.go).
Also, when user logout, the token are invalidated in jwt repository.

The signed token use public and private rsa key.
In every docker build, two primes keys are generated in directory ./cert.
For more information, please see [jwt](pkg/app/authJwt/domain/jwt.go)

## CQRS Pattern. <a name="cqrs_pattern"></a>

### Command. <a name="Command"></a>

The app implements the pattern Command Bus to dispatch commands from
http layer to application layer. The bus it is synchronous.
Also, every command return the same type called CommandResponse.
See [Command Response](api/apiResponse.go)

### Query. <a name="Query"></a>

In the same way, the app use Query bus to communicate http layer
and service layer. Equal to command, this bus it is synchronous.
The response type of query is called QueryResponse.
See [Query Response](api/apiResponse.go)

### Enqueues and Events. <a name="Enqueues_and_Events"></a>

The unique event generated in app is UserRegistered.
And the unique consumer of this event is the email sender service.
In the infrastructure layer, I chose use Rabbit MQ as delivery chanel.
RabbitMq is an implementation on Event Bus.
See [Event Bus](pkg/kernel/cqrs/event/eventBus.go)

## Email sender. <a name="Email_sender"></a>

As described in the previous point.
The application sends a welcome email to the user when the registration is completed successfully.
I chose an SMTP Gmail server for this requisite.

## Data layer. <a name="data_layer"></a>

I choose one Relational DB and other non-sql.
The non-sql DB is for situations where we need high availability for
read model and also don´t have a high level of write operations.
In this way, the non-sql DB can scale horizontally.

### MySQL.

- [User Repository](pkg/app/user/domain/userRepository.go)
  for store the user data.
- [Sent Email Repository](pkg/app/emailSender/domain/sentEmailRepository.go)
  for store the status of emails that application sends. This repository is only for internal
  control.

### Redis.

- [Messages Repository](pkg/kernel/cqrs/returnLog/domain/message/messageRepository.go)
  All the response messages are stored in messages repository.
- [Jwt Repository](pkg/app/authJwt/domain/jwtRepository.go)
  The status of the Jwt are stored in jwt repository.

## Application logs. <a name="app_logs"></a>

All uses cases receives a pointer to [Return Log](pkg/kernel/cqrs/returnLog/domain/returnLog.go).
Return Log can store messages that can be represented success, client error or internal error.
When internal error occurs, a background task starts and write the Error
information in ./logs/ directory.

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

API documentation are built using open api 3.0
specification. Please see [swagger yaml](swagger.yaml) for more information.

On the other hand, below is the list of response messages of application. 

| Pkg   | Id  | Text  | ClientErrorType |
|-------|-----|-------|-----------------|
| user | 0   | user %v logged successful | 0               |
| user  | 1   | user %v created successful  | 0               |
| user  | 2   | user %v updated successful  | 0               |
| user  | 3   | user %v deleted successful  | 0               |
| user  | 4   | value %v is invalid as %v attribute  | 1               |
| user  | 5   | attribute %v are mandatory  | 1               |
| user  | 6   | attribute %v can´t be greater than %v characters  | 1               |
| user  | 7   | attribute %v can´t contain special characters (%v)  | 1               |
| user  | 8   | password must be contain at least one special character like %$#&  | 1               |
| user  | 9   | password must be contain at least one number  | 1               |
| user  | 10  | attribute %v can´t be smaller than %v characters  | 1               |
| user  | 11  | password must be contain at least one upper case  | 1               |
| user  | 12  | password must be contain at least one lower case  | 1               |
| user  | 13  | %v attribute dont´t must contain %v  | 1               |
| user  | 14  | user with component: %v and value: %v already exists  | 1               |
| user  | 15  | email, alias or password are not correct. Repeat the access data.  | 1               |
| user  | 16  | user logged out successful  | 1               |
| user  | 17  | none of the input values correspond to a registered user  | 1               |
| Authorization  | 0   | Unauthorized  | 2               |
| Authorization  | 1   | user is not logged  | 2               |
| email_sender  | 0   | welcome email send successful to user email %v  | 1               |

# Usage. <a name="Usage"></a>

## Dependencies.<a name="Dependencies"></a>

- Docker.
- Docker-compose.

## Testing. <a name="Testing"></a>

Integration tests run on docker containers
created automatically in handlers_test.go
Containers created for testing are not always
automatically destroyed.
That's why the script takes care of destroying
existing test containers.
To run test.sh just run:

```shell
make go_tests
```

and use:

```shell
make stop_test_services
```

in case that you need stop the test DB containers.

## Install. <a name="Install"></a>

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

The app have a demo for simulate correct interaction between Client and Server.
To run demo, first run the application, and them:

```shell
make run_demo
```