package cas

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

type AttributesStyle int

const (
	CasThreeZeroNamedAttributesStyle AttributesStyle = iota
	CasThreeZeroAnyAttributesStyle
	RubyCasAttributesStyle
)

func (as AttributesStyle) String() string {
	switch as {
	case CasThreeZeroNamedAttributesStyle:
		return "CasThreeZeroNamedAttributesStyle"
	case CasThreeZeroAnyAttributesStyle:
		return "CasThreeZeroAnyAttributesStyle"
	case RubyCasAttributesStyle:
		return "RubyCasAttributesStyle"
	default:
		return ""
	}
}

type TestServer struct {
	serviceTickets map[string]*TestTicket
}

type TestTicket struct {
	Name                string
	Service             string
	Username            string
	ProxyGrantingTicket string
	Attributes          UserAttributes
	AttributesStyle     AttributesStyle
}

func (ts *TestServer) NewTicket(ticket string) *TestTicket {
	t := &TestTicket{
		Name:            ticket,
		Attributes:      make(UserAttributes),
		AttributesStyle: CasThreeZeroNamedAttributesStyle,
	}

	return t
}

func (ts *TestServer) AddTicket(ticket *TestTicket) {
	if ts.serviceTickets == nil {
		ts.serviceTickets = make(map[string]*TestTicket)
	}

	ts.serviceTickets[ticket.Name] = ticket
}

func (ts *TestServer) Close() {
	ts.serviceTickets = nil
}

func (ts *TestServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	switch r.URL.Path {
	case "/validate":
		ticket := query.Get("ticket")
		service := query.Get("service")

		if t, ok := ts.serviceTickets[ticket]; ok {
			if t.Service == service {
				fmt.Fprintf(w, "yes\n%v\n", t.Username)
			} else {
				fmt.Fprintf(w, "no\n\n")
			}
		} else {
			fmt.Fprintf(w, "no\n\n")
		}

		return
	case "/proxyValidate":
		ticket := query.Get("ticket")
		service := query.Get("service")

		var serviceResponse *xmlServiceResponse
		if t, ok := ts.serviceTickets[ticket]; ok {
			if t.Service == service {
				serviceResponse = successServiceResponse(t.Username, t.ProxyGrantingTicket)

				switch t.AttributesStyle {
				case CasThreeZeroNamedAttributesStyle:
					var userAttributes []*xmlNamedAttribute
					for key, values := range t.Attributes {
						for _, value := range values {
							userAttributes = append(userAttributes, &xmlNamedAttribute{
								Name:  key,
								Value: value,
							})
						}
					}

					serviceResponse.Success.Attributes = &xmlAttributes{
						AuthenticationDate:                     time.Now().UTC(),
						LongTermAuthenticationRequestTokenUsed: false,
						IsFromNewLogin:                         true,
					}

					if len(userAttributes) > 1 {
						serviceResponse.Success.Attributes.UserAttributes = &xmlUserAttributes{
							Attributes: userAttributes,
						}
					}
				case CasThreeZeroAnyAttributesStyle:
					var userAttributes []*xmlAnyAttribute
					for key, values := range t.Attributes {
						for _, value := range values {
							userAttributes = append(userAttributes, &xmlAnyAttribute{
								XMLName: xml.Name{Local: key},
								Value:   value,
							})
						}
					}

					serviceResponse.Success.Attributes = &xmlAttributes{
						AuthenticationDate:                     time.Now().UTC(),
						LongTermAuthenticationRequestTokenUsed: false,
						IsFromNewLogin:                         true,
					}

					if len(userAttributes) > 1 {
						serviceResponse.Success.Attributes.UserAttributes = &xmlUserAttributes{
							AnyAttributes: userAttributes,
						}
					}
				case RubyCasAttributesStyle:
					var userAttributes []*xmlAnyAttribute
					for key, values := range t.Attributes {
						for _, value := range values {
							userAttributes = append(userAttributes, &xmlAnyAttribute{
								XMLName: xml.Name{Local: key},
								Value:   value,
							})
						}
					}

					if len(userAttributes) > 1 {
						serviceResponse.Success.ExtraAttributes = userAttributes
					}
				}
			} else {
				serviceResponse = failureServiceResponse("INVALID_SERVICE", fmt.Sprintf("Ticket %s is not recognized for service %s", ticket, service))
			}
		} else {
			serviceResponse = failureServiceResponse("INVALID_TICKET", fmt.Sprintf("Ticket %s not recognized", ticket))
		}

		e := xml.NewEncoder(w)
		e.Indent("", "  ")
		e.Encode(serviceResponse)
		return
	default:
		http.NotFound(w, r)
		return
	}
}
