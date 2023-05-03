<div style="text-align: justify;">

# Frequently Asked Questions

## What to do if Chronos is not responsive on start and closes after some time?
Is the UI being displayed?
If it is not being displayed then it means that chronos is not able to contact the server for the list of search hostnames. This can happen due to the server in the middle of startup or if the node in the connection string is unresponsive.

If it is being displayed then the first thing to be done is to check the logs for additional error information. This can happen due to the following reasons:
Username or password for the cluster is wrong.
Connection string is incorrect.
Cluster has no active search nodes.
Search nodes send an invalid response to /api/chronosInit call.
One of the flags defined is wrong. Refer to the man page to identify the proper way to define flags for Chronos



## What to do if Chronos is abruptly closing without any warning?
If the program is shutting down in the middle of displaying the UI, then the first thing to do is to check the logs for the problem.
This can happen if the cluster rebalanced all the search nodes out or if one of the nodes sent an invalid response on /api/statsStream.


## What are threshold values?
Threshold values are the limits within which the stat value is supposed to be. If the stat crosses these thresholds, then it is considered erratic and an alert is generated for that stat. 
The thresholds defined are :
- _max_val to set the upper limit for the stat.
- _min_val to set the lower limit for the stat.
- _max_change to set the maximum percent change the stat can undergo.
- _max_change_val (along with _max_change) to set the amount of time for the percent change calculation.


## What flags does Chronos need to run?
Chronos requires 3 essential flags to run. They are 
- -username <Username for the cluster> (default – “Administrator”)
- -password <Password for the cluster> (default – “123456”)
- -connection_string <Connection string for the cluster> (default – “couchbase://127.0.0.1:12000”)

Chronos also takes a number of additional flags to enhance other functionalities such as alerts.
- -alert_TTL <TTL for an alert in seconds>
- -alert_data_padding <Amount of extra data for an alert in seconds>
- -\<stat Name> + <threshold> <Threshold value> (Examples: -pct_cpu_gc_max_val 0.05, -utilization:billableUnitsRate_min_val 2, -num_bytes_used_ram_max_change 0.5

Chronos also accepts stat thresholds for any index level stats for indexes already created before it is launched. Any indexes created after launch will have its stats displayed but will need a restart to accept threshold values for them.


## What to do if I have 3 nodes, but the graph is only displaying one line?
Chronos displays only one line when multiple lines overlap (have the same values at the time).


## How to display a graph for a stat?
You can select the stat using the mouse or the arrow keys or vim movement keys (hjkl) to. Once the required stat is highlighted, use “a” key to display a graph on the left side and “d” key to display the graph on the right side.


## How to hide or unhide a line on the graph?
You can toggle the lines using the center table. Select the node matching the color of the line and press enter to toggle it. To toggle a node for the graph that is not selected, you first need to go to its stat, select it again and then toggle the node.


## How to display a legend for a graph?
You can display the legend for the selected graph by pressing “p”.


## How to generate a report for an alert?
You can generate a report for an alert by selecting the alert and then pressing enter.
	
</div>