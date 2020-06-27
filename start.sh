#!/bin/sh

set -e

exec /app/server server --configpath ./config/global.json
