package dingtalk

import (
	"encoding/json"
	"errors"
	"reflect"

	"strings"
)

const (
	ERROR_FAILED_TO_PARSE = "failed to parse Dingtalk ProcessForm Data!"
	ERROR_REQUEST_FAILED  = "Dingtalk ProcessForm Request failed!"
)

type DingtalkProcessFormParser struct {
	components map[string]string
	formdata   *DingtalkProcessFormData
}

type DingtalkProcessFormData struct {
	// 审批发起人的部门ID
	OriginatorDeptId string `json:"originator_dept_id"`
	// 审批发起人的部门名称
	OriginatorDeptName string `json:"originator_dept_name"`
	// 审批发起人的钉钉用户ID
	OriginatorUserId string `json:"originator_userid"`
	// 审批结果
	Result string `json:"result"`
	// 审批状态
	Status string `json:"status"`
	// 审批标题
	Title string `json:"title"`
	// 回调地址
	Callback string `json:"callback"`
}

func NewDingtalkProcessFormParser() *DingtalkProcessFormParser {
	formparser := new(DingtalkProcessFormParser)
	formparser.components = make(map[string]string)
	formparser.formdata = new(DingtalkProcessFormData)
	return formparser
}

func getValueFloat64(v map[string]interface{}, fildname string) (float64, error) {
	if value, ok := v[fildname].(float64); ok {
		return value, nil
	} else {
		return -1, errors.New("type is mismatch!")
	}
}

func (parser *DingtalkProcessFormParser) GetValue(fieldname string) string {
	return parser.components[fieldname]
}

func (parser *DingtalkProcessFormParser) getComponentsValue(v interface{}) {
	form_component_values := v.(map[string]interface{})["form_component_values"]
	if reflect.ValueOf(form_component_values).Kind() == reflect.Invalid {
		return
	}
	nums := reflect.ValueOf(form_component_values).Len()
	for i := 0; i < nums; i++ {
		x := reflect.ValueOf(form_component_values).Index(i).Interface()
		name := reflect.ValueOf(x).MapIndex(reflect.ValueOf("name"))

		if name.Kind() == reflect.Invalid {
			callback_url := reflect.ValueOf(x).MapIndex(reflect.ValueOf("value"))
			parser.formdata.Callback = reflect.ValueOf(callback_url.Interface()).String()
			if strings.Index(parser.formdata.Callback, "__callback") == 0 {
				parser.formdata.Callback = parser.formdata.Callback[10:]
			} else {
				parser.formdata.Callback = ""
			}
		} else {
			value := reflect.ValueOf(x).MapIndex(reflect.ValueOf("value"))
			field_name := reflect.ValueOf(name.Interface()).String()
			parser.components[field_name] = reflect.ValueOf(value.Interface()).String()
		}
	}
}

func getValueString(v interface{}, fieldname string) string {
	r := reflect.ValueOf(v).MapIndex(reflect.ValueOf(fieldname))
	if r.Kind() == reflect.Invalid {
		return ""
	} else {
		return reflect.ValueOf(r.Interface()).String()
	}
}

func (parser *DingtalkProcessFormParser) processMapData(t map[string]interface{}) (*DingtalkProcessFormData, error) {
	var (
		process_instance map[string]interface{}
		formdata         *DingtalkProcessFormData
	)
	if errcode, err := getValueFloat64(t, "errcode"); err != nil {
		return nil, errors.New(ERROR_FAILED_TO_PARSE)
	} else {
		if errcode != 0 {
			return nil, errors.New(ERROR_REQUEST_FAILED)
		}
	}

	if t["process_instance"] == nil {
		return nil, errors.New(ERROR_FAILED_TO_PARSE)
	}
	if reflect.TypeOf(t["process_instance"].(map[string]interface{})).Kind() != reflect.Map {
		return nil, errors.New(ERROR_FAILED_TO_PARSE)
	}
	process_instance = t["process_instance"].(map[string]interface{})

	formdata = parser.formdata
	formdata.OriginatorDeptId = getValueString(process_instance, "originator_dept_id")
	formdata.OriginatorDeptName = getValueString(process_instance, "originator_dept_name")
	formdata.OriginatorUserId = getValueString(process_instance, "originator_userid")
	formdata.Result = getValueString(process_instance, "result")
	formdata.Status = getValueString(process_instance, "status")
	formdata.Title = getValueString(process_instance, "title")

	parser.getComponentsValue(process_instance)
// 	fmt.Println(parser.components)
	return formdata, nil
}

func (parser *DingtalkProcessFormParser) Parse(v interface{}) (*DingtalkProcessFormData, error) {
	var (
		x   map[string]interface{}
		err error
	)
	switch t := v.(type) {
	case string:
		if err = json.Unmarshal([]byte(t), &x); err != nil {
			return nil, errors.New(ERROR_FAILED_TO_PARSE)
		}
		return parser.processMapData(x)
	case *string:
		if err = json.Unmarshal([]byte(*t), &x); err != nil {
			return nil, errors.New(ERROR_FAILED_TO_PARSE)
		}
		return parser.processMapData(x)
	case map[string]interface{}:
		return parser.processMapData(t)
	default:
		return nil, errors.New(ERROR_FAILED_TO_PARSE)
	}
}
