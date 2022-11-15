package znet

import "zinx/ziface"

// BaseRouter 实现route时， 先嵌入这个BaseRoute基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct{}

// 这里之所以BaseRouter的方法都为空，
// 是因为有的Router不希望有PreHandle或PostHandle
// 所以Router全部继承BaseRouter的好处是，不需要实现PreHandle和PostHandle也可以实例化

// PreHandle 在处理业务之前的钩子方法
func (b *BaseRouter) PreHandle(req ziface.IRequest) {}

// Handle 在处理业务的钩子方法
func (b *BaseRouter) Handle(req ziface.IRequest) {}

// PostHandle 在处理业务之后的钩子方法
func (b *BaseRouter) PostHandle(req ziface.IRequest) {}
