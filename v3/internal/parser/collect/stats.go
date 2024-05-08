package collect

import "time"

type Stats struct {
	NumPackages int
	NumServices int
	NumMethods  int
	NumEnums    int
	NumModels   int
	StartTime   time.Time
	EndTime     time.Time
}

func (stats *Stats) Start() {
	stats.StartTime = time.Now()
	stats.EndTime = stats.StartTime
}

func (stats *Stats) Stop() {
	stats.EndTime = time.Now()
}

func (stats *Stats) Elapsed() time.Duration {
	return stats.EndTime.Sub(stats.StartTime)
}

func (stats *Stats) Add(other *Stats) {
	stats.NumPackages += other.NumPackages
	stats.NumServices += other.NumServices
	stats.NumMethods += other.NumMethods
	stats.NumEnums += other.NumEnums
	stats.NumModels += other.NumModels
}
