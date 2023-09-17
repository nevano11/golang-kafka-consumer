package filter

import "fmt"

const (
	ParamLike     = 1
	ParamEq       = 2
	ParamPageNum  = 3
	ParamPageSize = 4
)

type Option struct {
	ParamType int
	Name      string
	Val       string
}

func NewOption(name, val string, dtype int) Option {
	return Option{
		Name:      name,
		ParamType: dtype,
		Val:       val,
	}
}

func (o *Option) String() string {
	return fmt.Sprintf("{%s, %d, %s}", o.Name, o.ParamType, o.Val)
}
