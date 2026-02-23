#!/bin/sh
set -e

mkdir -p /app/data

if [ ! -f /app/data/FGG.db ]; then
  sqlite3 /app/data/FGG.db < /app/src/dbaccess/FGG.sql
fi

exec "$@"
