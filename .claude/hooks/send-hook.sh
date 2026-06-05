#!/bin/bash
hook=$(cat)

curl -s -X POST http://localhost:10987/hooks \
  -H "Content-Type: application/json" \
  -d "$hook"
