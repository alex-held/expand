<h1 align="center">Welcome to expand 👋</h1>
<p>
  <a href="https://twitter.com/0_alexheld" target="_blank">
    <img alt="Twitter: 0_alexheld" src="https://img.shields.io/twitter/follow/0_alexheld.svg?style=social" />
  </a>

</p>


![Go_Version](https://img.shields.io/github/go-mod/go-version/alex-held/expand?style=flat-square)
[![Go Reference](https://pkg.go.dev/badge/github.com/alex-held/expand.svg)](https://pkg.go.dev/github.com/alex-held/expand)


> Expands variables recursivly using a dependency graph algorithm

### 🏠 [Homepage](https://github.com/alex-held/expand)

## Install

```sh
go get github.com/alex-held/expand
```

## Usage

```go
func main() {
	vars := map[string]string{
	    "a": "$HOME/$b/$c",
		"b": "$c/foo",
		"c": "bar",
		"HOME": "/home/user" // override environment variable
    }   
	
	expansions, err := expand.Expand(vars)
	if err != nil {
	   panic(err)
	}
	
	println(expansions.MustGet("a"))
	
	// Outputs:
	// /home/user/bar/foo/bar
}
```

## Run tests

```sh
go test -v ./...
```

## Author

👤 **Alexander Held**

* Website: https://alexheld.io
* Twitter: [@0_alexheld](https://twitter.com/0\alexheld)
* Github: [@alex-held](https://github.com/alex-held)

## Show your support

Give a ⭐️ if this project helped you!
