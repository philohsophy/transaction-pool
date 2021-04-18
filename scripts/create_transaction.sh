#!/bin/bash
curl -X POST -H "Content-Type: application/json" \
-d @transaction.json \
http://localhost:8010/transactions