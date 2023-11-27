package main

import (
	"fmt"
	"time"
)

type Ttype struct {
	id         int
	cT         string
	fT         string
	taskRESULT []byte
}

func main() {
	taskCreator := func(a chan Ttype) {
		go func() {
			for {
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 {
					ft = "Some error occurred"
				}
				a <- Ttype{cT: ft, id: int(time.Now().Unix())}
				time.Sleep(time.Second) // Added sleep to prevent excessive task creation
			}
		}()
	}

	superChan := make(chan Ttype, 10)

	go taskCreator(superChan)

	taskWorker := func(a Ttype) Ttype {
		tt, _ := time.Parse(time.RFC3339, a.cT)
		if tt.After(time.Now().Add(-20 * time.Second)) {
			a.taskRESULT = []byte("task has been succeeded")
		} else {
			a.taskRESULT = []byte("something went wrong")
		}
		a.fT = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)

		return a
	}

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	taskSorter := func(t Ttype) {
		if string(t.taskRESULT) == "task has been succeeded" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
		}
	}

	go func() {

		for t := range superChan {
			t = taskWorker(t)
			go taskSorter(t)
		}
		close(doneTasks)
		close(undoneTasks)
	}()

	result := map[int]Ttype{}
	errs := []error{}

	go func() {
		for r := range doneTasks {
			result[r.id] = r
		}
		for r := range undoneTasks {
			errs = append(errs, r)
		}
	}()

	time.Sleep(time.Second * 3)

	println("Errors:")
	for _, r := range errs {
		println(r)
	}

	println("Done tasks:")
	for _, r := range result {
		println(r.id)
	}
}
