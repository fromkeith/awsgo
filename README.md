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
* Describe Table
* Get Item
* Put Item
* Query
* Scan
* Update Item
* Update Table

### S3
* Get Object
* Head Object
* Put Object

### SES
* Send Email
* Some helpers for dealing with the SNS notifications

### SQS
* Batch Send Message
* Change Message Visibility
* Delete Message
* Receive Message
* Send Message

### Cloud Search
* Document Batch Write


### EC2
* Describe Instances
* Get Instance Metadata

### SNS
* Publish

## SWF
* Poll For Activity Task
* Poll For Decision Task
* Record Activity Task Heartbeat
* Respond Activity Task Canceled
* Respond Activity Task Completed
* Respond Activity Task Failed
* Respond Decision Task Completed
* Start Workflow Execution
* A custom workflow execution helper, to make the incrmental decision process easier. See swf/swfhelper


## Depends on
    github.com/pmylund/sortutil
    - for sorting


## License
See LICENSE file.

