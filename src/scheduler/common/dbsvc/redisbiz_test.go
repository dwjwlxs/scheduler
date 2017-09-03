package dbsvc

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestSetex(t *testing.T) {
	host, _ := os.Hostname()
	_, err := Setex("testSetex", 3, host)
	if err != nil {
		t.Errorf("TestSetex() failed: %#v", err)
	}
}

func TestExists(t *testing.T) {
	key := "testExists2"
	host, _ := os.Hostname()
	_, _ = Setex(key, 60, host)
	e, err := Exists(key)
	if err != nil || e != 1 {
		t.Errorf("TestExists() failed: %#v", err)
	}
}

func TestPush(t *testing.T) {
	list := "testList1"
	ele := "anele_pushed"
	size, err := Push(list, ele)
	if err != nil || size < 1 {
		t.Errorf("TestPush() failed: %v", err)
	}
}

func TestPop(t *testing.T) {
	list := "testList"
	ele := "anele_testpop"
	_, _ = Push(list, ele)
	rele, err := Pop(list)
	if err != nil || rele != ele {
		t.Errorf("TestPop() failed: %#v", err)
	}
}

func TestSetnx(t *testing.T) {
	// now := time.Now()
	ok, err := Setnx("testsetnx"+fmt.Sprintf("%v", (time.Now()).Unix()), "123456")
	if err != nil || ok != 1 {
		t.Errorf("Test failed: %#v", err)
	}
}

func TestDel(t *testing.T) {
	key := "testDel"
	_, _ = Setnx(key, "123456")
	_, err := Del(key)
	if err != nil {
		t.Errorf("Test failed: %#v", err)
	}
}

func TestGet(t *testing.T) {
	key := "testGet"
	_, _ = Setnx(key, "123456")
	value, err := Get(key)
	if err != nil || value == "" {
		t.Errorf("Test failed: %#v", err)
	}
}

func TestGetSet(t *testing.T) {
	key := "testGetSet"
	_, _ = Setnx(key, "123456")
	value, err := GetSet(key, "test")
	if err != nil || value == "" {
		t.Errorf("Test failed: %#v", err)
	}
}
