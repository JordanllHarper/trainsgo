package main

import (
	"fmt"
)

func main() {
	ts := &trainStoreLocal{}
	ss := &stationStoreLocal{}
	ls := &lineStoreLocal{}
	doThings(ts, ss, ls)
}

func doThings(ts trainStore, ss stationStore, ls lineStore) {
	st1 := newStation(position{0, 0}, "Station 1", 3)
	st2 := newStation(position{10, 10}, "Station 2", 2)

	err := ss.register(st1)
	nilErrOrPanic(err)

	err = ss.register(st2)
	nilErrOrPanic(err)

	err = ls.register(newLine(st1, st2))
	nilErrOrPanic(err)

	t1 := newTrain("Train 1", st1)
	err = ts.register(t1)
	nilErrOrPanic(err)

	{
		train, err := ts.getById(t1.id)
		nilErrOrPanic(err)
		fmt.Println(train)
	}
	{
		lines, err := ls.getAll()
		nilErrOrPanic(err)
		fmt.Println(lines)
	}
	{
		station, err := ss.getById(st1.id)
		nilErrOrPanic(err)
		fmt.Printf("%v\n", station.name)
	}
}

func nilErrOrPanic(err error) {
	if err != nil {
		panic(err)
	}
}
