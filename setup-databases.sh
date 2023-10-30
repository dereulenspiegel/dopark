#!/bin/bash
export PGPASSWORD=postgres
PSQL="psql -h localhost -p 5432 -U postgres"
$PSQL -c "DROP DATABASE IF EXISTS dopark;"
$PSQL -c "CREATE DATABASE dopark"
$PSQL -c "CREATE USER dopark WITH ENCRYPTED PASSWORD 'dopark';"
$PSQL -c "GRANT ALL PRIVILEGES ON DATABASE dopark to dopark;"
$PSQL -d dopark -c "GRANT ALL PRIVILEGES ON SCHEMA public to dopark;"
$PSQL -d dopark -c "CREATE EXTENSION IF NOT EXISTS Postgis;"