package copier_test

import (
	"errors"
	"github.com/jinzhu/copier"
	"reflect"
	"testing"
	"time"
)

type CustomFrom struct {
	Birthday *time.Time
}

type CustomTo struct {
	Birthday int64
}

// Convert from time.Time to int64
func fromTime2int64(to, from reflect.Value) (err error) {
	// log.Println(to.Addr().Type(), "->", from.Addr().Type())
	if _, ok := to.Addr().Interface().(*int64); ok {
		if fromTime, ok2 := from.Addr().Interface().(*time.Time); ok2 {
			var ts int64
			ts = fromTime.Unix()
			to.Set(reflect.Indirect(reflect.ValueOf(ts)))
		} else {
			err = errors.New("not from time.Time")
		}
	} else {
		err = errors.New("not to int64")
	}
	return err
}

// Convert from int64 to time.Time
func fromInt642Time(to, from reflect.Value) (err error) {
	// log.Println(to.Addr().Type(), "->", from.Addr().Type())
	if _, ok := to.Addr().Interface().(*time.Time); ok {
		if fromTimestamp, ok2 := from.Addr().Interface().(*int64); ok2 {
			t := time.Unix(*fromTimestamp, 0)
			to.Set(reflect.Indirect(reflect.ValueOf(t)))
		} else {
			err = errors.New("not from int64")
		}
	} else {
		err = errors.New("not to time.Time")
	}
	return err
}

func TestCopyCustomStruct(t *testing.T) {
	err := copier.RegisterCopyFunc(
		copier.CopierFunc{
			ToType:   reflect.TypeOf(int64(0)),
			FromType: reflect.TypeOf(time.Time{}),
			CopyFunc: fromTime2int64,
		},
		copier.CopierFunc{
			ToType:   reflect.TypeOf(time.Time{}),
			FromType: reflect.TypeOf(int64(0)),
			CopyFunc: fromInt642Time,
		},
	)
	if err != nil {
		t.Errorf("Failed to register copy functions.")
	}

	time1 := time.Now()
	f1 := CustomFrom{Birthday: &time1}
	t1 := CustomTo{}
	err = copier.Copy(&t1, &f1)
	if err != nil {
		t.Errorf("Error when attempting copy 1: " + err.Error())
	}
	if t1.Birthday != time1.Unix() {
		t.Errorf("Conversion failed from time.Time() to int64.")
	}

	time2 := time.Now()
	f2 := CustomFrom{}
	t2 := CustomTo{Birthday: time2.Unix()}
	err = copier.Copy(&f2, &t2)
	if err != nil {
		t.Errorf("Error when attempting copy 2: " + err.Error())
	}
	if f1.Birthday.Format("2006-01-02 03:04:05") != time2.Format("2006-01-02 03:04:05") {
		t.Errorf("Conversion failed from int64 to time.Time(): %s != %s", f1.Birthday.String(), time2.String())
	}
}
