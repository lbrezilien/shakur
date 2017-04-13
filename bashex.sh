#!/bin/bash
source bashpreex.sh
 

declare -a SCARYPHRASES=("git co master" "git co" "go run")
preexec() { 
  echo "This is the trial at work"
  for PHRASE in "${SCARYPHRASES[@]}"
  do
    if [[ $1 == $PHRASE ]]
    then
      echo "Are you sure you want to do this?"
    fi
  done
 }