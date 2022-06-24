#!/bin/zsh

scp -i ~/.ssh/spacedev.pem config.yaml ubuntu@34.204.77.196:/home/ubuntu/release

cd "/Users/evantu/Documents/GoProjects/src/avchem-server/cmd/avchem"

GOOS=linux GOARC=amd64 go build -o avchem

scp -i ~/.ssh/spacedev.pem avchem ubuntu@34.204.77.196:/home/ubuntu/release

ssh -i ~/.ssh/spacedev.pem ubuntu@34.204.77.196