#!/bin/bash

## the file
FILE=/tmp/scrapped-devp2p-nodes.csv

## clean it
echo "" > $FILE

## iterate over all scrapped peers
redis-cli --raw SMEMBERS devp2p-scrapped-peers:all | while read peerid; do
	# split the id into <hash>, <ip:port>
	split_id=$(echo $(echo $peerid | cut -d ':' -f1), $(echo $peerid | cut -d ':' -f2,3))

	## and just keep the last recorded status on each
	## If you want to have all statuses of each node, just remove the `| tail -n1` bit
	echo $split_id, $(redis-cli --raw SMEMBERS "devp2p-peerstatus:$peerid" | sort | tail -n1) >> $FILE
done