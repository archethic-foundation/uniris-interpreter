
# Uniris interpreter

Uniris interpreter used for smart contracts is built on top of golang and using syntax near of Javascript, Python and golang

## Operators

- Arithmetics: 
	- Addition: +
	- Substraction : - 
	- Multiplication: *
	- Division: /
- Comparison: 
    - Greater: >
    - Less : <
    - Greater and Equal: >=
    - Less and Equal: <=
    - Equal: ==
    - Different: !=
- Flow
    - if else
    ````
    if a > 10 {

    } 
    else {

    }
    ````
    - and
    ```
    if a == "10" and b == "20" {

    }
    ```
    - or
    ```
    if b == "20" or c == "20" {

    }
    ```
    - for
    ```
    for i=0; i < 10; i=i+1 {

    }
    ```
    - while
    ```
    while a < 10 {

    }
    ```


## Assignation

To assign variable no keyword are required

```
a=10
b="hello"
```
    

## Function

To create a function, *function* keyword is required
```
function hello(name) {
    return "Hello " + name
}


hello("John Doe")
```


## Print/Debug

To enable debugging, a *print* instruction is also available

```
print "hello"
a = 10
print a
```
## Natives functions

- now(): returns the current timestamp
