awsgo
=====

An _incomplete_ go library to talk to AWS services. 

This is by no means a fully, complete, or well tested library. It will evolve as needed. Pull requests are more then welcome. As well as any go related tips or feedback.


## Services
Supported operations for each service:
> Caveat: Not every operation has had each option fully tested. See examples for some minor tests / uses.

### Cloud Watch
* Put Metric

### DynamoDB
> Caveat: Only single value items are returned. Eg. If the column 'A' has 2 values for one item, that is not supported yet. (At least, not yet tested)

* Batch Get Item
* Get Item
* Put Item
* Update Item

### S3
* Put Item

### SES
* Send Email

### SQS
* Delete Message
* Receive Message
* Send Message


## Depends on
    github.com/pmylund/sortutil
    - for sorting


## License
See LICENSE file.

