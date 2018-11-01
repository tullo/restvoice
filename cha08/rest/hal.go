package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/rwirdemann/restvoice/cha05/domain"
)

type Link struct {
	Href string `json:"href"`
}

type Embedded struct {
	Bookings []domain.Booking `json:"bookings,omitempty"`
}

type HALInvoice struct {
	domain.Invoice
	Links    map[domain.Operation]Link `json:"_links"`
	Embedded *Embedded                 `json:"_embedded,omitempty"`
}

func NewHALInvoice(invoice domain.Invoice) HALInvoice {
	var links = make(map[domain.Operation]Link)
	links["self"] = Link{fmt.Sprintf("/invoice/%d", invoice.Id)}
	for _, o := range invoice.GetOperations() {
		if l, err := translate(o, invoice); err == nil {
			links[o] = l
		} else {
			log.Print(err)
		}
	}
	return HALInvoice{Invoice: invoice, Links: links}
}

func translate(operation domain.Operation, invoice domain.Invoice) (Link, error) {
	switch operation {
	case "book":
		return Link{fmt.Sprintf("/book/%d", invoice.Id)}, nil
	case "charge":
		return Link{fmt.Sprintf("/charge/%d", invoice.Id)}, nil
	case "payment":
		return Link{fmt.Sprintf("/payment/%d", invoice.Id)}, nil
	case "archive":
		return Link{fmt.Sprintf("/payment/%d", invoice.Id)}, nil
	default:
		return Link{}, errors.New(fmt.Sprintf("No translation found for operartion %s", operation))
	}
}

type HALInvoicePresenter struct {
	writer http.ResponseWriter
}

func NewHALInvoicePresenter(w http.ResponseWriter) HALInvoicePresenter {
	return HALInvoicePresenter{writer: w}
}

func (p HALInvoicePresenter) Present(i interface{}) {
	invoice := i.(HALInvoice)
	if len(invoice.Bookings) > 0 {
		invoice.Embedded = &Embedded{
			Bookings: invoice.Bookings,
		}
	}

	if b, err := json.Marshal(invoice); err == nil {
		p.writer.Write(b)
	}
}
