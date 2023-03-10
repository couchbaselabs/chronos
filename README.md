# Chronos

## FTS monitoring/diagnostic tool

The chronos tool is designed to get live stats from a couchbase server and display them in a terminal user interface.

It takes the node ip information and the stats that needs to be polled along with authentication information and other settings in the form of a config file

## Getting Started
- Clone the repo and start the tool using go run .
- Pass the configuration settings using the flag -config \<relative path to json>
- Configuration json syntax
```
    "username": <username for the cluster> (default "Administrator")
    "password": <password for the cluster> (default "123456")
    "nodes": {
        <fts node ip address with port>:[ list of stats to poll from this node as string ]
    }
```
- Terminal User Interface commands
    - left arrow key and right arrow key to navigate between nodes list and stats list
    - down arrow key and up arrow key to navigate within the lists
    - a to select a stat for the left graph
    - d to select a stat for the right graph
    - s to select or unselect a node for the stat
    - q to quit the program