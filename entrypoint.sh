#!/bin/sh

if [ "$TOR" = "true" ]; then
  tor
fi

/app/seamoon "$@"