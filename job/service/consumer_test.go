package service

import (
	"testing"
	"time"
)

func TestService_initConsumers(t *testing.T) {
	t.Logf("开始")
	svr.initConsumers()
	<-time.After(3 * time.Minute)
	//<-time.After(5 * time.Second)
	svr.ctxCancel()
}
