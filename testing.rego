package example

default allow = false

allow {
    test = hello("bob")
    test == "hello, " 
}