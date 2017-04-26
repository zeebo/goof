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
fmt.Printf("wrote %d bytes\n", n)
```

Caveat: you have to have called `fmt.Printf` elsewhere in your binary.

Goof lets you get access to globals in your binary with just the string of
their name. How?

```go
var troop goof.Troop
rv, err := troop.Global("net/http.DefaultServeMux")
if err != nil { // couldn't find it
	return err
}
// rv contains an addressable reflect.Value of the default ServeMux!
```

Caveat: the global must be used elsewhere in the binary somehow.

Goof lets you get access to all of the `reflect.Type`s in your binary. How?

```go
var troop goof.Troop
types, err := troop.Types()
if err != nil { // something went wrong getting them
	return err
}
for _, typ := range types {
	fmt.Println(typ)
}
```

Caveat: the types must be possible outputs to `reflect.TypeOf(val)` in your binary.

## Usage

You should probably just make a single `Troop` in your binary and use that
everywhere since it does a lot of caching and work on first use.

## How?

It loads up the dwarf information of any binary it's loaded in and then does
a bunch of unsafe tom foolery to perform these dirty deeds. How unsafe is it?

- Reusing needles unsafe.
- Jumping into a shark tank with a steak swimming suit unsafe.
- Carnival ride unsafe.
- Driving on the wrong side of the highway blindfolded unsafe.

## Should I use this?

Do you really have to ask? OF COURSE! If you do, please let me know what terrible
idea this enabled. I'm very interested.

## Testimonials

> "I can't wait to get some goof in my [manhole](https://github.com/jtolds/go-manhole)!" - [@jtolds](https://github.com/jtolds)

> "README is hilarious :joy:"

> "Now I just need to come up with something horrendously risky to use this for..."
