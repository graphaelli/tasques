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

// NewReportOnClaimedTaskParams creates a new ReportOnClaimedTaskParams object
// with the default values initialized.
func NewReportOnClaimedTaskParams() *ReportOnClaimedTaskParams {
	var ()
	return &ReportOnClaimedTaskParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewReportOnClaimedTaskParamsWithTimeout creates a new ReportOnClaimedTaskParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewReportOnClaimedTaskParamsWithTimeout(timeout time.Duration) *ReportOnClaimedTaskParams {
	var ()
	return &ReportOnClaimedTaskParams{

		timeout: timeout,
	}
}

// NewReportOnClaimedTaskParamsWithContext creates a new ReportOnClaimedTaskParams object
// with the default values initialized, and the ability to set a context for a request
func NewReportOnClaimedTaskParamsWithContext(ctx context.Context) *ReportOnClaimedTaskParams {
	var ()
	return &ReportOnClaimedTaskParams{

		Context: ctx,
	}
}

// NewReportOnClaimedTaskParamsWithHTTPClient creates a new ReportOnClaimedTaskParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewReportOnClaimedTaskParamsWithHTTPClient(client *http.Client) *ReportOnClaimedTaskParams {
	var ()
	return &ReportOnClaimedTaskParams{
		HTTPClient: client,
	}
}

/*ReportOnClaimedTaskParams contains all the parameters to send to the API endpoint
for the report on claimed task operation typically these are written to a http.Request
*/
type ReportOnClaimedTaskParams struct {

	/*XTASQUESWORKERID
	  Worker ID

	*/
	XTASQUESWORKERID string
	/*ID
	  The id of the Task

	*/
	ID string
	/*NewReport
	  The request body

	*/
	NewReport *models.TaskNewReport
	/*Queue
	  The Queue of the Task

	*/
	Queue string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the report on claimed task params
func (o *ReportOnClaimedTaskParams) WithTimeout(timeout time.Duration) *ReportOnClaimedTaskParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the report on claimed task params
func (o *ReportOnClaimedTaskParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the report on claimed task params
func (o *ReportOnClaimedTaskParams) WithContext(ctx context.Context) *ReportOnClaimedTaskParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the report on claimed task params
func (o *ReportOnClaimedTaskParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the report on claimed task params
func (o *ReportOnClaimedTaskParams) WithHTTPClient(client *http.Client) *ReportOnClaimedTaskParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the report on claimed task params
func (o *ReportOnClaimedTaskParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithXTASQUESWORKERID adds the xTASQUESWORKERID to the report on claimed task params
func (o *ReportOnClaimedTaskParams) WithXTASQUESWORKERID(xTASQUESWORKERID string) *ReportOnClaimedTaskParams {
	o.SetXTASQUESWORKERID(xTASQUESWORKERID)
	return o
}

// SetXTASQUESWORKERID adds the xTASQUESWORKERId to the report on claimed task params
func (o *ReportOnClaimedTaskParams) SetXTASQUESWORKERID(xTASQUESWORKERID string) {
	o.XTASQUESWORKERID = xTASQUESWORKERID
}

// WithID adds the id to the report on claimed task params
func (o *ReportOnClaimedTaskParams) WithID(id string) *ReportOnClaimedTaskParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the report on claimed task params
func (o *ReportOnClaimedTaskParams) SetID(id string) {
	o.ID = id
}

// WithNewReport adds the newReport to the report on claimed task params
func (o *ReportOnClaimedTaskParams) WithNewReport(newReport *models.TaskNewReport) *ReportOnClaimedTaskParams {
	o.SetNewReport(newReport)
	return o
}

// SetNewReport adds the newReport to the report on claimed task params
func (o *ReportOnClaimedTaskParams) SetNewReport(newReport *models.TaskNewReport) {
	o.NewReport = newReport
}

// WithQueue adds the queue to the report on claimed task params
func (o *ReportOnClaimedTaskParams) WithQueue(queue string) *ReportOnClaimedTaskParams {
	o.SetQueue(queue)
	return o
}

// SetQueue adds the queue to the report on claimed task params
func (o *ReportOnClaimedTaskParams) SetQueue(queue string) {
	o.Queue = queue
}

// WriteToRequest writes these params to a swagger request
func (o *ReportOnClaimedTaskParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// header param X-TASQUES-WORKER-ID
	if err := r.SetHeaderParam("X-TASQUES-WORKER-ID", o.XTASQUESWORKERID); err != nil {
		return err
	}

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	if o.NewReport != nil {
		if err := r.SetBodyParam(o.NewReport); err != nil {
			return err
		}
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
