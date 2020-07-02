package main
import(
	"fmt"
)

// 重写函数reverse，使用数组指针作为参数而不是slice

func main(){
	a := [...]int{0,1,2,3,4,5}
	reverse(&a)
	fmt.Println(a)
	b := []int{0,1,2,3,4,5}
	fmt.Println(rotate(b,2))
	s := []string{"ab","ab","bc"}
	no(s)

}

func reverse(s *[6]int){
	for i, j := 0, len(*s) - 1; i < j; i,j = i+1, j-1{
		(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
	}
}

// 编写一个rotate函数，通过一次循环完成旋转
func rotate(s []int, idx int)[]int{
	tmp := make([]int, len(s))
	for i := 0; i < len(s); i++{
		tmp[i] = s[(i+idx)%len(s)]
	}
	return tmp
}

// 写一个函数在原地完成消除[]string中相邻重复的字符串的操作
func no (s []string) []string{
	i := 1
	for idx, str := range s{
		if idx == 0{
			continue
		}
		if s[idx-1] == str{
			
		}
	}
	return s
}
