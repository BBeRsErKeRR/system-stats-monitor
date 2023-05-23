package stats

import (
	"fmt"
	"sort"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

func getAvgCPU(cpuTimesMetric []storage.Metric) (storage.CPUTimeStat, error) {
	var sumUser, sumSystem, sumIdle float64
	for _, metric := range cpuTimesMetric {
		stat := metric.StatInfo.(storage.CPUTimeStat)
		sumUser += stat.User
		sumSystem += stat.System
		sumIdle += stat.Idle
	}
	totalLen := float64(len(cpuTimesMetric))
	return storage.CPUTimeStat{
		User:   sumUser / totalLen,
		System: sumSystem / totalLen,
		Idle:   sumIdle / totalLen,
	}, nil
}

func getAvgDiskIo(statsItems []storage.Metric) ([]storage.DiskIoStatItem, error) {
	buff := make(map[string]*storage.DiskIoStatItem)
	buffLen := make(map[string]float64)
	for _, metric := range statsItems {
		stat := metric.StatInfo.(storage.DiskIoStatItem)
		val, ok := buff[stat.Device]
		if !ok {
			buff[stat.Device] = &stat
			buffLen[stat.Device] = 1
		} else {
			buffLen[stat.Device]++
			val.Tps += stat.Tps
			val.KbReadS += stat.KbReadS
			val.KbWriteS += stat.KbWriteS
		}
	}
	ioStats := make([]storage.DiskIoStatItem, 0, len(buff))
	for key, val := range buff {
		val.Tps /= buffLen[key]
		val.KbReadS /= buffLen[key]
		val.KbWriteS /= buffLen[key]
		ioStats = append(ioStats, *val)
	}
	return ioStats, nil
}

func getAvgLoad(lastLoadStats []storage.Metric) (storage.LoadStat, error) {
	var sumLoad1, sumLoad5, sumLoad15 float64
	for _, metric := range lastLoadStats {
		stat := metric.StatInfo.(storage.LoadStat)
		sumLoad1 += stat.Load1
		sumLoad5 += stat.Load5
		sumLoad15 += stat.Load15
	}
	totalLen := len(lastLoadStats)
	return storage.LoadStat{
		Load1:  sumLoad1 / float64(totalLen),
		Load5:  sumLoad5 / float64(totalLen),
		Load15: sumLoad15 / float64(totalLen),
	}, nil
}

func getUniqueDu(intSlice []storage.Metric) []storage.UsageStatItem {
	keys := make(map[string]bool)
	list := make([]storage.UsageStatItem, 0, len(intSlice))
	sort.Slice(intSlice, func(i, j int) bool {
		return intSlice[i].Date.Before(intSlice[j].Date)
	})

	for _, fact := range intSlice {
		stat := fact.StatInfo.(storage.UsageStatItem)
		entry := fmt.Sprintf("%v/%v", stat.Path, stat.Fstype)
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, stat)
		}
	}
	return list
}

func getAvgNetworkStates(nsStats []storage.Metric) (storage.NetworkStatesStat, error) {
	avgNs := make(map[string]int32)

	for _, fact := range nsStats {
		stat := fact.StatInfo.(storage.NetworkStatesStat)
		for name, counter := range stat.Counters {
			_, ok := avgNs[name]
			if ok {
				avgNs[name] += counter
			} else {
				avgNs[name] = counter
			}
		}
	}
	lengthStat := int32(len(nsStats))
	for name, counter := range avgNs {
		avgNs[name] = counter / lengthStat
	}
	return storage.NetworkStatesStat{
		Counters: avgNs,
	}, nil
}

func getUniqueNetworkStatistics(intSlice []storage.Metric) []storage.NetworkStatsItem {
	keys := make(map[string]bool)
	list := make([]storage.NetworkStatsItem, 0, len(intSlice))

	sort.Slice(intSlice, func(i, j int) bool {
		return intSlice[i].Date.Before(intSlice[j].Date)
	})

	for _, fact := range intSlice {
		stat := fact.StatInfo.(storage.NetworkStatsItem)
		entry := fmt.Sprintf("%v/%v/%v/%v", stat.Command, stat.Protocol, stat.PID, stat.Port)
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, stat)
		}
	}
	return list
}

func getUniqueBps(intSlice []storage.Metric, period int64) []storage.BpsItem {
	unique := make(map[string]storage.BpsItem)
	list := make([]storage.BpsItem, 0, len(intSlice))

	for _, fact := range intSlice {
		stat := fact.StatInfo.(storage.BpsItem)
		entry := fmt.Sprintf("%s->%s", stat.Source, stat.Destination)
		item, value := unique[entry]
		if !value {
			unique[entry] = stat
		} else {
			item.Numbers += stat.Numbers
			unique[entry] = item
		}
	}
	seconds := float64(period)
	for _, elem := range unique {
		elem.Bps = elem.Numbers / seconds
		list = append(list, elem)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Bps > list[j].Bps
	})

	return list
}

func getUniqueProtocolTalkers(intSlice []storage.Metric) []storage.ProtocolTalkerItem {
	var sumBytes float64
	unique := make(map[string]storage.ProtocolTalkerItem)
	list := make([]storage.ProtocolTalkerItem, 0, len(intSlice))

	for _, fact := range intSlice {
		stat := fact.StatInfo.(storage.ProtocolTalkerItem)
		entry := stat.Protocol
		item, value := unique[entry]
		if !value {
			unique[entry] = stat
		} else {
			item.SendBytes += stat.SendBytes
			unique[entry] = item
		}
		sumBytes += stat.SendBytes
	}

	for _, elem := range unique {
		elem.BytesPercentage = (elem.SendBytes / sumBytes) * 100.0
		list = append(list, elem)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Protocol < list[j].Protocol
	})

	return list
}
