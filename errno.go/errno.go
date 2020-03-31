package errno

import "fmt"

//CopeError 框架自定义错误,可以被映射
type CopeError struct {
	Code    int
	Message string
}

func (c *CopeError) Error() string {
	return fmt.Sprintf("cope error, code:%d message:%v", c.Code, c.Message)
}

var (
	NotHaveInstance = &CopeError{
		Code:    1,
		Message: "not have instance",
	}
	ConnFailed = &CopeError{
		Code:    2,
		Message: "connect failed",
	}
	InvalidNode = &CopeError{
		Code:    3,
		Message: "invalid node",
	}
	AllNodeFailed = &CopeError{
		Code:    4,
		Message: "all node failed",
	}
)

//IsConnectError 连接是否超时
func IsConnectError(err error) bool {
	res, ok := err.(*CopeError)
	if !ok {
		return false
	}
	if res.Error() == AllNodeFailed.Error() {
		return true
	}

	return false
}
