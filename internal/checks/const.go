package checks

const (
	// LevelCritical is used to determine which components are critical to a clusters operations.
	LevelCritical = "critical"

	// ReasonCrashLoopBackOff is used for determining if a Pod has been restarted too many times.
	ReasonCrashLoopBackOff = "CrashLoopBackOff"

	// LabelClusterHealthz is used when querying for resources which need to be checked.
	LabelClusterHealthz = "skpr.io/cluster-healthz"
)
