package main

import (
    "github.com/robfig/cron"
    _ "github.com/ying32/govcl/pkgs/winappres"
    "github.com/ying32/govcl/vcl"
    "log"
)

var (
    bc = cron.New()
    failSlice []*User
    logger *log.Logger
)


const (
    //cron表达式
    //数据读取后创建定时器的时间
    morningEndSepc = "0 30 08 * * ?"
    noonEndSpec = "0 20 13 * * ?"
    eveningEndSpec = "0 30 20 * * ?"

    //处理失败队列的时间
    morningProcessFailSepc = "0 00 8 * * ?"
    noonProcessFailSepc = "0 00 13 * * ?"
    eveningProcessFailSepc = "0 30 19 * * ?"

    //新的处理
    morningSpec2 = "0 10 07 * * ?"
    noonSpec2 = "0 55 11 * * ?"
    eveningSpec2 = "0 35 18 * * ?"
)

func main() {

    go createLog()



    go bc.Start()
    //go WebStart()
    vcl.RunApp(&Form1)
}

func GetCronTime(hour,min string) string {
    return "0 "+min+" "+hour+" * * ?"
}