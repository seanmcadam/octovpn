package layers

import (
	"testing"

	"github.com/seanmcadam/bufferpool"
)

func TestCompile(t *testing.T){}

func TestCreateLayerPair(t *testing.T){
 _, _ = CreateLayerPair()
}

func TestLayerPair(t *testing.T){
	pool := bufferpool.New()

a, b := CreateLayerPair()

abuf := pool.Get().Append([]byte("A"))
bbuf := pool.Get().Append([]byte("B"))


a.Send(abuf)
b.Send(bbuf)

<-a.RecvCh()
<-b.RecvCh()

}