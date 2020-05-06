#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE USER ${POSTGRES_DB} PASSWORD '${POSTGRES_PASSWORD}';
    CREATE DATABASE ${POSTGRES_DB};
    GRANT ALL PRIVILEGES ON DATABASE ${POSTGRES_DB} TO ${POSTGRES_DB};

EOSQL