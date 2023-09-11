# Build stage
FROM golang:1.21-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
# curl is necessary to download goland-migrate
RUN apk add curl
# download pre-built binary of the migrate CLI
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.18 
WORKDIR /app
# copy application executable built in the build stage
COPY --from=builder /app/main .
# copy migrate tool downloaded to the container in the build stage
COPY --from=builder app/migrate ./migrate
# We must copy file with config values to the final container because Viper
# read these values at runtime not during compile in build stage
COPY app.env .
# we have to copy this script because it is later run by entrypoint command
COPY start.sh .
# we have to copy the files with migrations to the final image so that golang-migrate can run them
COPY db/migration ./migration

# The command below doest'n actually publish the application at the specified port
# outside of the container internal networ, it only tells docker on which port to publish the app
# for the internal network usage
# it also serves as a information on which port the app should listen when it is run
EXPOSE 8080
# CMD runs a command passed as a first item of the list, next items are parameters for this command
CMD ["/app/main"]
# But when CMD is used together with ENTRYPOINT it serves as a list of additional parameters passed to a script specified in ENTRYPOINT
# so it is basically the same as ENTRYPOINT [ "/app/start.sh", "/app/main"], but dividing it to command and entrypoint gives us more flexibility
# to replace it with other command at runtime
ENTRYPOINT [ "/app/start.sh"]

# Both CMD and ENTRYPOINT will be executed everytime the docker start is executed
# Running docker run is a combination of docker create and docker start so the commands are executed 
# And also when we start an existing stopped container,they are executed