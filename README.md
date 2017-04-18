# Goof

Goof lets you call functions in your binary with just the string of their
name. How?

```go
var troop goof.Troop
out, err := troop.Call("fmt.Printf", "hello %s", []interface{}{"world"})
if err != nil { // some error calling the function
	return err
}
n, err := out[0].(int), out[1].(error)
if err != nil {
	return err
}
fmt.Println("wrote", n, "bytes")
```
