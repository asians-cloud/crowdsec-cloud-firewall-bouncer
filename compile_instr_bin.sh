#!/usr/bin/env bash
go test . -tags testbincover -coverpkg=./... -c -o crowdsec-cloud-firewall-bouncer_instr-bin -ldflags="-X github.com/asians-cloud/crowdsec-cloud-firewall-bouncer.isTest=true"