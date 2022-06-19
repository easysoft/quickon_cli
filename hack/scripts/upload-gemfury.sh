#!/bin/bash

for filename in dist/*.deb; do
    if [[ "$filename" == *"arm"* ]]; then
        echo "Skipping $filename"
        continue
    fi
    echo "Pushing $filename to apt repository"
    # curl -F package=@$filename https://${FURY_TOKEN}@push.fury.io/qucheng/
    #curl -F package=@$filename https://${FURY_TOKEN}@push.fury.io/qucheng/
done
for filename in dist/*.rpm; do
    if [[ "$filename" == *"arm"* ]]; then
      echo "Skipping $filename"
       continue
    fi
    echo "Pushing $filename to rpm repository"
    # curl -F package=@$filename https://${FURY_TOKEN}@push.fury.io/qucheng/
    # curl -F package=@$filename https://${FURY_TOKEN}@push.fury.io/qucheng/
done
