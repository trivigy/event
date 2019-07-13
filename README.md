# Event
[![CircleCI branch](https://img.shields.io/circleci/project/github/trivigy/event/master.svg?label=master&logo=circleci)](https://circleci.com/gh/trivigy/workflows/event)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE.md)
[![](https://godoc.org/github.com/trivigy/event?status.svg&style=flat)](http://godoc.org/github.com/trivigy/event)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/trivigy/event.svg?style=flat&color=e36397&label=release)](https://github.com/trivigy/event/releases/latest)

## Introduction
Event is a simple locking primitives which allows for sending notifications 
across goroutines when an event has occurred.

## Example
```go
package main

import (
	"fmt"
	"sync"
	"time"
	
	"github.com/pkg/errors"
	"github.com/trivigy/event"
)

func main() {
	mutex := sync.Mutex{}
    	results := make([]error, 0)
    
    	start := sync.WaitGroup{}
    	finish := sync.WaitGroup{}
    	
    	N := 5
    	start.Add(N)
    	finish.Add(N)
    	for i := 0; i < N; i++ {
    		go func() {
    			start.Done()
    			
    			value := event.Wait(nil)
    			mutex.Lock()
    			results = append(results, value)
    			mutex.Unlock()
    
    			finish.Done()
    		}()
    	}
    
    	start.Wait()
    	time.Sleep(100 * time.Millisecond)
    	event.Set()
    	finish.Wait()
    	
    	fmt.Printf("%+v\n", results)
}
```
