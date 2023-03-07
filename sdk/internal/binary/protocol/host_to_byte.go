package protocol

func HostToByte(host, index byte) byte {
	return index*32 + host
}
