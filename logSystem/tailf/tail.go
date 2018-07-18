package tailf

import (
	"logSystem/conf"
	"github.com/hpcloud/tail"
	"time"
	"github.com/astaxie/beego/logs"
	"sync"
)

const (
	StatusNormal = 1
	StatusDelete = 2
)

// 存放单个tail对象
type TailObj struct {
	// tail对象
	tail *tail.Tail
	// 配置信息
	conf conf.CollectConf
	// 任务状态   配置是正常还是被删除了
	status int
	//
	exitChan chan int
}

//
type TextMsg struct {
	Msg   string // 读的这一行文本消息
	Topic string // 文本写到哪个topic里面
}

// 管理所有tail对象
type TailObjMgr struct {
	tailObjs []*TailObj
	msgChan  chan *TextMsg // 需要两个字段
	lock     sync.Mutex    // TailObjMgr.tailObjs 在多个goroutine间更新, 所以要加锁
}

var tailObjMgr *TailObjMgr

func GetSingleTail() (msg *TextMsg) {
	msg = <-tailObjMgr.msgChan
	return
}

/*
更新配置
*/
func UpdateConfig(confs []conf.CollectConf) (err error) {

	// 有可能有多个goroutine在更新同一份数据,需要加锁,更新完成后解锁
	tailObjMgr.lock.Lock()
	defer tailObjMgr.lock.Unlock()
	// 对比 "最新配置中" 与 "旧配置中" 对象是否一致, 对于已经运行的实例,他们的路径是不是一样的
	// 如果一样,说明已经在running
	// 否则, 找出新增的, 再添加一个任务
	for _, oneConf := range confs {
		var isRunning = false
		for _, obj := range tailObjMgr.tailObjs {
			if oneConf.LogPath == obj.conf.LogPath { // 说明日志已经运行了,没必要再跑了
				isRunning = true

				break
			}

		}
		if isRunning {
			continue
		}
		// 创建一个新的任务
		createNewTask(oneConf)
	}

	var tailObjs []*TailObj
	for _, obj := range tailObjMgr.tailObjs {
		obj.status = StatusDelete
		for _, oneConf := range confs {
			if oneConf.LogPath == obj.conf.LogPath {
				obj.status = StatusNormal
				break
			}
		}

		// 如果被删除了,需要停止任务, 通过channel停止
		if obj.status == StatusDelete {
			obj.exitChan <- 1
			continue
		}
		tailObjs = append(tailObjs, obj)
	}

	tailObjMgr.tailObjs = tailObjs

	return
}

func createNewTask(conf conf.CollectConf) {
	obj := &TailObj{
		conf:     conf,
		exitChan: make(chan int, 1),
	}

	tail, err := tail.TailFile(conf.LogPath, tail.Config{
		ReOpen:    true,
		Follow:    true,
		MustExist: false,
		Poll:      true,
	})

	if err != nil {
		logs.Error("collect filename[%s] failed, err:%v", conf.LogPath, err)
		return
	}

	obj.tail = tail

	tailObjMgr.tailObjs = append(tailObjMgr.tailObjs, obj)

	// 上面的代码是初始化配置

	// goroutine才是真正的去读取日志的内容,读取的日志内容放到channel里面
	// 然后主线程(main.go)去channel里面读消息,写到kafka里面
	go readFromTail(obj)

}

// 初始化tail包含多个日志
func InitTail(conf []conf.CollectConf, chanSize int) (err error) {

	tailObjMgr = &TailObjMgr{
		msgChan: make(chan *TextMsg, chanSize),
	}

	if len(conf) == 0 {
		logs.Error("invalid config fo log collect,conf:%v", conf)
		return
	}

	// 每个日志文件初始化一个tail对象去读取对应的业务日志
	for _, v := range conf {

		createNewTask(v)

		/*
		obj := &TailObj{
			conf: v,
		}

		tail, errtail := tail.TailFile(v.LogPath, tail.Config{
			ReOpen:    true,
			Follow:    true,
			MustExist: false,
			Poll:      true,
		})

		if errtail != nil {
			logs.Error("invaild config for log collect, cong:%v, err:%v", conf, errtail)
			continue
		}

		obj.tail = tail

		tailObjMgr.tailObjs = append(tailObjMgr.tailObjs, obj)

		// 上面的代码是初始化配置

		// goroutine才是真正的去读取日志的内容,读取的日志内容放到channel里面
		// 然后主线程(main.go)去channel里面读消息,写到kafka里面
		go readFromTail(obj)
		*/
	}

	return
}

func readFromTail(tailObj *TailObj) {
	for true {
		select {
		case line, ok := <-tailObj.tail.Lines:
			if !ok {
				logs.Warn("tail file close reopen, filename:%s\n", tailObj.tail.Filename)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			textMsg := &TextMsg{
				Msg:   line.Text,
				Topic: tailObj.conf.Topic,
			}
			tailObjMgr.msgChan <- textMsg

		case <-tailObj.exitChan:
			logs.Warn("tail obj will exited, conf:%v", tailObj.conf)
			return
		}
	}
}
