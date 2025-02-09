#!/bin/sh

set -e

host="$1"
shift
cmd="$@"

until pg_isready -h "$host" -p 5432; do
  echo "Waiting for Postgres to be ready..."
  sleep 1
done

echo "Postgres is ready - executing command"
exec $cmd