package octolib

type Counter64Chan chan uint64
type Counter64 struct {
	CountCh Counter64Chan
	closech chan interface{}
}

func NewCounter64()(c *Counter64){
	c = &Counter64{
		CountCh: make(Counter64Chan, 5),
		closech: make(chan interface{}),
	}
	c.goCount()
	return c
}

func(c *Counter64)GetCountCh()(ch Counter64Chan){
	return c.CountCh
}

func(c *Counter64)Close(){
	close(c.closech)
}

//
// goCounter()
// Generates a UniqueID (int) and returns via supplied channel
// Runs forever
//
func (c *Counter64)goCount() {
	go func(c *Counter64) {
		defer func(){
			for{
				select{
				case	<-c.CountCh:
				default:
					close(c.CountCh)
					return
				}
			}
		}()

		var counter uint64 = 1
		for {
			select{
			case c.CountCh <- counter:
				counter += 1
			case <-c.closech:
				return
			}
		}
	}(c)
}
