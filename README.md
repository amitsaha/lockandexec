# `lockandexec` - Lock and execute

Get a lock and execute your program - local as well as "distributed" mode. 
Having a single running copy of your program can be a requirement in various situations. 
I am sure you have come up with one and that's why you have landed here. Regardless of 
the approach you end up with, you will have to implement some sort of a mutual exclusion
strategy. On a single compute instance (VM/container), you can get away with using
a keeping a PID file, Unix domain socket or the Linux only abstract Unix domain socket.
However, if you wanted to ensure a single copy of this program across multiple
compute instances separated by network boundaries, we need a different solution.

This repository contains the following programs to demonstrate the above.

## Locking using abstract unix domain sockets


## Locking using AWS DynamoDB
