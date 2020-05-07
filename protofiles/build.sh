#!/usr/bin/env bash
protoc -I . hellouserinfo.proto --go_out=plugins=grpc:.