#!/usr/bin/env bash
cd chat
go build -o chat

./chat -addr=":8080" \
    -clientId=$CHATTER_BOX_CLI_ID \
    -clientSecret=$CHATTER_BOX_CLI_SECRET \
    -securityKey=$CHATTER_BOX_SECRET_KEY
