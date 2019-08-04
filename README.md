# Windivert2-Go
An partial go binding for Windivert2

## Code Example
[RabbitHole](https://github.com/RabbitYilia/RabbitHole)

## Something different

### DivertInit()
Load WinDivert.dll

### DivertPacket Struct
``` 
type DivertPacket struct {
	Data []byte //PacketData
	Addr WINDIVERTADDRESS //Only for network mode
}
type WINDIVERTADDRESS struct {
	Timestamp int64
	Flag      uint64
	IfIdx     uint64
}
```
