#!/bin/bash

function launch {
	cd `dirname $1`
	./`basename $1` &
	cd -
	echo "$1 has been launched"
	sleep 5
}

set +e

pathToMongodb=$(which mongod)
set -e
if [ ! -x "$pathToMongodb" ] ; then
    echo -e "You need to install mongodb."
    echo -e "If you are using ubuntu\n\tsudo apt install mongodb"
    echo -e "For arm64 use the packaged mongo"
    exit 1
fi

# Mongo
mongod --config config/mongodb.conf &

sleep 8
mongo --host=127.0.0.1 config/init_mongo.js

launch cmd/core-metadata/core-metadata
launch cmd/core-command/core-command
launch cmd/core-data/core-data
launch cmd/export-distro/export-distro
launch cmd/export-client/export-client 

