package Divert2

import (
	"log"
	"syscall"
	"unsafe"
)

var dll *syscall.LazyDLL
var DivertOpen *syscall.LazyProc
var DivertRecv *syscall.LazyProc
var DivertSend *syscall.LazyProc
var DivertShutdown *syscall.LazyProc
var DivertClose *syscall.LazyProc
var DivertCalcCheckSums *syscall.LazyProc
var RXChan chan *DivertPacket
var TXChan chan *DivertPacket
var EndFlag bool

type DivertPacket struct {
	Data []byte
	Addr WINDIVERTADDRESS
}

type WINDIVERTADDRESS struct {
	Timestamp int64
	Flag      uint64
	IfIdx     uint64
}

func SendOut(Handle uintptr, Data []byte) {
	thisPacket := DivertPacket{Data: Data, Addr: WINDIVERTADDRESS{Timestamp: 0, Flag: 131072}}
	WinDivertSend(Handle, &thisPacket)
}

func DivertInit() {
	dll = syscall.NewLazyDLL("WinDivert.dll")
	DivertOpen = dll.NewProc("WinDivertOpen")
	DivertRecv = dll.NewProc("WinDivertRecv")
	DivertSend = dll.NewProc("WinDivertSend")
	DivertShutdown = dll.NewProc("WinDivertShutdown")
	DivertClose = dll.NewProc("WinDivertClose")
	DivertCalcCheckSums = dll.NewProc("WinDivertHelperCalcChecksums")
	RXChan = make(chan *DivertPacket, 1000)
	TXChan = make(chan *DivertPacket, 1000)
	EndFlag = false
}

func WinDivertOpen(filter string, layer int, priority int16, flags uint64) (uintptr, error) {
	str := make([]byte, len(filter)+1)
	copy(str, filter)
	r, _, err := DivertOpen.Call(uintptr(unsafe.Pointer(&str[0])), uintptr(layer), uintptr(priority), uintptr(flags))
	if int(r) == -1 {
		return 0, err
	}
	return uintptr(r), nil
}

func WinDivertRecv(Handle uintptr) (*DivertPacket, error) {
	thisPacket := DivertPacket{Addr: WINDIVERTADDRESS{}}
	thisPacket.Data = make([]byte, 65535)
	packetLen := uint(0)
	r, _, err := DivertRecv.Call(Handle, uintptr(unsafe.Pointer(&thisPacket.Data[0])), uintptr(65535), uintptr(unsafe.Pointer(&packetLen)), uintptr(unsafe.Pointer(&thisPacket.Addr)))
	if int(r) == 0 {
		return nil, err
	}
	thisPacket.Data = thisPacket.Data[:packetLen]
	return &thisPacket, nil
}

func WinDivertSend(Handle uintptr, packet *DivertPacket) error {
	r, _, err := DivertSend.Call(Handle, uintptr(unsafe.Pointer(&packet.Data[0])), uintptr(len(packet.Data)), 0, uintptr(unsafe.Pointer(&packet.Addr)))
	if int(r) == 0 {
		log.Println(packet)
		return err
	}
	return nil
}

func WinDivertShutdown(Handle uintptr, How uint) error {
	r, _, err := DivertShutdown.Call(Handle, uintptr(How))
	if int(r) == 0 {
		return err
	}
	return nil
}

func WinDivertClose(Handle uintptr) error {
	r, _, err := DivertClose.Call(Handle)
	if int(r) == 0 {
		return err
	}
	return nil
}
