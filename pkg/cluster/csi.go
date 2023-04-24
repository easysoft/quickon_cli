package cluster

type CSI struct {
}

func NewCSI() *CSI {
	return &CSI{}
}

func (c *CSI) Check() bool {
	return false
}

func (c *CSI) Upgrade() error {
	return nil
}
