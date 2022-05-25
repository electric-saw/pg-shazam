#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRESQL_USERNAME" --dbname "$POSTGRESQL_DATABASE" <<-EOSQL
    create user test superuser password '123';
EOSQL
