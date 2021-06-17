#!/bin/bash

EXECUTION_TYPE=$1

if [ "$EXECUTION_TYPE" == "error" ]; then
   echo "It was finished with a error" >> /dev/stderr
   exit 1
fi 

echo "It was well executed"