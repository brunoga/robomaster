package controller

type enableFunction int

const (
	blooded enableFunction = iota + 2
	movementControl
	gunControl
	offControl
)

type functionEnableInfo struct {
	Id     enableFunction `json:"id"`
	Enable bool           `json:"enable"`
}

type functionEnableParameter struct {
	List []*functionEnableInfo `json:"list"`
}

func (f *functionEnableParameter) setFunctionEnable(id enableFunction,
	enable bool) {
	for _, info := range f.List {
		if info.Id == id {
			info.Enable = enable
			return
		}
	}

	f.List = append(f.List, &functionEnableInfo{
		enableFunction(id),
		enable,
	})
}
