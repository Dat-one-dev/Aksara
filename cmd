#!/bin/bash

case "$1" in
    run)
        ./scripts/gr
        ;;
    git)
        ./scripts/github "$2"
        ;;
    sus)
        ./scripts/sus
        ;;
    fetch)
        ./scripts/fetch
        ;;
    *)
        echo "Usage: ./cmd {run|git|sus|fetch}"
        ;;
esac
