# shallenge

This project tries to find the lowest SHA256 possible given a string based on the rules of https://shallenge.quirino.net/.

It uses goroutines to calculate different strings hash concurrently.

## String format

* {username}/{nonce}
    * username: 1-32 characters from a-zA-Z0-9_-
    * nonce: 1-64 characters from Base64 (a-zA-Z0-9+/)

The hash of the full string will be considered, not just the nonce

For example for the following string: "0x00cl/i5+8250U/Hi/HN/OinJeJj/eam3VRRI9d" the hash is "0000000003b98cfcf91b9b3e637354916d9cb960638373758d81072349bcb778"

## Compiling

```
$ go build
```

## Running

```
$ ./shallenge
```

### Running options

```
$ ./shallenge -h
  -h    Print this help message.
  -n int
        Number of workers to spawn (default 8)
  -p string
        String to add as suffix to the random generated word (default "i5+8250U/Hello+HN/")
  -v    Print the version number.
```