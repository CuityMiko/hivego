//调度模块，负责从元数据库读取并解析调度信息。
//将需要执行的任务发送给执行模块，并读取返回信息。
//package schedule
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

//全局变量定义
var (
	gPort   string  // 监听端口号
	gDbConn *sql.DB //数据库链接

	gScds []Schedule //全局调度列表

	gchExScd chan ExecSchedule //执行的调度结构
)

//初始化工作
func init() {
	//从配置文件中获取数据库连接、服务端口号等信息
	gPort = ":8123"

}

//StartSchedule函数是调度模块的入口函数。程序初始化完成后，它负责连接元数据库，
//获取调度信息，在内存中构建Schedule结构。完成后，会调用Schedule的Timer方法。
//Timer方法会根据调度周期及启动时间，按时启动，随后会依据Schedule信息构建执行结构
//并送入chan中。
//模块的另一部分在不断的检测chan中的内容，将取到的执行结构体后创建新的goroutine
//执行。
func StartSchedule() error {
	// 连接数据库
	cnn, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/hive?charset=utf8")
	checkErr(err)
	gDbConn = cnn

	defer gDbConn.Close()

	//连接元数据库，初始化调度信息至内存
	tasks, err := getAllTasks() //获取Task列表
	checkErr(err)

	jobs, err := getAllJobs() //获取Job列表
	checkErr(err)

	schedules, err := getAllSchedules() //获取Schedule列表
	checkErr(err)

	//构建调度链信息
	for _, scd := range schedules {
		//设置调度中的job
		scd.job = jobs[scd.jobId]

		//设置job链
		for j := scd.job; j.nextJobId != 0; {
			j.nextJob = jobs[j.nextJobId]
			j.preJob = jobs[j.preJobId]
			j = j.nextJob

		}

	}

	fmt.Println(schedules, jobs, tasks)

	fmt.Println(schedules[0].job.nextJob.nextJob.nextJob.name)

	//当构建完成一个调度后，调用它的Timer方法。

	//从chan中得到需要执行的调度，启动一个线程执行

	return nil
}

func main() {
	StartSchedule()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}
