#!/bin/bash

function killit {
	echo "Killing $1"
	killall -9 $1
}

killit core-metadata
killit core-command
killit core-data
killit export-distro
killit export-client

