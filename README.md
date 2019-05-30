# Event

[![Discord](https://img.shields.io/discord/428990244952735764.svg?style=flat&logo=discord&colorB=green)](https://discord.gg/M9nxJ3g)
[![CircleCI branch](https://img.shields.io/circleci/project/github/syncaide/event/master.svg?label=master&logo=circleci)](https://circleci.com/gh/syncaide/workflows/event)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE.md)
[![](https://godoc.org/github.com/syncaide/event?status.svg&style=flat)](http://godoc.org/github.com/syncaide/event)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/syncaide/event.svg?style=flat&color=e36397&label=release)](https://github.com/syncaide/event/releases/latest)

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
	"github.com/syncaide/event"
)

func main() {
	mutex := sync.Mutex{}
    	results := make([]bool, 0)
    
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
