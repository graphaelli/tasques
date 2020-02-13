// Code generated by go-swagger; DO NOT EDIT.

package tasks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/lloydmeta/tasques/models"
)

// NewMarkClaimedTaskFailedParams creates a new MarkClaimedTaskFailedParams object
// with the default values initialized.
func NewMarkClaimedTaskFailedParams() *MarkClaimedTaskFailedParams {
	var ()
	return &MarkClaimedTaskFailedParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewMarkClaimedTaskFailedParamsWithTimeout creates a new MarkClaimedTaskFailedParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewMarkClaimedTaskFailedParamsWithTimeout(timeout time.Duration) *MarkClaimedTaskFailedParams {
	var ()
	return &MarkClaimedTaskFailedParams{

		timeout: timeout,
	}
}

// NewMarkClaimedTaskFailedParamsWithContext creates a new MarkClaimedTaskFailedParams object
// with the default values initialized, and the ability to set a context for a request
func NewMarkClaimedTaskFailedParamsWithContext(ctx context.Context) *MarkClaimedTaskFailedParams {
	var ()
	return &MarkClaimedTaskFailedParams{

		Context: ctx,
	}
}

// NewMarkClaimedTaskFailedParamsWithHTTPClient creates a new MarkClaimedTaskFailedParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewMarkClaimedTaskFailedParamsWithHTTPClient(client *http.Client) *MarkClaimedTaskFailedParams {
	var ()
	return &MarkClaimedTaskFailedParams{
		HTTPClient: client,
	}
}

/*MarkClaimedTaskFailedParams contains all the parameters to send to the API endpoint
for the mark claimed task failed operation typically these are written to a http.Request
*/
type MarkClaimedTaskFailedParams struct {

	/*XTASQUESWORKERID
	  Worker ID

	*/
	XTASQUESWORKERID string
	/*Failure
	  The request body

	*/
	Failure *models.TaskFailure
	/*ID
	  The id of the Task

	*/
	ID string
	/*Queue
	  The Queue of the Task

	*/
	Queue string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the mark claimed task failed params
func (o *MarkClaimedTaskFailedParams) WithTimeout(timeout time.Duration) *MarkClaimedTaskFailedParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the mark claimed task failed params
func (o *MarkClaimedTaskFailedParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the mark claimed task failed params
func (o *MarkClaimedTaskFailedParams) WithContext(ctx context.Context) *MarkClaimedTaskFailedParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the mark claimed task failed params
func (o *MarkClaimedTaskFailedParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the mark claimed task failed params
func (o *MarkClaimedTaskFailedParams) WithHTTPClient(client *http.Client) *MarkClaimedTaskFailedParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the mark claimed task failed params
func (o *MarkClaimedTaskFailedParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithXTASQUESWORKERID adds the xTASQUESWORKERID to the mark claimed task failed params
func (o *MarkClaimedTaskFailedParams) WithXTASQUESWORKERID(xTASQUESWORKERID string) *MarkClaimedTaskFailedParams {
	o.SetXTASQUESWORKERID(xTASQUESWORKERID)
	return o
}

// SetXTASQUESWORKERID adds the xTASQUESWORKERId to the mark claimed task failed params
func (o *MarkClaimedTaskFailedParams) SetXTASQUESWORKERID(xTASQUESWORKERID string) {
	o.XTASQUESWORKERID = xTASQUESWORKERID
}

// WithFailure adds the failure to the mark claimed task failed params
func (o *MarkClaimedTaskFailedParams) WithFailure(failure *models.TaskFailure) *MarkClaimedTaskFailedParams {
	o.SetFailure(failure)
	return o
}

// SetFailure adds the failure to the mark claimed task failed params
func (o *MarkClaimedTaskFailedParams) SetFailure(failure *models.TaskFailure) {
	o.Failure = failure
}

// WithID adds the id to the mark claimed task failed params
func (o *MarkClaimedTaskFailedParams) WithID(id string) *MarkClaimedTaskFailedParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the mark claimed task failed params
func (o *MarkClaimedTaskFailedParams) SetID(id string) {
	o.ID = id
}

// WithQueue adds the queue to the mark claimed task failed params
func (o *MarkClaimedTaskFailedParams) WithQueue(queue string) *MarkClaimedTaskFailedParams {
	o.SetQueue(queue)
	return o
}

// SetQueue adds the queue to the mark claimed task failed params
func (o *MarkClaimedTaskFailedParams) SetQueue(queue string) {
	o.Queue = queue
}

// WriteToRequest writes these params to a swagger request
func (o *MarkClaimedTaskFailedParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// header param X-TASQUES-WORKER-ID
	if err := r.SetHeaderParam("X-TASQUES-WORKER-ID", o.XTASQUESWORKERID); err != nil {
		return err
	}

	if o.Failure != nil {
		if err := r.SetBodyParam(o.Failure); err != nil {
			return err
		}
	}

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	// path param queue
	if err := r.SetPathParam("queue", o.Queue); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}