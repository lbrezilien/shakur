#!/bin/bash
 

declare -a SCARYPHRASES=("git co master" "git co" "go run")
loadpreex(){
  source bashpreex.sh

}
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