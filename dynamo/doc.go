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

PutItem

As defined: http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_PutItem.html

    // build the request
    putItem := dynamo.NewPutItemRequest()
    putItem.Item["MyKey"] = "Asd"
    putItem.Item["SomeNumbers"] = []float64{23.232,4342,112}
    putItem.Expected = dynamo.ExpectedItem{false, nil} // we don't expect it to exist
    putItem.TableName = "the.best.table"
    putItem.ReturnValues = dynamo.ReturnValues_ALL_NEW
    // set the region
    putItem.Host.Region = "us-west-2"
    // set the keys
    putItem.Key, err = awsgo.GetSecurityKeys()
    // test err ...
    // do the request!
    resp, err := putItem.Request()
    if err != nil {
        os.Exit(1)
    }
    if len(resp.BeforeAttributes) == 0 {
        // no items returned
        return
    }
    assert(resp.BeforeAttributes["MyKey"].(string) == "Asd")

DeleteItem

As defined: http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_DeleteItem.html

    // build the request
    deleteItem := dynamo.NewDeleteItemRequest()
    deleteItem.DeleteKey["MyKey"] = "Asd"
    deleteItem.TableName = "the.best.table"
    deleteItem.Expected = dynamo.ExpectedItem{true, "Asd"} // we expect it to exist
    deleteItem.ReturnValues = dynamo.ReturnValues_ALL_NEW
    // set the region
    deleteItem.Host.Region = "us-west-2"
    // set the keys
    deleteItem.Key, err = awsgo.GetSecurityKeys()
    // test err ...
    // do the request!
    resp, err := deleteItem.Request()
    if err != nil {
        os.Exit(1)
    }
    if len(resp.Attributes) == 0 {
        // no items returned
        return
    }
    assert(resp.Attributes["MyKey"].(string) == "Asd")

BatchWriteItem

As defined: http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_BatchWriteItem.html

    // build the request
    batchWrite := dynamo.NewBatchWriteItemRequest()
    batchWrite.AddPutRequest("test.table", map[string]interface{}{
        "MyKey" : "Asd",
        "SomeVal": []string{"gfdg", "Dfgd"},
        "SomeMore": 343.34,
    })
    batchWrite.AddDeleteRequest("worst.table", map[string]interface{}{
        "ItsKey" : 434334,
    })
    // set the region
    deleteItem.Host.Region = "us-west-2"
    // set the keys
    deleteItem.Key, err = awsgo.GetSecurityKeys()
    // test err ...
    // do the request!
    resp, err := deleteItem.Request()
    if err != nil {
        os.Exit(1)
    }
    if len(resp.UnprocessedItems) > 0 {
        // Need to make another request, since not all items were processed
        // For the next batchWrite2 set RequestItems = UnprocessedItems
    }


*/
package dynamo
