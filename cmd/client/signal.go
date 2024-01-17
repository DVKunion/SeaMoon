package client

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/DVKunion/SeaMoon/pkg/transfer"
)

type SigGroup struct {
	wg           *sync.WaitGroup
	WatchChannel chan os.Signal
	StartChannel chan transfer.Type
	StopChannel  chan transfer.Type
}

func NewSigGroup() *SigGroup {
	sg := &SigGroup{
		new(sync.WaitGroup),
		make(chan os.Signal, 1),
		make(chan transfer.Type, 1),
		make(chan transfer.Type, 1),
	}
	signal.Notify(sg.WatchChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	return sg
}

func (sg *SigGroup) StartHttpProxy() {
	if Config().Http.Status == "active" {
		return
	}
	sg.StartChannel <- transfer.HTTP
	Config().Http.Status = "active"
}

func (sg *SigGroup) StopHttpProxy() {
	if Config().Http.Status == "inactive" {
		return
	}
	sg.StopChannel <- transfer.HTTP
	Config().Http.Status = "inactive"
}

func (sg *SigGroup) StartSocksProxy() {
	if Config().Socks5.Status == "active" {
		return
	}
	sg.StartChannel <- transfer.SOCKS5
	Config().Socks5.Status = "active"
}

func (sg *SigGroup) StopSocksProxy() {
	if Config().Socks5.Status == "inactive" {
		return
	}
	sg.StopChannel <- transfer.SOCKS5
	Config().Socks5.Status = "inactive"
}

func (sg *SigGroup) StopProxy() {
	sg.StopHttpProxy()
	sg.StopSocksProxy()
	Config().Save()
}

//func (sg *SigGroup) Stop() {
//	close(sg.SocksStopChannel)
//	close(sg.HttpStopChannel)
//	close(sg.SocksStartChannel)
//	close(sg.HttpStartChannel)
//	close(sg.WatchChannel)
//}

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
