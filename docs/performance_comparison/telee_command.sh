#!/bin/bash

LOOP_START=1
LOOP_END=100
TIMEFORMAT=%R
TELEE_HOSTNAME=lab-cat29l-02f99-01
TELEE_COMMAND=$1

for _ in $(seq "$LOOP_START" "$LOOP_END"); \
do \
  time telee --hostname $TELEE_HOSTNAME --command "$TELEE_COMMAND" > /dev/null
  sleep 3
done

