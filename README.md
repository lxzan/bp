## BP

### Introduction

BP is the abbreviation of binary pool, built-in multiple memory pools, the size of the managed memory is 256, 512,
1024... bytes (default).

### Usage

```go
package main

import "github.com/lxzan/bp"

func main() {
	pool := bp.NewPool(128, 128*1024)
	println(pool.Get(1).Cap(), pool.Get(200).Cap(), pool.Get(600).Cap())
}
```
