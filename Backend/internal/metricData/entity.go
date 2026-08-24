package metricdata

import "time"

type server_metric struct {
	Time              time.Time `json:"time"` 
	CPUUtilization    float64   `json:"cpu_util"`
	MemoryUtilization float64   `json:"mem_util"`
	DiskReadBytes     int64     `json:"disk_r"`
	DiskWriteBytes    int64     `json:"disk_w"`
}
