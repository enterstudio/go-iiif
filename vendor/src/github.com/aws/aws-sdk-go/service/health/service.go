// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package health

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/private/protocol/jsonrpc"
)

// The AWS Health API provides programmatic access to the AWS Health information
// that is presented in the AWS Personal Health Dashboard (https://phd.aws.amazon.com/phd/home#/).
// You can get information about events that affect your AWS resources:
//
//    * DescribeEvents: Summary information about events.
//
//    * DescribeEventDetails: Detailed information about one or more events.
//
//    * DescribeAffectedEntities: Information about AWS resources that are affected
//    by one or more events.
//
// In addition, these operations provide information about event types and summary
// counts of events or affected entities:
//
//    * DescribeEventTypes: Information about the kinds of events that AWS Health
//    tracks.
//
//    * DescribeEventAggregates: A count of the number of events that meet specified
//    criteria.
//
//    * DescribeEntityAggregates: A count of the number of affected entities
//    that meet specified criteria.
//
// The Health API requires a Business or Enterprise support plan from AWS Support
// (http://aws.amazon.com/premiumsupport/). Calling the Health API from an account
// that does not have a Business or Enterprise support plan causes a SubscriptionRequiredException.
//
// For authentication of requests, AWS Health uses the Signature Version 4 Signing
// Process (http://docs.aws.amazon.com/general/latest/gr/signature-version-4.html).
//
// See the AWS Health User Guide (http://docs.aws.amazon.com/health/latest/ug/what-is-aws-health.html)
// for information about how to use the API.
//
// Service Endpoint
//
// The HTTP endpoint for the AWS Health API is:
//
//    * https://health.us-east-1.amazonaws.com
// The service client's operations are safe to be used concurrently.
// It is not safe to mutate any of the client's properties though.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/health-2016-08-04
type Health struct {
	*client.Client
}

// Used for custom client initialization logic
var initClient func(*client.Client)

// Used for custom request initialization logic
var initRequest func(*request.Request)

// Service information constants
const (
	ServiceName = "health"    // Service endpoint prefix API calls made to.
	EndpointsID = ServiceName // Service ID for Regions and Endpoints metadata.
)

// New creates a new instance of the Health client with a session.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a Health client from just a session.
//     svc := health.New(mySession)
//
//     // Create a Health client with additional configuration
//     svc := health.New(mySession, aws.NewConfig().WithRegion("us-west-2"))
func New(p client.ConfigProvider, cfgs ...*aws.Config) *Health {
	c := p.ClientConfig(EndpointsID, cfgs...)
	return newClient(*c.Config, c.Handlers, c.Endpoint, c.SigningRegion, c.SigningName)
}

// newClient creates, initializes and returns a new service client instance.
func newClient(cfg aws.Config, handlers request.Handlers, endpoint, signingRegion, signingName string) *Health {
	svc := &Health{
		Client: client.New(
			cfg,
			metadata.ClientInfo{
				ServiceName:   ServiceName,
				SigningName:   signingName,
				SigningRegion: signingRegion,
				Endpoint:      endpoint,
				APIVersion:    "2016-08-04",
				JSONVersion:   "1.1",
				TargetPrefix:  "AWSHealth_20160804",
			},
			handlers,
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

	// Run custom client initialization if present
	if initClient != nil {
		initClient(svc.Client)
	}

	return svc
}

// newRequest creates a new request for a Health operation and runs any
// custom request initialization.
func (c *Health) newRequest(op *request.Operation, params, data interface{}) *request.Request {
	req := c.NewRequest(op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(req)
	}

	return req
}
