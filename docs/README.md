<div style="text-align: justify;">

# Chronos

## FTS monitoring/diagnostic tool

The chronos tool is designed to get live stats from a couchbase server and display them in a terminal user interface.

It takes the node ip information and the stats that needs to be polled along with authentication information and other settings as flags
It takes the node ip information and the stats that needs to be polled along with authentication information and other settings as flags

## Getting Started
- Clone the repo and start the tool using go run .
- The tool requires some information to be passed along as flags to connect to the server and configure alerts for the stats.
    - -username \<Username for the cluster> (default 'Administrator')
    - -password \<Password for the cluster> (default '123456')
    - -connection_string \<Connection string for the cluster> (default 'couchbases://127.0.0.1:12000')
    - -report \<Path to generate reports> (default './')
    - -alert_TTL \<Amount of time (in seconds) an alert should be visible in the UI> (default 120, max 600, min 1, type int)
    - -alert_data_padding \<Amount of data (in seconds) an alert should store before and after its triggered> (default 20, max 60, min 1, type int)
    - -\<stat name>_min_val \<Minimum threshold value for the stat. An alert will be generated if the stat falls below this limit> (type float)
    - -\<stat name>_max_val \<Maximum threshold value for the stat. An alert will be generated if the stat goes above this limit> (type float)
    - -\<stat name>_max_change \<Maximum percent change the stat can undergo in a certain duration of time> (type float)
    - -\<stat name>_max_change_time \<The time for which the max change for the stat is calculated> (default 1, type int)

## Terminal Commands
- Terminal User Interface commands
    - Left arrow key and right arrow key to navigate between the tables
    - Down arrow key and up arrow key to navigate within the table
    - 'a' key to select a stat for the left graph
    - 'd' key to select a stat for the right graph
    - 'Enter' to toggle selection of a node or to print a report
    - 'q' key to quit the program

## Log Information

The tool will store logs for any anamoly or event that happens while the program is running. This includes
- Couchbase go sdk unable to connect to the cluster.
- Unable to get a list of search nodes from the cluster.
- No search nodes on the cluster.
- Invalid or incorrect flags.
- Invalid or incorrect chronos init options from the server.
- Invalid server response or incorrect status codes from the server.
- Being unable to parse the server response body.
- Being unable to initialize the UI.
- Server closing connection unexpectedly.
- Alerts expiring.

## Stats Supported
- batch_bytes_added
- batch_bytes_removed
- curr_batches_blocked_by_herder
- num_batches_introduced
- num_bytes_used_ram
- num_gocbcore_dcp_agents
- num_gocbcore_stats_agents
- pct_cpu_gc
- tot_batches_merged
- tot_batches_new
- tot_bleve_dest_closed
- tot_bleve_dest_opened
- tot_queryreject_on_memquota
- tot_rollback_full
- tot_rollback_partial
- total_gc
- total_queries_rejected_by_herder
- utilization:billableUnitsRate
- utilization:cpuPercent
- utilization:diskBytes
- utilization:memoryBytes

## Example Start Commands

- go run . -username Administrator -password asdasd -connection_string couchbase://127.0.0.1:12000
- go run . -username Test -password 123456 -connection_string couchbase://192.183.42.7:12000 -report ~/Desktop/ -alert_TTL 30 -alert_data_padding 10
- go run . -username Administrator -password 123456 -connection_string couchbase://192.173.39.128:12000 -report ~/Documents/Reports/ -alert_TTL 150 -alert_data_padding 50 -tot_query_reject_on_memquota_max_val 100 -pct_cpu_gc_max_change 0.5 -total_gc_max_change 0.5 -total_gc_max_change_time 2 -num_bytes_used_ram_min_val 50000

</div>