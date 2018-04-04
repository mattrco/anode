welford
=======

Package welford implements a one-pass algorithm for computing mean and
variance. Knuth attributes this algorithm to B.P. Welford in The Art
Of Computer Programming, Volume 2.

For an explanation of why you might want to calculate variance this
way, see http://www.johndcook.com/standard_deviation.html and
http://www.johndcook.com/blog/2008/09/26/comparing-three-methods-of-computing-standard-deviation/

For documentation on this package see
http://godoc.org/github.com/eclesh/welford

Quick start
===========

	$ go get github.com/eclesh/welford
	$ cd $GOPATH/src/github.com/eclesh/welford
	$ go test

License
=======

welford is licensed under the MIT license.
