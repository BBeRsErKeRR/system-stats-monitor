package storage

type CPUTimeStat struct {
	User   float64 `json:"user"`
	System float64 `json:"system"`
	Idle   float64 `json:"idle"`
}

type LoadStat struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

type NetworkStatesStat struct {
	Counters map[string]int32 `json:"counters"`
}

type NetworkStatsItem struct {
	Command  string
	PID      int32
	User     int32
	Protocol string
	Port     int32
}

type UsageStatItem struct {
	Path                   string  `json:"path"`
	Fstype                 string  `json:"fstype"`
	Used                   int64   `json:"used"`
	AvailablePercent       float64 `json:"available_percent"`        //nolint:tagliatelle
	InodesUsed             int64   `json:"inodes_used"`              //nolint:tagliatelle
	InodesAvailablePercent float64 `json:"inodes_available_percent"` //nolint:tagliatelle
}

type DiskIoStatItem struct {
	Device   string  `json:"device"`
	Tps      float64 `json:"tps"`
	KbReadS  float64 `json:"kb_read_s"`  //nolint:tagliatelle
	KbWriteS float64 `json:"kb_write_s"` //nolint:tagliatelle
}

type ProtocolTalkerItem struct {
	Protocol        string  `json:"protocol"`
	SendBytes       float64 `json:"send_bytes"`       //nolint:tagliatelle
	BytesPercentage float64 `json:"bytes_percentage"` //nolint:tagliatelle
}

type BpsItem struct {
	Source      string  `json:"source"`
	Destination string  `json:"destination"`
	Protocol    string  `json:"protocol"`
	Bps         float64 `json:"bps"`
	Numbers     float64 `json:"numbers"`
}
