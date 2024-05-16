// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package gops

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ergoapi/util/exstr"
	"github.com/google/gops/goprocess"
	"github.com/shirou/gopsutil/v3/process"
)

var develRe = regexp.MustCompile(`devel\s+\+\w+`)

func Processes() {
	ps := goprocess.FindAll()

	var maxPID, maxPPID, maxExec, maxVersion int
	for i, p := range ps {
		ps[i].BuildVersion = shortenVersion(p.BuildVersion)
		maxPID = max(maxPID, len(strconv.Itoa(p.PID)))
		maxPPID = max(maxPPID, len(strconv.Itoa(p.PPID)))
		maxExec = max(maxExec, len(p.Exec))
		maxVersion = max(maxVersion, len(ps[i].BuildVersion))
	}

	for _, p := range ps {
		buf := bytes.NewBuffer(nil)
		pid := strconv.Itoa(p.PID)
		fmt.Fprint(buf, pad(pid, maxPID))
		fmt.Fprint(buf, " ")
		ppid := strconv.Itoa(p.PPID)
		fmt.Fprint(buf, pad(ppid, maxPPID))
		fmt.Fprint(buf, " ")
		fmt.Fprint(buf, pad(p.Exec, maxExec))
		if p.Agent {
			fmt.Fprint(buf, "*")
		} else {
			fmt.Fprint(buf, " ")
		}
		fmt.Fprint(buf, " ")
		fmt.Fprint(buf, pad(p.BuildVersion, maxVersion))
		fmt.Fprint(buf, " ")
		fmt.Fprint(buf, p.Path)
		fmt.Fprintln(buf)
		buf.WriteTo(os.Stdout)
	}
}

func shortenVersion(v string) string {
	if !strings.HasPrefix(v, "devel") {
		return v
	}
	results := develRe.FindAllString(v, 1)
	if len(results) == 0 {
		return v
	}
	return results[0]
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}

func pad(s string, total int) string {
	if len(s) >= total {
		return s
	}
	return s + strings.Repeat(" ", total-len(s))
}

// ProcessInfo takes arguments starting with pid|:addr and grabs all kinds of
// useful Go process information.
func ProcessInfo(args []string) {
	pid := exstr.Str2Int32(args[0])
	var period time.Duration
	var err error
	if len(args) >= 2 {
		period, err = time.ParseDuration(args[1])
		if err != nil {
			secs, _ := strconv.Atoi(args[1])
			period = time.Duration(secs) * time.Second
		}
	}
	processInfo(pid, period)
}

func processInfo(pid int32, period time.Duration) {
	if period < 0 {
		log.Fatalf("Cannot determine CPU usage for negative duration %v", period)
	}
	p, err := process.NewProcess(pid)
	if err != nil {
		log.Fatalf("Cannot read process info: %v", err)
	}
	if v, err := p.Parent(); err == nil {
		fmt.Printf("parent PID:\t%v\n", v.Pid)
	}
	if v, err := p.NumThreads(); err == nil {
		fmt.Printf("threads:\t%v\n", v)
	}
	if v, err := p.MemoryPercent(); err == nil {
		fmt.Printf("memory usage:\t%.3f%%\n", v)
	}
	if v, err := p.CPUPercent(); err == nil {
		fmt.Printf("cpu usage:\t%.3f%%\n", v)
	}
	if period > 0 {
		if v, err := cpuPercentWithinTime(p, period); err == nil {
			fmt.Printf("cpu usage (%v):\t%.3f%%\n", period, v)
		}
	}
	if v, err := p.Username(); err == nil {
		fmt.Printf("username:\t%v\n", v)
	}
	if v, err := p.Cmdline(); err == nil {
		fmt.Printf("cmd+args:\t%v\n", v)
	}
	if v, err := elapsedTime(p); err == nil {
		fmt.Printf("elapsed time:\t%v\n", v)
	}
	if v, err := p.Connections(); err == nil {
		if len(v) > 0 {
			for _, conn := range v {
				fmt.Printf("local/remote:\t%v:%v <-> %v:%v (%v)\n",
					conn.Laddr.IP, conn.Laddr.Port, conn.Raddr.IP, conn.Raddr.Port, conn.Status)
			}
		}
	}
}

// cpuPercentWithinTime return how many percent of the CPU time this process uses within given time duration
func cpuPercentWithinTime(p *process.Process, t time.Duration) (float64, error) {
	cput, err := p.Times()
	if err != nil {
		return 0, err
	}
	time.Sleep(t)
	cput2, err := p.Times()
	if err != nil {
		return 0, err
	}
	// nolint:staticcheck
	return 100 * (cput2.Total() - cput.Total()) / t.Seconds(), nil
}

// elapsedTime shows the elapsed time of the process indicating how long the
// process has been running for.
func elapsedTime(p *process.Process) (string, error) {
	crtTime, err := p.CreateTime()
	if err != nil {
		return "", err
	}
	etime := time.Since(time.Unix(crtTime/1000, 0))
	return fmtEtimeDuration(etime), nil
}

// fmtEtimeDuration formats etime's duration based on ps' format:
// [[DD-]hh:]mm:ss
// format specification: http://linuxcommand.org/lc3_man_pages/ps1.html
func fmtEtimeDuration(d time.Duration) string {
	days := d / (24 * time.Hour)
	hours := d % (24 * time.Hour)
	minutes := hours % time.Hour
	seconds := math.Mod(minutes.Seconds(), 60)
	var b strings.Builder
	if days > 0 {
		fmt.Fprintf(&b, "%02d-", days)
	}
	if days > 0 || hours/time.Hour > 0 {
		fmt.Fprintf(&b, "%02d:", hours/time.Hour)
	}
	fmt.Fprintf(&b, "%02d:", minutes/time.Minute)
	fmt.Fprintf(&b, "%02.0f", seconds)
	return b.String()
}
