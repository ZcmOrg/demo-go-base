package errors

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego"
)

func main() {
	/*
		package errors:实现了创建错误值的函数
		New方法 :将字符串 text 包装成一个 error 对象返回
			func New(text string) error {
				return &errorString{text}
			}
	*/
	// 例子：
	err := errors.New("标的类型有误")
	if err != nil {
		beego.Info(err)
	}
	// 延伸1：fmt包有个返回error类型的：内部也是调用errors包的New方法
	/*
		func Errorf(format string, a ...interface{}) error {
			return errors.New(Sprintf(format, a...))
		}
	*/
	const name, id = "bimmler", 17
	err1 := fmt.Errorf("user %q (id %d) not found", name, id)
	if err1 != nil {
		fmt.Println(err1)
	}
	// 延伸2：sql调用
	/*
		o := orm.NewOrm()
		sql := ` select * from cg_product where product_code = ? `
		err = o.Raw(sql, product_code).QueryRow(&cgPro)
		若cgPro值为空，err会返回错误<QuerySeter> no row found，另外orm.ErrNoRows的值就是errors.New("<QuerySeter> no row found")


		o := orm.NewOrm()
		sql := ` select * from cg_trade where trade_uuid =?  `
		_, err = o.Raw(sql, tradeUuid).QueryRows(&cgtrades)
		若cgtrades值为空，err不会返回错误，是返回nil
	*/
}
