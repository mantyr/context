# Set of contexts for all occasions

[![Build Status](https://travis-ci.org/mantyr/context.svg?branch=master)][build_status]
[![GoDoc](https://godoc.org/github.com/mantyr/context?status.png)][godoc]
[![Go Report Card](https://goreportcard.com/badge/github.com/mantyr/context?v=1)][goreport]
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](LICENSE.md)

## Installation

    $ go get github.com/mantyr/context


### WaitGroupContext returns a context with WaitGroup functions
1. You can add an element only if the context is not closed
2. You can delete the element at any time
3. You can close the context at any time
4. You can wait for the context to close
5. You can wait for the completion of operations related to the context


## Example

    package main

    import (
        "context"
        ct "github.com/mantyr/context"
    }

    func main() {
        ctx, cancel := ct.WaitGroupCancel(context.Background())
        ctx.Add(1)
        ctx.Add(1)
        ctx.Delete()
        ctx.Delete()
        cancel()

        select {
        case <- ctx.Done()
        case <- ctx.Wait()
        }
    }

## Author

[Oleg Shevelev][mantyr]

[mantyr]: https://github.com/mantyr


[build_status]: https://travis-ci.org/mantyr/context
[godoc]:        http://godoc.org/github.com/mantyr/context
[goreport]:     https://goreportcard.com/report/github.com/mantyr/context