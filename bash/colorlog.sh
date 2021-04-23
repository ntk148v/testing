#!/bin/bash

# tput colors
BLACK=0
RED=1
GREEN=2
BLUE=4
WHITE=7

# make text have cool color
function log_info {
  tput setaf $BLUE; echo -e $1
}

function log_success {
  tput setaf $GREEN
  echo -e "✔ " $1
}

function log_error {
  tput setaf $RED
  echo -e "✘" $1
}

# How to use
# source colorlog.sh
# a_command && log_success "OK" || log_error "Failed"
# for example:
# cat nonexistfile && log_success "OK" || log_error "Failed"
# cat existfile && log_success "OK" || log_error "Failed"
