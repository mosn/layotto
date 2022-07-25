
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
