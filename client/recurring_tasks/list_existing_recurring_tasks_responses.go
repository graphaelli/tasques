// Code generated by go-swagger; DO NOT EDIT.

package recurring_tasks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/lloydmeta/tasques/models"
)

// ListExistingRecurringTasksReader is a Reader for the ListExistingRecurringTasks structure.
type ListExistingRecurringTasksReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListExistingRecurringTasksReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListExistingRecurringTasksOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewListExistingRecurringTasksOK creates a ListExistingRecurringTasksOK with default headers values
func NewListExistingRecurringTasksOK() *ListExistingRecurringTasksOK {
	return &ListExistingRecurringTasksOK{}
}

/*ListExistingRecurringTasksOK handles this case with default header values.

OK
*/
type ListExistingRecurringTasksOK struct {
	Payload []*models.RecurringTask
}

func (o *ListExistingRecurringTasksOK) Error() string {
	return fmt.Sprintf("[GET /recurring_tasques][%d] listExistingRecurringTasksOK  %+v", 200, o.Payload)
}

func (o *ListExistingRecurringTasksOK) GetPayload() []*models.RecurringTask {
	return o.Payload
}

func (o *ListExistingRecurringTasksOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
