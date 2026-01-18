#!/bin/sh
set -e

mkdir -p /app/data

if [ ! -f /app/data/FGG.db ]; then
  echo "Initializing database..."
  sqlite3 /app/data/FGG.db < /app/db_access/FGG.sql
fi

exec "$@"
