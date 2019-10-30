Cluster Healthz
===============

HTTP endpoint which executes standard checks on a cluster.

## Checks

* **NodeNotReady** - Check if all Nodes are `Ready`
* **NodeScalingIssue** - Check if all HorizontalPodAutoscaler objects have condition `ScalingActive=True`