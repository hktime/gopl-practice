package main
import(
	"fmt"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"hash"
)

// 编写一个函数，用于在默认情况下输出其标准输入的SHA256散列
// 但也支持一个输出SHA384或SHA512的命令行标记

var name string

func init(){
	// 参数处理，默认sha256
	flag.StringVar(&name, "name", "sha256", "The hashing function")
}

func main(){
	// 所有hash函数实现的公共接口
	var h hash.Hash
	flag.Parse()
	fmt.Printf("The hashing function is %s\n", name)
	switch name{
	case "sha256":
		h = sha256.New()
	case "sha384":
		h = sha512.New384()
	case "sha512":
		h = sha512.New()
	default:
		panic(name)
	}
	h.Write([]byte("hhhh"))
	fmt.Printf("%x\n", h.Sum(nil))
}

