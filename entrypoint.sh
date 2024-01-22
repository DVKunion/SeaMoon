#!/bin/sh

if [ "$SEAMOON_TOR" = "true" ] || [ "$SEAMOON_TOR" = true ]; then
  tor
fi

/app/seamoon "$@"
