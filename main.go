package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

const (
	bashex = `
declare -a SCARYPHRASES=(%v)

shakur-preexec() {
  echo "You've entered a command that you asked shakur to watch!"
  for PHRASE in "${SCARYPHRASES[@]}"
  do
    if [[ $1 == $PHRASE ]]
    then
      echo "Are you sure you want to execute this?"
      read ANSWER
      while [  "$ANSWER" != "Y" ] && [  "$ANSWER" != "N" ]; do
        echo "please enter Y (yes) or N (no)"
        read ANSWER
      done

      if [[ "$ANSWER" == "N" ]]
      then
        echo "Ok, Shakur will be closing this terminal window if 5 seconds, you can cancel this by pressing ^C (CTRL+C)"
        echo "5....."
        sleep 1
        echo "4...."
        sleep 1
        echo "3..."
        sleep 1
        echo "2.."
        sleep 1
        echo "1 , Bye!"
        sleep 1
        exit
      else
        echo "okay, applying command"
      fi
    fi
  done
 }

 shakur-load(){
  preexec_functions+=(shakur-preexec)
}`
	bashpreex = `
# bash-preexec.sh -- Bash support for ZSH-like 'preexec' and 'precmd' functions.
# https://github.com/rcaloras/bash-preexec
#
#
# 'preexec' functions are executed before each interactive command is
# executed, with the interactive command as its argument. The 'precmd'
# function is executed before each prompt is displayed.
#
# Author: Ryan Caloras (ryan@bashhub.com)
# Forked from Original Author: Glyph Lefkowitz
#
# V0.3.1
#

# General Usage:
#
#  1. Source this file at the end of your bash profile so as not to interfere
#     with anything else that's using PROMPT_COMMAND.
#
#  2. Add any precmd or preexec functions by appending them to their arrays:
#       e.g.
#       precmd_functions+=(my_precmd_function)
#       precmd_functions+=(some_other_precmd_function)
#
#       preexec_functions+=(my_preexec_function)
#
#  3. If you have anything that's using the Debug Trap, change it to use
#     preexec. (Optional) change anything using PROMPT_COMMAND to now use
#     precmd instead.
#
#  Note: This module requires two bash features which you must not otherwise be
#  using: the "DEBUG" trap, and the "PROMPT_COMMAND" variable. prexec_and_precmd_install
#  will override these and if you override one or the other this will most likely break.

# Avoid duplicate inclusion
if [[ "$__bp_imported" == "defined" ]]; then
    return 0
fi
__bp_imported="defined"

# Should be available to each precmd and preexec
# functions, should they want it.
__bp_last_ret_value="$?"
__bp_last_argument_prev_command="$_"

# Command to set our preexec trap. It's invoked once via
# PROMPT_COMMAND and then removed.
__bp_trap_install_string="trap '__bp_preexec_invoke_exec \"\$_\"' DEBUG;"

# Remove ignorespace and or replace ignoreboth from HISTCONTROL
# so we can accurately invoke preexec with a command from our
# history even if it starts with a space.
__bp_adjust_histcontrol() {
    local histcontrol
    histcontrol="${HISTCONTROL//ignorespace}"
    # Replace ignoreboth with ignoredups
    if [[ "$histcontrol" == *"ignoreboth"* ]]; then
        histcontrol="ignoredups:${histcontrol//ignoreboth}"
    fi;
    export HISTCONTROL="$histcontrol"
}

# This variable describes whether we are currently in "interactive mode";
# i.e. whether this shell has just executed a prompt and is waiting for user
# input.  It documents whether the current command invoked by the trace hook is
# run interactively by the user; it's set immediately after the prompt hook,
# and unset as soon as the trace hook is run.
__bp_preexec_interactive_mode=""

__bp_trim_whitespace() {
    local var=$@
    var="${var#"${var%%[![:space:]]*}"}"   # remove leading whitespace characters
    var="${var%"${var##*[![:space:]]}"}"   # remove trailing whitespace characters
    echo -n "$var"
}

# This function is installed as part of the PROMPT_COMMAND;
# It sets a variable to indicate that the prompt was just displayed,
# to allow the DEBUG trap to know that the next command is likely interactive.
__bp_interactive_mode() {
    __bp_preexec_interactive_mode="on";
}

# This function is installed as part of the PROMPT_COMMAND.
# It will invoke any functions defined in the precmd_functions array.
__bp_precmd_invoke_cmd() {

    # Save the returned value from our last command
    __bp_last_ret_value="$?"

    # For every function defined in our function array. Invoke it.
    local precmd_function
    for precmd_function in "${precmd_functions[@]}"; do

        # Only execute this function if it actually exists.
        # Test existence of functions with: declare -[Ff]
        if type -t "$precmd_function" 1>/dev/null; then
            __bp_set_ret_value "$__bp_last_ret_value" "$__bp_last_argument_prev_command"
            $precmd_function
        fi
    done
}

# Sets a return value in $?. We may want to get access to the $? variable in our
# precmd functions. This is available for instance in zsh. We can simulate it in bash
# by setting the value here.
__bp_set_ret_value() {
    return $1
}

__bp_in_prompt_command() {

    local prompt_command_array
    IFS=';' read -ra prompt_command_array <<< "$PROMPT_COMMAND"

    local trimmed_arg
    trimmed_arg=$(__bp_trim_whitespace "$1")

    local command
    for command in "${prompt_command_array[@]}"; do
        local trimmed_command
        trimmed_command=$(__bp_trim_whitespace "$command")
        # Only execute each function if it actually exists.
        if [[ "$trimmed_command" == "$trimmed_arg" ]]; then
            return 0
        fi
    done

    return 1
}

# This function is installed as the DEBUG trap.  It is invoked before each
# interactive prompt display.  Its purpose is to inspect the current
# environment to attempt to detect if the current command is being invoked
# interactively, and invoke 'preexec' if so.
__bp_preexec_invoke_exec() {

    # Save the contents of $_ so that it can be restored later on.
    # https://stackoverflow.com/questions/40944532/bash-preserve-in-a-debug-trap#40944702
    __bp_last_argument_prev_command="$1"

    # Checks if the file descriptor is not standard out (i.e. '1')
    # __bp_delay_install checks if we're in test. Needed for bats to run.
    # Prevents preexec from being invoked for functions in PS1
    if [[ ! -t 1 && -z "$__bp_delay_install" ]]; then
        return
    fi

    if [[ -n "$COMP_LINE" ]]; then
        # We're in the middle of a completer. This obviously can't be
        # an interactively issued command.
        return
    fi
    if [[ -z "$__bp_preexec_interactive_mode" ]]; then
        # We're doing something related to displaying the prompt.  Let the
        # prompt set the title instead of me.
        return
    else
        # If we're in a subshell, then the prompt won't be re-displayed to put
        # us back into interactive mode, so let's not set the variable back.
        # In other words, if you have a subshell like
        #   (sleep 1; sleep 2)
        # You want to see the 'sleep 2' as a set_command_title as well.
        if [[ 0 -eq "$BASH_SUBSHELL" ]]; then
            __bp_preexec_interactive_mode=""
        fi
    fi

    if  __bp_in_prompt_command "$BASH_COMMAND"; then
        # If we're executing something inside our prompt_command then we don't
        # want to call preexec. Bash prior to 3.1 can't detect this at all :/
        __bp_preexec_interactive_mode=""
        return
    fi

    local this_command
    this_command=$(HISTTIMEFORMAT= history 1 | { read -r _ this_command; echo "$this_command"; })

    # Sanity check to make sure we have something to invoke our function with.
    if [[ -z "$this_command" ]]; then
        return
    fi

    # If none of the previous checks have returned out of this function, then
    # the command is in fact interactive and we should invoke the user's
    # preexec functions.

    # For every function defined in our function array. Invoke it.
    local preexec_function
    for preexec_function in "${preexec_functions[@]}"; do

        # Only execute each function if it actually exists.
        # Test existence of function with: declare -[fF]
        if type -t "$preexec_function" 1>/dev/null; then
            __bp_set_ret_value $__bp_last_ret_value
            $preexec_function "$this_command"
        fi
    done

    # Restore the last argument of the last executed command
    : "$__bp_last_argument_prev_command"
}

# Returns PROMPT_COMMAND with a semicolon appended
# if it doesn't already have one.
__bp_prompt_command_with_semi_colon() {

    # Trim our existing PROMPT_COMMAND
    local trimmed
    trimmed=$(__bp_trim_whitespace "$PROMPT_COMMAND")

    # Take our existing prompt command and append a semicolon to it
    # if it doesn't already have one.
    local existing_prompt_command
    if [[ -n "$trimmed" ]]; then
        existing_prompt_command=${trimmed%${trimmed##*[![:space:]]}}
        existing_prompt_command=${existing_prompt_command%;}
        existing_prompt_command=${existing_prompt_command/%/;}
    else
        existing_prompt_command=""
    fi

    echo -n "$existing_prompt_command"
}

__bp_install() {

    # Remove setting our trap from PROMPT_COMMAND
    PROMPT_COMMAND="${PROMPT_COMMAND//$__bp_trap_install_string}"

    # Remove this function from our PROMPT_COMMAND
    PROMPT_COMMAND="${PROMPT_COMMAND//__bp_install;}"

    # Exit if we already have this installed.
    if [[ "$PROMPT_COMMAND" == *"__bp_precmd_invoke_cmd"* ]]; then
        return 1;
    fi

    # Adjust our HISTCONTROL Variable if needed.
    __bp_adjust_histcontrol

    # Issue #25. Setting debug trap for subshells causes sessions to exit for
    # backgrounded subshell commands (e.g. (pwd)& ). Believe this is a bug in Bash.
    #
    # Disabling this by default. It can be enabled by setting this variable.
    if [[ -n "$__bp_enable_subshells" ]]; then

        # Set so debug trap will work be invoked in subshells.
        set -o functrace > /dev/null 2>&1
        shopt -s extdebug > /dev/null 2>&1
    fi;

    local existing_prompt_command
    existing_prompt_command=$(__bp_prompt_command_with_semi_colon)

    # Install our hooks in PROMPT_COMMAND to allow our trap to know when we've
    # actually entered something.
    PROMPT_COMMAND="__bp_precmd_invoke_cmd; ${existing_prompt_command} __bp_interactive_mode;"
    eval "$__bp_trap_install_string"

    # Add two functions to our arrays for convenience
    # of definition.
    precmd_functions+=(precmd)
    preexec_functions+=(preexec)

    # Since this is in PROMPT_COMMAND, invoke any precmd functions we have defined.
    __bp_precmd_invoke_cmd
    # Put us in interactive mode for our first command.
    __bp_interactive_mode
}

# Sets our trap and __bp_install as part of our PROMPT_COMMAND to install
# after our session has started. This allows bash-preexec to be inlucded
# at any point in our bash profile. Ideally we could set our trap inside
# __bp_install, but if a trap already exists it'll only set locally to
# the function.
__bp_install_after_session_init() {

    # Make sure this is bash that's running this and return otherwise.
    if [[ -z "$BASH_VERSION" ]]; then
        return 1;
    fi

    local existing_prompt_command
    existing_prompt_command=$(__bp_prompt_command_with_semi_colon)

    # Add our installation to be done last via our PROMPT_COMMAND. These are
    # removed by __bp_install when it's invoked so it only runs once.
    PROMPT_COMMAND="${existing_prompt_command} $__bp_trap_install_string __bp_install;"
}

# Run our install so long as we're not delaying it.
if [[ -z "$__bp_delay_install" ]]; then
    __bp_install_after_session_init
fi;
`
)

