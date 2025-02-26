package os

import realOS "os"

// ReadFile reads the named file and returns the contents.
// A successful call returns err == nil, not err == EOF.
// Because ReadFile reads the whole file, it does not treat an EOF from Read
// as an error to be reported.
var ReadFile = realOS.ReadFile

// Stat returns a [FileInfo] describing the named file.
// If there is an error, it will be of type [*PathError].
var Stat = realOS.Stat

// IsNotExist returns a boolean indicating whether its argument is known to
// report that a file or directory does not exist. It is satisfied by
// [ErrNotExist] as well as some syscall errors.
//
// This function predates [errors.Is]. It only supports errors returned by
// the os package. New code should use errors.Is(err, fs.ErrNotExist).
var IsNotExist = realOS.IsNotExist

// Create creates or truncates the named file. If the file already exists,
// it is truncated. If the file does not exist, it is created with mode 0o666
// (before umask). If successful, methods on the returned File can
// be used for I/O; the associated file descriptor has mode O_RDWR.
// The directory containing the file must already exist.
// If there is an error, it will be of type *PathError.
var Create = realOS.Create
