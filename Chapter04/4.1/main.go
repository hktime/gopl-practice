package main
import(
	"fmt"
	"crypto/sha256"
)

// 编写一个函数，用于统计SHA256散列中不同的位数

func main(){
	c1 := sha256.Sum256([]byte("x"))
	c2 := sha256.Sum256([]byte("X"))
	fmt.Println(popCount(c1,c2))
}

func popCount(c1, c2 [32]uint8)int{
	var count int = 0
	for i := 0; i < 32; i++{
		x, y := c1[i], c2[i]
		var one uint8 = 1
		for j := 0; j  < 8; j++{
			if x&one == y & one{
				count++
			}
			one <<= 1
		}
	}
	return count
}