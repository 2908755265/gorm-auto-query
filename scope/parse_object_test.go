package scope

import (
	"fmt"
	"testing"
)

type testObj struct {
	name  string     `scope:"op=rlike"`
	age   int32      `scope:"op=eq"`
	ids   []int32    `scope:"op=in"`
	cash  float64    `scope:"name=userCash"`
	fnums *[]float64 `scope:"op=in"`
}

func TestParse(t *testing.T) {
	fns := make([]float64, 0)
	fns = append(fns, 1.0)
	to := testObj{
		name:  "张三", //
		age:   18,
		ids:   []int32{1}, //1, 2, 3
		fnums: nil,        //&fns,
	}
	cs, err := parseObject(to)
	if err != nil {
		panic(err)
	}
	for _, c := range cs {
		fmt.Printf("%s %s %v\n", c.key, c.op, c.val)
	}

	cash := 1.2

	to.name = "李四"
	to.cash = cash
	to.ids = append(to.ids, 1)
	fns = append(fns, cash)
	to.fnums = &fns
	cs, err = parseObject(&to)
	if err != nil {
		panic(err)
	}
	for _, c := range cs {
		fmt.Printf("%s %s %v\n", c.key, c.op, c.val)
	}
}

func TestAscii(t *testing.T) {
	println('a' - 'A')
}

func TestSnake(t *testing.T) {
	println(parseToSnake([]byte("TestFieldName")))
}

func BenchmarkParseObject(b *testing.B) {
	fns := make([]float64, 0)
	fns = append(fns, 1.2)
	for i := 0; i < b.N; i++ {
		to := testObj{
			name:  "张三", //
			age:   18,
			ids:   []int32{1}, //1, 2, 3
			cash:  1.1,
			fnums: &fns, //&fns,
		}
		_, err := parseObject(to)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkAdd(b *testing.B) {
	var a = 1
	for i := 0; i < b.N; i++ {
		a++
	}
}
