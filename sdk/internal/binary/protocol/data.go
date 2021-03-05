package protocol

import "fmt"

type Data struct {
	name    string
	cmdSet  int
	cmdID   int
	retCode int
}

func NewData(name string, cmdSet, cmdID int) *Data {
	return &Data{
		name,
		cmdSet,
		cmdID,
		0,
	}
}

func (d *Data) CmdSet() int {
	return d.cmdSet
}

func (d *Data) CmdID() int {
	return d.cmdID
}

func (d *Data) CmdKey() int {
	return (d.cmdSet << 8) + d.cmdID
}

func (d *Data) PackReq() []byte {
	return nil
}

func (d *Data) UnpackResp(buf []byte, offset int) error {
	d.retCode = int(buf[offset])

	if d.retCode != 0 {
		return fmt.Errorf("response unpacking error")
	}

	return nil
}

func (d *Data) SetRetCode(retCode int) {
	d.retCode = retCode
}

func (d *Data) RetCode() int {
	return d.retCode
}

func (d *Data) String() string {
	return fmt.Sprintf("<%s cmset:%02x, cmdid:%02x>", d.name, d.cmdSet, d.cmdID)
}
