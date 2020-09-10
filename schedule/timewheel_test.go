package schedule

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeWheel_AddTask(t *testing.T) {
	tw := New(1*time.Second, 180)
	tw.Start()
	tw.AddTask(3*time.Second, 2, "", TaskData{"a": "b"}, job)
	tw.AddForeverTaskBySeconds(3, "test3", TaskData{"a": "b"}, job2)
	tw.AddTaskByCron("0/2 * * * * *", "test2", nil, job2)
	select {}

}

func TestTimeWheel_AddJobByCron(t *testing.T) {
	tw := NewSecondWheel(180)
	tw.Start()
	tw.AddTaskByCron("0/3 * * * * *", "test2", nil, job2)
	select {}

}

func job(t TaskData) {
	fmt.Println("test1")
}

func job2(t TaskData) {

	now := time.Now()
	fmt.Println("test2->", now.Second(), ":", now.String(), ">", t)
}
