package rid_test

import (
	"fmt"

	"github.com/xiahuaxiahua0616/miniblog/internal/pkg/rid"
)

func ExampleResourceID_String() {
	// 定义一个资源标识符，例如用户资源
	userID := rid.UserID

	// 调用String方法，将ResourceID类型转换为字符串类型
	idString := userID.String()

	// 输出结果
	fmt.Println(idString)

	// Output:
	// user
}
