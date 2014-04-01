/*

Dynamo

This package contains objects needed to do request to DynamoDB.

GetItem

As defined: http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_GetItem.html.
However the "Key" attribute has been renamed to "Search".

    // build the request
    getItemRequest := dynamo.NewGetItemRequest()
    getItemRequest.AttributesToGet = []{"MyKey", "MyVal", "MyOtherVal"}
    getItemRequest.Search["MyKey"] = "Hello"
    getItemRequest.ConsistentRead = true
    getItemRequest.TableName = "my.table.name"
    getItemRequest.ReturnConsumedCapacity = dynamo.ConsumedCapacity_INDEXES
    // set the region
    getItemRequest.Host.Region = "us-west-2"
    // set the keys
    getItemRequest.Key, err = awsgo.GetSecurityKeys()
    // test err ...
    // do the request!
    resp, err := getItemRequest.Request()
    if err != nil {
        os.Exit(1)
    }
    if len(resp.Item) == 0 {
        // no items returned
        return
    }
    myKey := resp.Item["MyKey"].(string)
    myVal := resp.Item["MyVal"].(float64)
    myOtherVal := resp.Item["MyOtherVal"].([]string)

*/
package dynamo
