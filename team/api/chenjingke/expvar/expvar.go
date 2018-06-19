package expvar

import (
	"encoding/json"
	"log"
	"os"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
)

/*
	expvar:
		公共变量包:辅助调试全局变量
	功能：
		1.它支持对变量的基本操作，修改、查询这些;
		2.整形类型，可以用来做计数器;
		3.操作都是线程安全的;
		4.此外还提供了调试接口，/debug/vars。它能够展示所有通过这个包创建的变量;
		5.所有的变量都是Var类型，可以自己通过实现这个接口扩展其它的类型;
		6.Handler()方法可以得到调试接口的http.Handler，和自己的路由对接;
*/
//	源码：

var (
	mutex   sync.RWMutex
	vars    = make(map[string]Var)
	varKeys []string // sorted
)

/*
	1.varKeys是全局变量所有的变量名，而且是有序的；
	2.vars根据变量名保存了对应的数据。当然mutex就是这个 Map 的锁；
	3.这三个变量组合起来其实是一个有序线程安全哈希表的实现。
*/

type Var interface {
	String() string
}

type Int struct {
	i int64
}

func (v *Int) Value() int64 {
	return atomic.LoadInt64(&v.i)
}

func (v *Int) String() string {
	return strconv.FormatInt(atomic.LoadInt64(&v.i), 10)
}

func (v *Int) Add(delta int64) {
	atomic.AddInt64(&v.i, delta)
}

func (v *Int) Set(value int64) {
	atomic.StoreInt64(&v.i, value)
}

/*
	1.这个包里面的所有类型都实现了这个接口；
	2.以 Int 类型举例。实现非常的简单，注意Add和Set方法是线程安全的。别的类型实现也一样
*/

func Publish(name string, v Var) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, existing := vars[name]; existing {
		log.Panicln("Reuse of exported var name:", name)
	}
	vars[name] = v
	varKeys = append(varKeys, name)
	sort.Strings(varKeys)
}

func NewInt(name string) *Int {
	v := new(Int)
	Publish(name, v)
	return v
}

/*
	1.Do方法，利用一个闭包，按照varKeys的顺序遍历所有全局变量；
	2.expvarHandler方法是http.Handler类型，将所有变量通过接口输出，里面通过Do方法，把所有变量遍历了一遍。挺巧妙；
	3.通过http.HandleFunc方法把expvarHandler这个外部不可访问的方法对外，这个方法用于对接自己的路由；
	4.输出数据的类型，fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)，可以发现，值输出的字符串，所以输出的内容是String()的结果。
		这里有一个技巧，虽然调用的字符串的方法，但是由于输出格式%s外面并没有引号，所有对于 JSON 来说，输出的内容是对象类型。
		相当于在 JSON 编码的时候做了一次类型转换。
*/

type Func func() interface{}

func (f Func) Value() interface{} {
	return f()
}

func (f Func) String() string {
	v, _ := json.Marshal(f())
	return string(v)
}

func cmdline() interface{} {
	return os.Args
}

/*
	1.它可以把任何类型转换成Var类型；
	2.Func定义的是函数，它的类型是func() interface{}
	3.Func(cmdline)，使用的地方需要看清楚，参数是cmdline而不是cmdline()，所以这个写法是类型转换。
		转换完之后cmdline方法就有了String()方法，在String()方法里又调用了f()，通过 JSON 编码输出。
		这个小技巧在前面提到的http.HandleFunc里面也有用到。
*/
