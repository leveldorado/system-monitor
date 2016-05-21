package systemMonitor

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
)

var numberRegexp = regexp.MustCompile(`[0-9]+`)

const (
	procMemInfo = "/proc/meminfo"
)

func loadProcMemInfo() ([]byte, error) {
	return ioutil.ReadFile(procMemInfo)
}

type ProcMemInfo struct {
	MemTotal     uint
	MemFree      uint
	MemAvailable uint
}

func GetProcMemInfo() (*ProcMemInfo, error) {
	data, err := loadProcMemInfo()
	if err != nil {
		return nil, err
	}
	procMemInfo := &ProcMemInfo{}
	err = procMemInfo.FromText(data)
	return procMemInfo, err
}

func (p ProcMemInfo) InMb() *ProcMemInfo {
	mbInfo := &ProcMemInfo{}
	mbInfo.MemAvailable = p.MemAvailable / 1024
	mbInfo.MemFree = p.MemFree / 1024
	mbInfo.MemTotal = p.MemTotal / 1024
	return mbInfo
}

func (p *ProcMemInfo) FromText(data []byte) error {
	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		lineData := bytes.Split(line, []byte(`:`))
		if len(lineData) != 2 {
			fmt.Println(string(line))
			continue
		}
		number, err := strconv.ParseUint(string(numberRegexp.Find(lineData[1])), 10, 64)
		if err != nil {
			return err
		}
		switch string(lineData[0]) {
		case "MemTotal":
			p.MemTotal = uint(number)
		case "MemFree":
			p.MemFree = uint(number)
		case "MemAvailable":
			p.MemAvailable = uint(number)
		}
	}
	return nil
}
