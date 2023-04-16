package main

import (
	"log"
	"sync"
	"time"
)

type TaskState uint8

const (
	TaskEnd TaskState = iota
	TaskBegin
	TaskIng
)

type Task struct {
	ID uint32
	RequestChan chan *Pack
	ResponseChan chan *Pack
	State TaskState
	_packs Packs
	_updateTime int64
}

func NewTask(pack *Pack,responseChan chan *Pack) *Task {
	task := &Task{
		ID:         pack.ID,
		_packs:      Packs{pack},
		_updateTime: time.Now().UnixMilli(),
		State: TaskBegin,
		RequestChan: make(chan *Pack),
		ResponseChan: responseChan,
	}
	return task
}

func (t *Task) run()  {
	defer close(t.RequestChan)
	var ti = time.NewTicker(time.Millisecond*300)
	for  {
		select {
		case pack := <- t.RequestChan:
			d := time.Now().UnixMilli()-t._updateTime
			if d <= 30 { d = 100 }
			t._updateTime = d

			t._packs.Append(pack)
			if !t._packs.CheckIntegrality() {
				ti.Reset(time.Millisecond * time.Duration(d))
				continue
			}
			t.State = TaskEnd
			ti.Stop()
			t.ResponseChan <- t._packs.Marge()
		case <-ti.C:
			log.Println("超时任务：",t.ID)
		}
	}
}

func  LoadTaskSyncMayCache(m *sync.Map,id uint32) (_task *Task,exits bool) {
	v,ok := m.Load(id)
	if ok {
		_task,ok = v.(*Task)
		if ok {
			return _task,true
		}
	}
	return nil,false
}

func  StoreTaskSyncMayCache(m *sync.Map,_task *Task)  {
	m.Store(_task.ID,_task)
}

