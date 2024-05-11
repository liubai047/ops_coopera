package delayQueue

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

var delay = NewDelayQe[map[string]interface{}]()

func consumer() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	for {
		msg, err := delay.Watch(ctx)
		if err != nil {
			fmt.Println(fmt.Errorf("watch item error:%w", err).Error())
			return
		}
		v := msg.Content()
		fmt.Printf("%#v, %s\n", v, time.Now().Format(time.DateTime))
	}
}

func product(s string) {
	err1 := delay.Add(NewMessageItem(s+"k1", map[string]interface{}{
		"nowTime": s + "k1",
	}, time.Now().Add(10*time.Second)))

	err2 := delay.Add(NewMessageItem(s+"k2", map[string]interface{}{
		"nowTime": s + "k2",
	}, time.Now().Add(5*time.Second)))

	err3 := delay.Add(NewMessageItem(s+"k3", map[string]interface{}{
		"nowTime": s + "k3",
	}, time.Now().Add(15*time.Second)))

	err4 := delay.Add(NewMessageItem(s+"k4", map[string]interface{}{
		"nowTime": s + "k4",
	}, time.Now().Add(15*time.Second)))

	fmt.Printf("%#v\n", errors.Join(err1, err2, err3, err4))
}

func Test_Queue(t *testing.T) {
	fmt.Println(time.Now().Format(time.DateTime))
	go consumer()
	go product("1")
	go product("2")
	go product("3")
	go product("4")
	go product("5")
	go product("6")
	go product("7")
	go product("8")
	go product("9")
	go product("0")
	time.Sleep(time.Second * 60)
	fmt.Println(time.Now().Format(time.DateTime))
}
