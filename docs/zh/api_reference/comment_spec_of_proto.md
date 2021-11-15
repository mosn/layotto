
# proto文件注释规范

避免在注释符号`//`之间添加空行，否则生成工具protoc-gen-doc会生成格式错乱的文档。
一个坏示例:
```
message  BadCase{
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
message  GoodCase{
  // XXXXXXXX
  // XX
  // XXXXXX
  field A
}
```
或者你可以直接使用另一种注释符号:`/*  */`

假如你想添加一些注释在proto文件里，但不想让它们出现在生成的文档里，你可以在注释里使用`@exclude`前缀。
示例：只包括id字段的注释
```
/**
 * @exclude
 * This comment won't be rendered
 */
message ExcludedMessage {
  string id   = 1; // the id of this message.
  string name = 2; // @exclude the name of this message

  /* @exclude the value of this message. */
  int32 value = 3;
}
```
