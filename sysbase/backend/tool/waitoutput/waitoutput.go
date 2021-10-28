package waitoutput

import "sync"

type message struct {
	cmd  string
	data chan string
}

type WaitOutput struct {
	sync.RWMutex
	Num     int
	Message []*message
	Cmds    []string
}

func (wo *WaitOutput) SetDataChan(cmd string) {
	wo.Lock()
	defer wo.Unlock()

	has := false
	for _, v := range wo.Message {
		if v.cmd == cmd {
			has = true
		}
	}
	if !has {
		m := &message{
			cmd:  cmd,
			data: make(chan string, 1),
		}
		wo.Num = wo.Num + 1
		wo.Message = append(wo.Message, m)
		wo.Cmds = append(wo.Cmds, cmd)
	}
}

func (wo *WaitOutput) GetDataChan(cmd string) chan string {
	wo.RLock()
	defer wo.RUnlock()

	for _, v := range wo.Message {
		if v.cmd == cmd {
			return v.data
		}
	}
	return nil
}

func (wo *WaitOutput) IsRunning(cmd string) bool {
	wo.RLock()
	defer wo.RUnlock()

	for _, v := range wo.Message {
		if v.cmd == cmd {
			return true
		}
	}
	return false
}

func (wo *WaitOutput) DeleteByCmd(cmd string) {
	wo.Lock()
	defer wo.Unlock()

	var ms []*message
	has := false
	for _, v := range wo.Message {
		if v.cmd == cmd {
			has = true
			if !isStringChanClosed(v.data) {
				close(v.data)
			}
		} else {
			ms = append(ms, v)
		}
	}
	if has && wo.Num > 0 {
		wo.Message = ms
		wo.Cmds = append(wo.Cmds[:wo.Num-1], wo.Cmds[wo.Num:]...)
		wo.Num = wo.Num - 1
	}
}

func isStringChanClosed(ch <-chan string) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}
