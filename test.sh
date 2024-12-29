# The below command executes all the test scripts present in the module.
# The 3 dots are given to iterate over all the folders present inside lib and execute every test script present there.
# A test script is identified as a go file ending with "_test".
# The -v command line option prints all the verbose logs generated during test case execution.
# The -cover command line option prints the total code coverage metrics for each package.
go test ./lib/... -v -cover
