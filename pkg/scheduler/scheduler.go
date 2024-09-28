package scheduler

// So much to improve
// make it sateless multi instance
// definintions from db or config?
// for now it gets the job done for only the monolithic apps

import (
	"fmt"
	"time"
)

type Task struct {
	hour, min, sec int
	taskFunc       func()
	nextRun        time.Time
	name           string
}

type Scheduler struct {
	tasks []Task
}

func New() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) AddTask(task Task) {
	now := time.Now()
	nextRun := time.Date(now.Year(), now.Month(), now.Day(), task.hour, task.min, task.sec, 0, now.Location())
	if nextRun.Before(now) {
		nextRun = nextRun.Add(24 * time.Hour)
	}
	task.nextRun = nextRun
	s.tasks = append(s.tasks, task)
}

func (s *Scheduler) Add(name string, hour, min, sec int, taskFunc func()) {
	task := Task{
		hour:     hour,
		min:      min,
		sec:      sec,
		taskFunc: taskFunc,
	}
	s.AddTask(task)
}

func (s *Scheduler) Start() {
	go s.schedule()
}

func (s *Scheduler) schedule() {
	for {
		now := time.Now()
		closestRun := now.Add(-2 * time.Hour)
		for _, t := range s.tasks {
			if time.Now().After(t.nextRun) {
				// works async to dont block
				go t.taskFunc()
				fmt.Printf("Task %s triggered at %s\n", t.name, time.Now().Format("2006-01-02 15:04:05"))

				// Recalculate next run time
				now = time.Now()
				nextRun := time.Date(now.Year(), now.Month(), now.Day(), t.hour, t.min, t.sec, 0, now.Location())
				if nextRun.Before(now) {
					nextRun = nextRun.Add(24 * time.Hour)
				}
				t.nextRun = nextRun
				if now.After(closestRun) || nextRun.After(closestRun) {
					closestRun = nextRun
				}
			}

			sleepDuration := closestRun.Sub(now)
			time.Sleep(sleepDuration)
		}
	}
}
