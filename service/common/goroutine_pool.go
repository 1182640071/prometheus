package common

import "sync"

type GPool struct {
	chanCount chan int
	wg *sync.WaitGroup
}

func NewPool(size int) *GPool{

	return &GPool{
		chanCount:make(chan int,size),
		wg:&sync.WaitGroup{},
	}
}

func (p *GPool)Add(){

	p.chanCount <- 1
	p.wg.Add(1)
}

func (p *GPool) JobDone(){
	<- p.chanCount
	p.wg.Done()
}