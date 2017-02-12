export PRECOM="/Users/leslyb/development/go/src/github.com/lbrezilien/shakur"
source $PRECOM/bash-preexec.sh
declare -a SCARYPHRASES=("git co master" "git co" "go run")
preexec() { 
  
  for PHRASE in "${SCARYPHRASES[@]}"
  do
    if [[ $1 == $PHRASE ]]
    then
      echo "Are you sure you want to do this?"
    fi
  done
 }