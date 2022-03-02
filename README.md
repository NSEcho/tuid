# tuid

tuid is a small tool to fetch id of twitter user using twitter developer API or to monitor sets of users to changes in their names or screen names(@something).

You need to create your own twitter developer application, and obtain the `Bearer` token. After you did that, you need to export it into `TUID_TOKEN` environment variable.

# Installation

```bash
$ go install github.com/lateralusd/tuid@latest
```

It will be located inside your go path, usually it is `~/go/bin`.

# Running

Create file containing screen names(@something) each on the newline and pass it to the tuid.

```bash
$ echo lateralusd_ > usersPath
$ TUID_TOKEN="..." ~/go/bin/tuid -users ./usersPath
2022/03/02 14:15:19 Read 1 users
2022/03/02 14:16:20 Name changed from lateralusd to lateralusdd
```
