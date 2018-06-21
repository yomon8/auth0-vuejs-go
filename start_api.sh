#!/usr/bin/env bash
docker build -t auth0-golang-api ./api
docker run -p 50000:50000 --name auth0-golang-api --rm -d auth0-golang-api