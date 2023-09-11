# Simple Data Store

[![Go Report Card](https://goreportcard.com/badge/schafer14/sds?style=flat-square)](https://goreportcard.com/report/schafer14/sds)
[![Go Reference](https://pkg.go.dev/badge/schafer14/sds.svg)](https://pkg.go.dev/schafer14/sds)
[![LICENSE](https://img.shields.io/github/license/schafer14/sds.svg?style=flat-square)](https://github.com/schafer14/sds/blob/master/LICENSE)

Simple Data Store (sds) is a simple storage interface I use that allows me 
to **focus on my domain model** instead of writing database logic. Additionally, 
this interface **supports my testing workflows** by providing different storage 
backends for different types of tests. This means I never need to mock a database
for a test, I always have the real thing. 

The data interface is small and opinionated. The goal is to support 90% of 
database needs, and allow the more complicated data models to be implemented 
with a more sophisticated library. 

I have been building implementations for this interface as the need arises. 
Currently there are implementations for an in memory data base, bbolt, and
mongodb. If you are interested in adding another database there is a test 
suite to ensure your implementation works correctly. I plan to write implementations
for Firestore and Dynamo in the future, but only as the need arises. If you 
would like a database to be supported please leave a comment.

I recommending quickly reading the "how to" notes for any implementation you 
plan to use. They are in the [./wiki/how-tos](how to) directory.

## More Information 

- [./wiki/getting-started.md](Getting started guide)
- [./wiki/how-tos](How to documentation)
- [https://pkg.go.dev/schafer14/sds](API reference)
- [./wiki/explantions](Explanations)

## Feedback

Please feel free to raise issues with any feedback. Questions, suggestions, 
improvements, and comments are all welcome.



