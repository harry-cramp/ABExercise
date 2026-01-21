#!/bin/bash

echo "Starting New Phones server..."
/app/server/abexercise &

echo "Starting React webapp..."
cd /app/web
npm run dev -- --host

wait
