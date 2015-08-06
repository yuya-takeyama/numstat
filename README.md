# numstat

Show some statistics of input numbers.

## Installation

```
$ go get github.com/yuya-takeyama/numstat
```

## Usage

Given following data,

```
1
2
3
4
5
6
7
8
9
10
```

```
$ cat numbers.txt | numstat
Max: 10
Min: 1
Sum: 55
Avg: 5.5
```

### Show as JSON

```
$ cat numbers.txt | numstat --json
{"avg":5.5,"max":10,"min":1,"sum":55}
```

## License

The MIT License

## Author

Yuya Takeyama
