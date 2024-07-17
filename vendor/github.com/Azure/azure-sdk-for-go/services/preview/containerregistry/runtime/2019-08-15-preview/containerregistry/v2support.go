package containerregistry

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// V2SupportClient is the metadata API definition for the Azure Container Registry runtime
type V2SupportClient struct {
	BaseClient
}

// NewV2SupportClient creates an instance of the V2SupportClient client.
func NewV2SupportClient(loginURI string) V2SupportClient {
	return V2SupportClient{New(loginURI)}
}

// Check tells whether this Docker Registry instance supports Docker Registry HTTP API v2
func (client V2SupportClient) Check(ctx context.Context) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/V2SupportClient.Check")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.CheckPreparer(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerregistry.V2SupportClient", "Check", nil, "Failure preparing request")
		return
	}

	resp, err := client.CheckSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "containerregistry.V2SupportClient", "Check", resp, "Failure sending request")
		return
	}

	result, err = client.CheckResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerregistry.V2SupportClient", "Check", resp, "Failure responding to request")
		return
	}

	return
}

// CheckPreparer prepares the Check request.
func (client V2SupportClient) CheckPreparer(ctx context.Context) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"url": client.LoginURI,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithCustomBaseURL("{url}", urlParameters),
		autorest.WithPath("/v2/"))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CheckSender sends the Check request. The method will close the
// http.Response Body if it receives an error.
func (client V2SupportClient) CheckSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// CheckResponder handles the response to the Check request. The method always
// closes the http.Response Body.
func (client V2SupportClient) CheckResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}
