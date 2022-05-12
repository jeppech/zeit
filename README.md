# zeit
A simple package for working with time, not timestamps


```golang
package main

import (
  "fmt"
  "github.com/jeppech/zeit"
)

func main() {
  now := time.Now().UTC()
  fmt.Println(now) // 2022-05-06 07:19:35.072174 +0000 UTC

  z_now := zeit.Now() // 07:19:35
  fmt.Println(z_now)


}
```
