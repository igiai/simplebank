#!/bin/sh

# we add below line to make sure that the script will exit immediately
# if a command returns a non-zero status
set -e

echo "run db migration"
# If Dockerfile is not run with docker-compose, DB_SOURCE env var is not defined, so we must extract env var values from the app.env
# file so that it can be used in the migrate step here
source /app/app.env
cat /app/app.env
# When ran with dokcer-composee DB_SOURCE is defined in the docker-compose.yaml and the below command will later use its value
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
# the command below means: take all parameters passed to the script and run it
# in our case it will be a binary with our app
exec "$@"