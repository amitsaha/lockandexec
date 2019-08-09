package main

import "fmt"
import "net"
import "os"
import "os/exec"
import "syscall"
import "runtime"
import "time"

// https://golang.org/src/cmd/vendor/golang.org/x/sys/unix/syscall_linux.go?h=%40
// https://golang-examples.tumblr.com/post/92025745979/pitfall-of-abstract-unix-domain-socket-address-in
// ss -a -xp | grep myprogram
// https://utcc.utoronto.ca/~cks/space/blog/linux/SocketAbstractNamespace?showcomments#comments

// https://stackoverflow.com/questions/51151973/catching-bind-address-already-in-use-in-golang
func isAddressInUse(err error) bool {

	opErr, ok := err.(*net.OpError)
	if !ok {
		return false
	}
	sysCallErr, ok := opErr.Err.(*os.SyscallError)
	if !ok {
		return false
	}
	errErrno, ok := sysCallErr.Err.(syscall.Errno)
	if !ok {
		return false
	}
	if errErrno == syscall.EADDRINUSE {
		return true
	}
	const WSAEADDRINUSE = 10048
	if runtime.GOOS == "windows" && errErrno == WSAEADDRINUSE {
		return true
	}

	return false
}

func acquireLock() error {
	_, err := net.Listen("unix", "@lockandexec")
	return err
}

func main() {

	if len(os.Args) == 1 {
		fmt.Printf("Usage: lockandexec <prog args..>\n")
		os.Exit(1)
	}

	for {
		if err := acquireLock(); err != nil {
			if isAddressInUse(err) {
				fmt.Printf("Waiting to acquire lock\n")
				time.Sleep(30 * time.Second)

			} else {
				fmt.Print(err)
				os.Exit(1)
			}
		} else {
			fmt.Printf("Acquired lock..\n")
			break
		}

	}

	// do the work
	c := exec.Command(os.Args[1])
	for _, a := range os.Args[2:] {
		c.Args = append(c.Args, a)
	}

	outerr, err := c.CombinedOutput()
	if err != nil {
		fmt.Print(string(outerr))
		fmt.Print(err)
		os.Exit(1)
	} else {
		fmt.Printf("%s", outerr)
	}

	//once we exit, our "lock" is automatically released

}
