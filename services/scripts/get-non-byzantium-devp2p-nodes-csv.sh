#!/bin/bash

## the file
FILE=/tmp/non-byzantium-devp2p-nodes.csv

## clean it
echo "" > $FILE

## iterate over all byzantium peers
redis-cli --raw SMEMBERS devp2p-non-byzantium-peers:all | while read peerid; do
	# split the id into <hash>, <ip:port>
	echo $(echo $peerid | cut -d ':' -f1), $(echo $peerid | cut -d ':' -f2,3) >> $FILE
done