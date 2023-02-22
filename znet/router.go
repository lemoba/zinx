package znet

import "github.com/lemoba/zinx/ziface"

// 实现router时，先嵌入BaseRouter基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct{}

/*
这里之所有BaseRouter都为空，是因为有的Router不希望有PreHandle、PostHandle这个两个业务
*/
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

func (br *BaseRouter) Handle(request ziface.IRequest) {}

func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
