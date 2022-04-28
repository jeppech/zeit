# zeit
A simple package for working with time, not timestamps


```golang
package main

import (
  "fmt"
  "github.com/jeppech/zeit"
)

func main() {
  now := time.Now()
  fmt.Println(now)

  z_now := zeit.Now()
  fmt.Println(z_now)

  
}
```
