### Probably

Probably is a Bloom filter implementation for Golang.
It's quite simple to use and it doesn't have any external dependency.

#### Usage

To create a Bloom filter it's required to specify the size of the set and a desired false
positive rate.

```go
import (
    "github.com/probably/filters"
)

func main() {
    filter := filters.NewBloomFilter(10, 0.01)
}

```
After this it's possible to add values and check if they are in the set like that
```go
filter.Add([]byte("hello"))
filter.Contains([]byte("world"))
```

It's possible to join two Bloom filters with the same configurataion: same size and false positive rate.
Probably provides two methods for this, `Merge` and `Union`.
Merge joins two Bloom filters and stores the result in the first Bloom filter.
```go
filter1 := filters.NewBloomFilter(10, 0.01)
filter2 := filters.NewBloomFilter(10, 0.01)

// add values here

filter1.Merge(filter2)
```
Whereas Union joins two Bloom filters but returns the new Bloom filter as a result.
```go
filter1 := filters.NewBloomFilter(10, 0.01)
filter2 := filters.NewBloomFilter(10, 0.01)

// add values here

result := filter1.Union(filter2)
```
To reset a filter state, call `Clear` method on the filter
```go
filter.Clear()
```

#### Implementation details

As an underlying data structure, Probably uses a bit array that is implemented using
byte slices. So in case a user wants to allocate 9 bits then Probably will create
a 2 byte slice to hold the data in the filter.

To generate k hashes Probably uses only two hash functions: MD5 and SHA1.
For more information about this, please refer to this [amazing paper](https://www.eecs.harvard.edu/~michaelm/postscripts/tr-02-05.pdf) by Adam Kirsch.

#### Future plans
It would be nice to extend Probably with other probabilistic data structures like HyperMinHash or Cuckoo filter

#### Useful links

- [Building a Better Bloom filter (paper) by Adam Kirch](https://www.eecs.harvard.edu/~michaelm/postscripts/tr-02-05.pdf)
- [Slides from The Univesity of Texas at Austin](https://www.cs.utexas.edu/users/lam/396m/slides/Bloom_filters.pdf)