type config struct {
	Preventatives []string `yaml:"preventatives"`
}

func (c *config) getConfig() *config {

	yamlFile, err := ioutil.ReadFile("shakur.config.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func main() {
	var c config
	checkForAndDealWithFlags(c)
	c.getConfig()
	initString := ``
	for item := range c.Preventatives {
		initString = initString + fmt.Sprintf(` "%v" `, c.Preventatives[item])
	}
	loadBashFiles(initString)
}

func loadBashFiles(opt string) {
	fmt.Printf(`%s
                %s &&
                shakur-load
                `, bashpreex, fmt.Sprintf(bashex, opt))
}

func checkForAndDealWithFlags(c config) {
	watch := flag.String("watch", "", "adds command to the yaml file")
	flag.Parse()

	if *watch != "" {
		c.getConfig()
		config := ``
		for item := range c.Preventatives {
			if item != (len(c.Preventatives) - 1) {
				config = config + fmt.Sprintf(` "%v", `, c.Preventatives[item])
			} else {
				config = config + fmt.Sprintf(` "%v" `, c.Preventatives[item])
			}
		}
		var commandArray string
		if len(config) == 0 {
			commandArray = fmt.Sprintf(`["%v"]`, *watch)
		} else {
			commandArray = fmt.Sprintf(`[%v, "%v"]`, config, *watch)
		}
		configFile := fmt.Sprintf("preventatives: %v", commandArray)
		err := ioutil.WriteFile("shakur.config.yml", []byte(configFile), 0644)
		if err != nil {
			fmt.Println(err)
		}
	}
}
