## Testdata

_The "testdata" folder holds various combinations of inputs and expected outputs to check that the library still works._

**You found a problematic HTML snippet? Great!**

Add your problematic HTML snippet to one of the `input.html` files in the `testdata` folder. Then run `go test -update` and have a look at which `.golden` output files changed in GIT.

_Note:_ Adding big block of HTML is not that helpful, so please:

1. narrow it down to the basic HTML structure that still causes the problem,
2. remove any unnecessary attributes (e.g. data-attributes) and
3. generalise the content (e.g. update links, replace text with lorem ipsum).

**=> Make sure it has been changed enough to be considered your own work!**

You can rerun `go test -update` until you are happy with the test case that you want to commit. Thanks for expanding the test cases!
