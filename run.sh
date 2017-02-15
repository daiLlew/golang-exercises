#!/usr/bin/env bash
cd chat
go build -o chat
./chat -addr=":3000"
