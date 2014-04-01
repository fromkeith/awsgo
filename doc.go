/*

awsgo

Awsgo is a wrapper around AWS API calls. It is divided into separate namespaces for each
AWS service. Check each package for their documentation.

This package holds the base/shared functionality between all requests. Most likely, from this package,
you will only be using a few methods, and now the AwsRequest object directly.


Credentials

Awsgo supports two 'automatic' methods for getting your security keys.

The first is via a json file.
THis is useful when testing locally. It expects a file named 'awskeys.json' to exist in your
working directory. This json should marshall to awsgo.Credentials struct.

The second is via the EC2 Security Role. This is done via a request to
http://169.254.169.254/latest/meta-data/iam/security-credentials. The EC2 Security Role
is only used if 'awskeys.json' does not exist.

In order to use these 'automatic' methods of Key population, you should call awsgo.GetSecurityKeys()


*/
package awsgo
