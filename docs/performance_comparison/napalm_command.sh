#!/bin/bash

LOOP_START=1
LOOP_END=100
TIMEFORMAT=%R
TELEE_HOSTNAME=lab-cat29l-02f99-01
TELEE_COMMAND=$1

for _ in $(seq "$LOOP_START" "$LOOP_END"); \
do \
  time napalm -u "$TELEE_USERNAME" -p "$TELEE_PASSWORD" -v ios --optional_args "secret='$TELEE_PRIVPASSWORD'" \
    $TELEE_HOSTNAME call cli --method-kwargs "commands=['$TELEE_COMMAND']" > /dev/null
    sleep 3
done
