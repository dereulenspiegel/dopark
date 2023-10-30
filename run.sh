#!/bin/bash
export DOPARK_DB_URL="postgres://dopark:dopark@localhost/dopark?sslmode=disable&binary_parameters=yes"
export DOPARK_INTERVAL="10s"
./dopark
