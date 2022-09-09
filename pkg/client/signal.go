package client

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type ControlSignal int8

const (
	HttpProxyStartSignal ControlSignal = iota + 1
	HttpProxyStopSignal
	SocksProxyStartSignal
	SocksProxyStopSignal
)

type SigGroup struct {
	wg                *sync.WaitGroup
	WatchChannel      chan os.Signal
	HttpStartChannel  chan ControlSignal
	HttpStopChannel   chan ControlSignal
	SocksStartChannel chan ControlSignal
	SocksStopChannel  chan ControlSignal
}

func NewSigGroup() *SigGroup {
	sg := &SigGroup{
		new(sync.WaitGroup),
		make(chan os.Signal, 1),
		make(chan ControlSignal, 1),
		make(chan ControlSignal, 1),
		make(chan ControlSignal, 1),
		make(chan ControlSignal, 1),
	}
	signal.Notify(sg.WatchChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	return sg
}

func (sg *SigGroup) StartHttpProxy() {
	if Config().Http.Status == "active" {
		return
	}
	sg.HttpStartChannel <- HttpProxyStartSignal
	Config().Http.Status = "active"
}

func (sg *SigGroup) StopHttpProxy() {
	if Config().Http.Status == "inactive" {
		return
	}
	sg.HttpStopChannel <- HttpProxyStopSignal
	Config().Http.Status = "inactive"
}

func (sg *SigGroup) StartSocksProxy() {
	if Config().Socks5.Status == "active" {
		return
	}
	sg.SocksStartChannel <- SocksProxyStartSignal
	Config().Socks5.Status = "active"
}

func (sg *SigGroup) StopSocksProxy() {
	if Config().Socks5.Status == "inactive" {
		return
	}
	sg.SocksStopChannel <- SocksProxyStopSignal
	Config().Socks5.Status = "inactive"
}

func (sg *SigGroup) StopProxy() {
	sg.StopHttpProxy()
	sg.StopSocksProxy()
	Config().Save()
}

func (sg *SigGroup) Stop() {
	close(sg.SocksStopChannel)
	close(sg.HttpStopChannel)
	close(sg.SocksStartChannel)
	close(sg.HttpStartChannel)
	close(sg.WatchChannel)
}

func (sg *SigGroup) Detection() {
	if Config().Http.Enabled {
		sg.StartHttpProxy()
	} else {
		sg.StopHttpProxy()
	}
	if Config().Socks5.Enabled {
		sg.StartSocksProxy()
	} else {
		sg.StopSocksProxy()
	}
}
