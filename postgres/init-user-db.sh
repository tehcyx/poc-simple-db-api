#!/bin/bash
set -e


# postgres auto configures user, db and secret from the env variables passed.
# psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
#     CREATE USER ${POSTGRES_USER};
#     CREATE DATABASE ${POSTGRES_DB};
#     GRANT ALL PRIVILEGES ON DATABASE ${POSTGES_DB} TO ${POSTGRES_USER};
# EOSQL