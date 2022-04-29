反射机制就是在运行时动态的调用对象的方法和属性，官方自带的reflect包就是反射相关的

1. reflect.TypeOf() 返回值Type类型

2. reflect.ValueOf() 返回值Value类型

3. value.Kind() 返回值Kind类型 注意与Type的不同

4. 修改反射对象，修改反射对象的前提条件是其值是可设置的
```
var a int = 10
v := reflect.ValueOf(&a)
e := v.Elem()
e.SetInt(15)
```

5. 遍历结构体字段内容
```
s := reflect.ValueOf(&student).Elem()
studentType := s.Type()
for i := 0; i < s.NumField(); i++ {
    f := s.Field(i)
    fmt.Printf("%d %s %s = %v\n", i, studentType.Field(i).Name, f.Type(), f.Interface())
}
```
