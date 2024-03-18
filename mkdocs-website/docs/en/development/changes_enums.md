#### Enums

In Go, enums are often defined as a type and a set of constants. For example:

```go
type MyEnum int

const (
    MyEnumOne MyEnum = iota
    MyEnumTwo
    MyEnumThree
)
```

Due to incompatibility between Go and JavaScript, custom types cannot be used in
this way. The best strategy is to use a type alias for float64:

```go
type MyEnum = float64

const (
    MyEnumOne MyEnum = iota
    MyEnumTwo
    MyEnumThree
)
```

In Javascript, you can then use the following:

```js
const MyEnum = {
  MyEnumOne: 0,
  MyEnumTwo: 1,
  MyEnumThree: 2,
};
```

- Why use `float64`? Can't we use `int`?
    - Because JavaScript doesn't have a concept of `int`. Everything is a
      `number`, which translates to `float64` in Go. There are also restrictions
      on casting types in Go's reflection package, which means using `int` doesn't
      work.