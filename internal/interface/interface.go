package _interface

type Interface interface{ private() }

type Impl struct{}

func (Impl) private() {}
