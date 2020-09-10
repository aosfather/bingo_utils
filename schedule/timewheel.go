package schedule

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

// time wheel struct
type TimeWheel struct {
	interval       time.Duration
	ticker         *time.Ticker
	slots          []*list.List
	currentPos     int
	slotNum        int
	addTaskChannel chan *task
	stopChannel    chan bool
	taskRecord     *sync.Map
}

type Jobhandle interface {
	TaskName() string
	Execute()
}

type SimpleJobHandle struct {
	Name    string
	Data    TaskData
	JobFunc Job
}

func (this *SimpleJobHandle) TaskName() string {
	return this.Name
}

func (this *SimpleJobHandle) Execute() {
	if this.JobFunc != nil {
		this.JobFunc(this.Data)
	}
}

// Job callback function
type Job func(TaskData)

// TaskData callback params
type TaskData map[interface{}]interface{}

// task struct
type task struct {
	interval    time.Duration
	times       int //-1:no limit >=1:run times
	circle      int
	key         interface{}
	job         Job
	taskData    TaskData
	schedule    Schedule
	lasttrigger *time.Time
}

func NewSecondWheel(slotNum int) *TimeWheel {
	return New(time.Second, slotNum)
}

func NewMinuteWheel(slotNum int) *TimeWheel {
	return New(time.Minute, slotNum)
}

// New create a empty time wheel
func New(interval time.Duration, slotNum int) *TimeWheel {
	if interval <= 0 || slotNum <= 0 {
		return nil
	}
	tw := &TimeWheel{
		interval:       interval,
		slots:          make([]*list.List, slotNum),
		currentPos:     0,
		slotNum:        slotNum,
		addTaskChannel: make(chan *task),
		stopChannel:    make(chan bool),
		taskRecord:     &sync.Map{},
	}

	tw.init()

	return tw
}

// Start start the time wheel
func (tw *TimeWheel) Start() {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.start()
}

// Stop stop the time wheel
func (tw *TimeWheel) Stop() {
	tw.stopChannel <- true
}

func (tw *TimeWheel) start() {
	for {
		select {
		case <-tw.ticker.C:
			tw.tickHandler()
		case task := <-tw.addTaskChannel:
			tw.addTask(task)
		case <-tw.stopChannel:
			tw.ticker.Stop()
			return
		}
	}
}

func (tw *TimeWheel) AddJobByCron(cron string, job Jobhandle) error {
	if cron == "" || job == nil {
		return errors.New("illegal task params")
	}
	return tw.AddTaskByCron(cron, job.TaskName(), nil, func(t TaskData) { job.Execute() })
}

func (tw *TimeWheel) AddTaskByCron(cron string, key interface{}, data TaskData, job Job) error {
	if cron == "" || key == nil || job == nil {
		return errors.New("illegal task params")
	}
	_, ok := tw.taskRecord.Load(key)
	if ok {
		return errors.New("duplicate task key")
	}

	p := NewParser(Second | Minute | Hour | Dom | Month | Dow)
	s, err := p.Parse(cron)
	if err != nil {
		return err
	}
	now := time.Now()
	happertime := s.Next(now)
	interval := happertime.Sub(now)
	tw.addTaskChannel <- &task{interval: interval, times: 1, key: key, taskData: data, job: job, schedule: s, lasttrigger: &happertime}

	return nil
}

func (tw *TimeWheel) AddForeverTaskByHours(t int64, key interface{}, data TaskData, job Job) error {
	return tw.AddTask(time.Duration(t)*time.Hour, -1, key, data, job)
}

func (tw *TimeWheel) AddForeverTaskByMinutes(t int64, key interface{}, data TaskData, job Job) error {
	return tw.AddTask(time.Duration(t)*time.Minute, -1, key, data, job)
}

func (tw *TimeWheel) AddForeverTaskBySeconds(t int64, key interface{}, data TaskData, job Job) error {
	return tw.AddTask(time.Duration(t)*time.Second, -1, key, data, job)
}

