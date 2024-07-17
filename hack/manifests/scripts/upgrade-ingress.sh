#!/bin/bash

z helm repo update
z helm upgrade -i ingress install/nginx-ingress-controller -n quickon-system
