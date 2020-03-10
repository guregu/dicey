## dicey [![GoDoc](https://godoc.org/github.com/guregu/dicey?status.svg)](https://godoc.org/github.com/guregu/dicey)
`import "github.com/guregu/dicey"` 

dicey is a dice-rolling library for Go, using D&D-style formulas such as `2d6`, `1d8+3d2`, `1d10+5`, `3d6-2`, or `4d2+5d9-8`.

You can add or subtract any number of dice or static bonuses.

### Example

```go
import (
	"fmt"
	"math/rand"
	"time"

	"github.com/guregu/dicey"
)

func main() {
	// seed rand if you don't want consistent rolls
	rand.Seed(time.Now().UnixNano())

	dice := dicey.MustParse("4d20+2d6")
	fmt.Println("Random roll:", dice.Roll())
	fmt.Println("Maximum possible roll:", dice.Max()) // Prints: 92
	fmt.Println("Minimum possible roll:", dice.Min()) // Prints: 6
}
```

### License

BSD