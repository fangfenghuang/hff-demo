Json(Javascript Object Nanotation)是一种数据交换格式，常用于前后端数据传输。
go语言本身为我们提供了json的工具包”encoding/json”


# 定制JSON序列化
Json Marshal：将数据编码成json字符串
```
    jsonStu, err := json.Marshal(stu)
```
json.Unmarshal: 反序列的函数，将json字符串解码到相应的数据结构；
```
    jsonstr := `{"a":"aaa","b":"bbb"}`
    map1 := make(map[string]interface{})
    err := json.Unmarshal([]byte(jsonstr), &map1)
```
## 临时忽略struct空字段： omitempty
```
json.Marshal(struct {
    *User
    Password bool `json:"password,omitempty"`  //忽略该字段
    }{
        User: user,
    })
```

## 忽略一些字段: -
```
json.Marshal(struct {
    *User
    Password bool `json:"-"`
}{
    User: user,
})
```

## 临时添加额外的字段: 
```
json.Marshal(struct {
    *User
    Token    string `json:"token"`  //添加该字段
    Password bool `json:"password,omitempty"`
    }{
        User: user,
        Token: token,
    })
```

## 临时粘合两个struct
```
type BlogPost struct {
    URL   string `json:"url"`
    Title string `json:"title"`
}
type Analytics struct {
    Visitors  int `json:"visitors"`
    PageViews int `json:"page_views"`
}
json.Marshal(struct{
    *BlogPost
    *Analytics
}{post, analytics})
```

## 一个json切分成两个struct: todo

## 用字符串传递数字: `json:",string"`
```
type TestObject struct {
    Field1 int    `json:",string"`
}
```