//
func (tw *TimeWheel) AddOneTimeSimpleTask(t time.Time, job Job) error {
	return tw.AddOneTimeTask(t, nil, job)
}

func (tw *TimeWheel) AddOneTimeTask(t time.Time, data TaskData, job Job) error {
	now := time.Now()
	interval := t.Sub(now)
	return tw.AddTask(interval, 1, &now, data, job)
}

func (tw *TimeWheel) AddCycleSimpleTask(interval time.Duration, times int, key interface{}, job Job) error {
	return tw.AddTask(interval, times, key, nil, job)
}

// AddTask add new task to the time wheel
func (tw *TimeWheel) AddTask(interval time.Duration, times int, key interface{}, data TaskData, job Job) error {
	if interval <= 0 || key == nil || job == nil || times < -1 || times == 0 {
		return errors.New("illegal task params")
	}

	_, ok := tw.taskRecord.Load(key)
	if ok {
		return errors.New("duplicate task key")
	}

	tw.addTaskChannel <- &task{interval: interval, times: times, key: key, taskData: data, job: job}
	return nil
}

// RemoveTask remove the task from time wheel
func (tw *TimeWheel) RemoveTask(key interface{}) error {
	if key == nil {
		return nil
	}

	value, ok := tw.taskRecord.Load(key)

	if !ok {
		return errors.New("task not exists, please check you task key")
	} else {
		// lazy remove task
		task := value.(*task)
		task.times = 0
		tw.taskRecord.Delete(task.key)
	}
	return nil
}

// UpdateTask update task times and data
func (tw *TimeWheel) UpdateTask(key interface{}, interval time.Duration, taskData TaskData) error {
	if key == nil {
		return errors.New("illegal key, please try again")
	}

	value, ok := tw.taskRecord.Load(key)

	if !ok {
		return errors.New("task not exists, please check you task key")
	}
	task := value.(*task)
	task.taskData = taskData
	task.interval = interval
	return nil
}

// time wheel initialize
func (tw *TimeWheel) init() {
	for i := 0; i < tw.slotNum; i++ {
		tw.slots[i] = list.New()
	}
}

//
func (tw *TimeWheel) tickHandler() {
	l := tw.slots[tw.currentPos]
	tw.scanAddRunTask(l)
	if tw.currentPos == tw.slotNum-1 {
		tw.currentPos = 0
	} else {
		tw.currentPos++
	}
}

// add task
func (tw *TimeWheel) addTask(task *task) {
	if task.times == 0 {
		return
	}

	pos, circle := tw.getPositionAndCircle(task.interval)
	task.circle = circle

	tw.slots[pos].PushBack(task)

	//record the task
	tw.taskRecord.Store(task.key, task)
}

// scan task list and run the task
func (tw *TimeWheel) scanAddRunTask(l *list.List) {

	if l == nil {
		return
	}

	for item := l.Front(); item != nil; {
		task := item.Value.(*task)

		if task.times == 0 {
			next := item.Next()
			l.Remove(item)
			tw.taskRecord.Delete(task.key)
			item = next
			continue
		}

		if task.circle > 0 {
			task.circle--
			item = item.Next()
			continue
		}

		go task.job(task.taskData)
		next := item.Next()
		l.Remove(item)
		item = next

		if task.times == 1 {
			if task.schedule != nil {
				happertime := task.schedule.Next(*task.lasttrigger)
				task.interval = happertime.Sub(*task.lasttrigger)
				task.lasttrigger = &happertime
				if task.interval > 0 {
					tw.addTask(task)
					continue
				}

			}

			task.times = 0
			tw.taskRecord.Delete(task.key)

		} else {
			if task.times > 0 {
				task.times--
			}
			tw.addTask(task)
		}
	}
}

// get the task position
func (tw *TimeWheel) getPositionAndCircle(d time.Duration) (pos int, circle int) {
	delaySeconds := int(d.Seconds())
	intervalSeconds := int(tw.interval.Seconds())
	circle = int(delaySeconds / intervalSeconds / tw.slotNum)
	pos = int(tw.currentPos+delaySeconds/intervalSeconds) % tw.slotNum
	return
}
