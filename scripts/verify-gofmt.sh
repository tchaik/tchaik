#!/bin/bash

lint=$(gofmt -l -e $*)
if [[ -z "$lint" ]]; then
    exit 0;
else
    echo "Formatting errors in the following files:"
    echo $lint
    exit 1;
fi
