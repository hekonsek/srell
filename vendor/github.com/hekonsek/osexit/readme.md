# osexit

`osexit` is micro library for Go dedicated for exiting Linux processes in a systematic manner.

## Usage

To close Linux process in case of error use `osexit.ExitOnError` function:

```
import "github.com/hekonsek/osexit"

err := someFunction()
osexit.ExitOnError(err)
```

This function checks if error is not `nil`. Then exits process with error code 1 (Unix general error code) and informative
message if needed:

```
$ ./myapp
Something went wrong: ops!
$ echo $?
1
```

If error is `nil` program continues execution.