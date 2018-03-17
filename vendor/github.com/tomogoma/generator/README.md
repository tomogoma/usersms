# generator
Random character generator that uses golang's crypto/rand package.

[godoc](https://godoc.org/github.com/tomogoma/generator)

## Usage

### Get the source
```bash
go get -u github.com/tomogoma/generator
```

### Import the package
```golang
import "github.com/tomogoma/generator"
```

### Generate random bytes from custom character set
```golang
charSet, _ := generator.NewCharSet(generator.AllChars)
randChars, _ := charSet.SecureRandomBytes(15)
```

### Generate random bytes from custom character set
```golang
r := generator.Random{}
randChars, _ := r.GenerateUpperCaseChars(6)
```
