#!/bin/bash

kubectl scale deploy market-cne-market-api -n quickon-system --replicas=0

kubectl apply -f /opt/quickon/deploy/cne-config.yaml

kubectl apply -f /opt/quickon/deploy/import-db.yaml

kubectl scale deploy market-cne-market-api -n quickon-system --replicas=1

helm upgrade -i market install/cne-market-api -n quickon-system
