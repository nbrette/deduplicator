#!/usr/bin/env bash

exit_on_error() {
    exit_code=$1
    last_command=${@:2}
    if [ $exit_code -ne 0 ]; then
        >&2 echo "\"${last_command}\" command failed with exit code ${exit_code}."
        exit $exit_code
    fi
}

# enable !! command completion
set -o history -o histexpand

go test --coverprofile=coverage.out ./...
exit_on_error $? !!
covered=0
total=0
while IFS='' read -r line || [[ -n "$line" ]]; do
    IFS=' ' read -r -a array <<< "$line"
    total=$(($total+${array[1]}))
    if [ "${array[2]}" = "1" ]; then
        covered=$(($covered+${array[1]}))
    fi
done < "coverage.out"
coverage=$(awk "BEGIN { pc=100*${covered}/${total}; i=int(pc); print (pc-i<0.5)?i:i+1 }")
echo "Global test coverage: $coverage%"
rm coverage.out
