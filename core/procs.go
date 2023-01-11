package core

import (
	"encoding/json"
	"processtool/log"
	"sync"

	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

type Process struct {
	Pid         int32                   `json:"pid,omitempty"`
	Name        string                  `json:"name,omitempty"`
	CmdLine     string                  `json:"cmd_line,omitempty"`
	CPUPercent  float64                 `json:"cpu_percent,omitempty"`
	MemPercent  float32                 `json:"mem_percent,omitempty"`
	MemoryInfo  *process.MemoryInfoStat `json:"memory_info,omitempty"`
	Connections []net.ConnectionStat    `json:"connections,omitempty"`
}

var (
	proc *Process
	once sync.Once
)

func Get() *Process {
	once.Do(func() { proc = &Process{} })
	return proc
}

func (p *Process) Processes() ([]Process, error) {
	var pg []Process
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}
	for _, v := range processes {
		name, err := v.Name()
		if err != nil {
			log.Debug(err)
			continue
		}
		cpuPercent, err := v.CPUPercent()
		if err != nil {
			log.Debug(err)
			continue
		}
		memPercent, err := v.MemoryPercent()
		if err != nil {
			log.Debug(err)
			continue
		}
		memoryInfo, err := v.MemoryInfo()
		if err != nil {
			log.Debug(err)
			continue
		}
		cmdLine, err := v.Cmdline()
		if err != nil {
			log.Debug(err)
			continue
		}
		connections, err := v.Connections()
		if err != nil {
			log.Debug(err)
			continue
		}
		pg = append(pg, Process{
			Pid:         v.Pid,
			Name:        name,
			CmdLine:     cmdLine,
			CPUPercent:  cpuPercent,
			MemPercent:  memPercent,
			MemoryInfo:  memoryInfo,
			Connections: connections,
		})
	}
	return pg, nil
}

func (p Process) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}
