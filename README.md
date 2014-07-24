awsgo
=====

An _incomplete_ go library to talk to AWS services. http://godoc.org/github.com/fromkeith/awsgo

This is by no means a fully, complete, or well tested library. It will evolve as needed. Pull requests are more then welcome. As well as any go related tips or feedback.


## Services
Supported operations for each service:
> Caveat: Not every operation has had each option fully tested. See examples for some minor tests / uses.

### Cloud Watch
* Put Metric
* Get Metric Statistics
* Logs:
    * Create Log Group
    * Create Log Stream
    * Get Log Events
    * Put Log Events
    * Put Rentention Policy

### DynamoDB
> Caveat: Binary values are not yet supported.

Godoc: http://godoc.org/github.com/fromkeith/awsgo/dynamo

* Batch Get Item
* Batch Write Item
* Delete Item
* Get Item
* Put Item
* Query
* Scan
* Update Item

### S3
* Put Object
* Get Object

### SES
* Send Email
* Some helpers for dealing with the SNS notifications

### SQS
* Change Message Visibility
* Delete Message
* Receive Message
* Send Message

### Cloud Search
* Document Batch Write


### EC2
* Describe Instances
* Get Instance Metadata

## Depends on
    github.com/pmylund/sortutil
    - for sorting


## License
See LICENSE file.

