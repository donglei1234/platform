package utils

import "time"

// TimestampDurationDayWithAnchor 计算从start时间开始到end时间结束经过的天数
// start: 开始时间戳
// end: 结束时间戳
// anchor (int): 锚点小时，格式为 5(凌晨5点为一天) 或 6(凌晨6点为一天)
// return: 经过的天数
func TimestampDurationDayWithAnchor(start, end int64, anchorHour int32) int32 {
	if start >= end {
		return 0
	}

	startTime := time.Unix(start, 0)
	endTime := time.Unix(end, 0)

	// 计算锚点时间
	anchorTime := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), int(anchorHour), 0, 0, 0, startTime.Location())

	// 计算天数
	days := int32(endTime.Sub(anchorTime).Hours() / 24)

	// 计算锚点时间是否超过结束时间
	if anchorTime.AddDate(0, 0, int(days)).After(endTime) {
		days--
	}

	return days
}
