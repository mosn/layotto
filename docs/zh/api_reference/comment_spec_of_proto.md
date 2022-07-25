
# proto文件注释规范

避免在注释符号`//`之间添加空行，否则生成工具protoc-gen-doc会生成格式错乱的文档。
一个坏示例:

```
// XXXXXXXX
message BadCase{
  // XXXXXXXX
  //
  // XX
  //
  // XXXXXX
  field A
}
```

一个好示例:

```
// XXXXXXXX
message GoodCase{
  // XXXXXXXX
  // XX
  // XXXXXX
  field A
}
```
