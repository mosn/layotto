
# Comment specification of proto file

Avoid adding empty lines between comments  symbols `//`.If there is a blank line in the comments, the tool(protoc-gen-doc) will generate malformed documents.
bad case:

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

good case:

```
// XXXXXXXX
message GoodCase{
  // XXXXXXXX
  // XX
  // XXXXXX
  field A
}
```

Or you can use another annotation symbol directly `/*  */`


If you want to have some comment in your proto files, but don't want them to be part of the docs, you can simply prefix the comment with `@exclude`.
Example: include only the comment for the id field

Attention: `/*  */` comments do not count towards passing ci `Proto Validation`. [refence](https://docs.buf.build/lint/rules#comments)

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