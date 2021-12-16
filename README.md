# 说明
钉钉审批表单解析

# 使用方法

```go
import(
	"github.com/vber/dingtalk_processform_parser"
)

func main() {
	var (
		v interface{}
	)

  // 直接传递字符串方式解析
	p := dingtalk.NewDingtalkProcessFormParser()
	d, err := p.Parse(&data)
	if err != nil {
		fmt.Println(err)
	} else {
    fmt.Println(parser.GetValue("开始日期"), data.Callback)
		fmt.Println(formdata.OriginatorDeptId, formdata.OriginatorDeptName, formdata.Status, formdata.Title, formdata)
	}

  // 直接传递map方式解析
	if err := json.Unmarshal([]byte(data), &v); err == nil {
		parser := dingtalk.NewDingtalkProcessFormParser()
		if data, err := parser.Parse(v); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(parser.GetValue("开始日期"), data.Callback)
			fmt.Println(formdata.OriginatorDeptId, formdata.OriginatorDeptName, formdata.Status, formdata.Title, formdata)
		}
	}
}
```
