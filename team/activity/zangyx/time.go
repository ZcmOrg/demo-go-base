package main

import (
	"fmt"
	"time"
)

func main() {
	// 返回时间点t对应的年份。
	fmt.Println(time.Now().Year())

	// Month代表一年的某个月。
	fmt.Println(time.Now().Month())

	//Weekday代表一周的某一天。
	fmt.Println(time.Now().Weekday())

	// 返回时间点t对应那一月的第几日。
	fmt.Println(time.Now().Day())

	// 返回t对应的那一天的第几小时，范围[0, 23]。
	fmt.Println(time.Now().Hour())

	// 返回t对应的那一小时的第几分种，范围[0, 59]。
	fmt.Println(time.Now().Minute())

	// 返回t对应的那一分钟的第几秒，范围[0, 59]。
	fmt.Println(time.Now().Second())

	// 返回t对应的那一秒内的纳秒偏移量，范围[0, 999999999]。
	fmt.Println(time.Now().Nanosecond())

	// Now返回当前本地时间。
	fmt.Println(time.Now())

}
