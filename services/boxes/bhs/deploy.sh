#!/bin/bash

################################################################################
#
# Deploy
#
#
# Run this script in your local machine.
#
# * PLEASE MODIFY the variable REMOTE_MACHINE
#   to your remote machine ssh access minutiae.
#
# * It copies into your remote machine's directory in /home the files
#   you need from this directory.
#
# * Runs remotely the `prepare_box` script (download basic elements).
#
# * Executes remotely the `run_box` script (starting the docker-compose spec).
#
################################################################################

# Variables

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

REMOTE_MACHINE="core@monkey.musteka.la"

REMOTE_HOME_DIR="/home/core/bhs"

# Actual commands: Copy of files

echo -e "\nCopying files to the remote machine\n"

ssh $REMOTE_MACHINE mkdir -p $REMOTE_HOME_DIR

scp $CURRENT_DIR/prepare_box $REMOTE_MACHINE:$REMOTE_HOME_DIR/prepare_box
scp $CURRENT_DIR/run_box $REMOTE_MACHINE:$REMOTE_HOME_DIR/run_box
scp $CURRENT_DIR/stop_box $REMOTE_MACHINE:$REMOTE_HOME_DIR/stop_box
scp $CURRENT_DIR/Dockerfile $REMOTE_MACHINE:$REMOTE_HOME_DIR/Dockerfile
scp $CURRENT_DIR/docker-compose.yml $REMOTE_MACHINE:$REMOTE_HOME_DIR/docker-compose.yml

# Actual commands: Execution

echo -e "\nExecuting preparation of box\n"

ssh $REMOTE_MACHINE $REMOTE_HOME_DIR/prepare_box

# TODO: stop_box and run_box
