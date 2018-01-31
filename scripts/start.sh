#!/bin/bash

function launch {
	cd `dirname $1`
	./`basename $1` &
	cd -
	echo "$1 has been launched"
	sleep 5
}

set +e


launch cmd/core-metadata/core-metadata
launch cmd/core-command/core-command
launch cmd/core-data/core-data
launch cmd/export-distro/export-distro
launch cmd/export-client/export-client 

