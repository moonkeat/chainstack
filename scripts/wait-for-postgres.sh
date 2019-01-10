#!/bin/sh
# wait-for-postgres.sh

set -e

host="$1"
shift
cmd="$@"

until PGPASSWORD=$POSTGRES_PASSWORD psql -h "$host" -U "postgres" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
sh -c "psql -h $host -U \"postgres\" -tc \"SELECT 1 FROM pg_database WHERE datname = 'chainstack'\" | grep -q 1 || psql -h $host  -U \"postgres\" -c \"CREATE DATABASE chainstack\""
sh -c "./db/goose -dir ./db/migrations postgres $DB_CONNSTRING up"
sh -c "./createuser -email admin@admin.com -password password -admin true"
exec "./main"
