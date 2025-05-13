package main

import (
	"fmt"
	"time"
)

func main() {
	ts := &trainStoreLocal{}
	ss := &stationStoreLocal{}
	ls := &lineStoreLocal{}
	sch := &localScheduler{}
	register(ss, ls, ts, sch)
	doThings(ts, ss, ls, sch)
}

func register(
	sReg registrar[station],
	lReg registrar[line],
	tReg registrar[train],
	sch scheduler,
) {
	st1 := newStation(position{0, 0}, "Station 1", 3)
	st2 := newStation(position{10, 10}, "Station 2", 2)

	err := sReg.register(st1)
	nilErrOrPanic(err)

	err = sReg.register(st2)
	nilErrOrPanic(err)

	err = lReg.register(newLine(st1, st2))
	nilErrOrPanic(err)

	t1 := newTrain("Train 1", st1)
	err = tReg.register(t1)
	nilErrOrPanic(err)

	err = sch.add(newScheduleEntry(st2.id, t1.id, newExpectedTimes(NewNoneTime(), NewSomeTime(time.Now().Add(5*time.Minute)))))
	nilErrOrPanic(err)
}

func doThings(ts store[train], ss store[station], ls store[line], schViewer scheduleViewer) {
	{
		trains, err := ts.all()
		nilErrOrPanic(err)
		for _, t := range trains {
			fmt.Printf("%v", t)
		}
	}
	{
		stations, err := ss.all()
		nilErrOrPanic(err)
		for _, l := range stations {
			fmt.Println(l)
		}
	}
	{
		lines, err := ss.all()
		nilErrOrPanic(err)
		for _, l := range lines {
			fmt.Println(l)
		}
	}
	{
		schedules, err := schViewer.all()
		nilErrOrPanic(err)
		for _, s := range schedules {
			fmt.Println(s)
		}
	}
}

func nilErrOrPanic(err error) {
	if err != nil {
		panic(err)
	}
}
