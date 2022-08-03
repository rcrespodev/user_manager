#FROM golang:1.18-alpine
#
## Set destination for COPY
#WORKDIR /app
#
## Copy the source code. Note the slash at the end, as explained in
## https://docs.docker.com/engine/reference/builder/#copy
#COPY . .
#
## Download Go modules
#RUN go mod tidy
#
## Build
#RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ./bin/main-linux cmd/web/main.go
#
## To actually open the port, runtime parameters
## must be supplied to the docker command.
#EXPOSE 8080
#
## (Optional) environment variable that our dockerised
## application can make use of. The value of environment
## variables can also be set via parameters supplied
## to the docker command on the command line.
##ENV HTTP_PORT=8081
#
## Run
##CMD [ "./bin/main-linux" ]
#CMD [ "go", "test -v ./..." ]


FROM golang:1.18-alpine as SRC

# Set destination for COPY
WORKDIR /app

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . ./

# Install dependencies
RUN set -ex; \
    apk get update; \
    apk add openssl; \
    apk update; \
    apk add --no-cache git; \
	openssl genrsa -out cert/id_rsa 4096; \
	openssl rsa -in cert/id_rsa -pubout -out cert/id_rsa.pub

# Build Go Binary
RUN set -ex; \
    CGO_ENABLED=0 GOOS=linux go build -o ./bin/main-linux cmd/web/main.go;

# Final image, no source code
FROM alpine:latest

# Install Root Ceritifcates
RUN set -ex; \
    apk update; \
    apk add --no-cache \
     ca-certificates

WORKDIR /bin/
COPY --from=src /app/bin/main-linux .

# Run Go Binary
CMD /bin/main-linux