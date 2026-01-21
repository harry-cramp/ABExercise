#!/bin/bash

# run go server
cd server && go build && ./abexercise &

# run web app
cd web && npm run dev
