package cas

import (
	"encoding/xml"
	"errors"
	"time"
)

var (
	ErrAuthenticationDateUnknown                     = errors.New("cas: authentication date not included in service response")
	ErrLongTermAuthenticationRequestTokenUsedUnknown = errors.New("cas: long term authentication request token used flag not included in service response")
	ErrIsFromNewLoginUnknown                         = errors.New("cas: is from new login flag not included in service response")
	ErrMemberOfUnknown                               = errors.New("cas: member of not included in service response")
	ErrUserAttributesUnknown                         = errors.New("cas: user attributes not included in service response")
)

type xmlServiceResponse struct {
	XMLName xml.Name `xml:"http://www.yale.edu/tp/cas serviceResponse"`

	Failure *xmlAuthenticationFailure
	Success *xmlAuthenticationSuccess
}

type xmlAuthenticationFailure struct {
	XMLName xml.Name `xml:"authenticationFailure"`
	Code    string   `xml:"code,attr"`
	Message string   `xml:",innerxml"`
}

type xmlAuthenticationSuccess struct {
	XMLName             xml.Name           `xml:"authenticationSuccess"`
	User                string             `xml:"user"`
	ProxyGrantingTicket string             `xml:"proxyGrantingTicket,omitempty"`
	Proxies             *xmlProxies        `xml:"proxies"`
	Attributes          *xmlAttributes     `xml:"attributes"`
	ExtraAttributes     []*xmlAnyAttribute `xml:",any"`
}

type xmlProxies struct {
	XMLName xml.Name `xml:"proxies"`
	Proxies []string `xml:"proxy"`
}

func (p *xmlProxies) AddProxy(proxy string) {
	p.Proxies = append(p.Proxies, proxy)
}

type xmlAttributes struct {
	XMLName                                xml.Name  `xml:"attributes"`
	AuthenticationDate                     time.Time `xml:"authenticationDate"`
	LongTermAuthenticationRequestTokenUsed bool      `xml:"longTermAuthenticationRequestTokenUsed"`
	IsFromNewLogin                         bool      `xml:"isFromNewLogin"`
	MemberOf                               []string  `xml:"memberOf"`
	UserAttributes                         *xmlUserAttributes
}

type xmlUserAttributes struct {
	XMLName       xml.Name             `xml:"userAttributes"`
	Attributes    []*xmlNamedAttribute `xml:"attribute"`
	AnyAttributes []*xmlAnyAttribute   `xml:",any"`
}

type xmlNamedAttribute struct {
	XMLName xml.Name `xml:"attribute"`
	Name    string   `xml:"name,attr,omitempty"`
	Value   string   `xml:",innerxml"`
}

type xmlAnyAttribute struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func (xsr *xmlServiceResponse) marshalXML(indent int) ([]byte, error) {
	if indent == 0 {
		return xml.Marshal(xsr)
	}

	indentStr := ""
	for i := 0; i < indent; i++ {
		indentStr += " "
	}

	return xml.MarshalIndent(xsr, "", indentStr)
}

func failureServiceResponse(code, message string) *xmlServiceResponse {
	return &xmlServiceResponse{
		Failure: &xmlAuthenticationFailure{
			Code:    code,
			Message: message,
		},
	}
}

func successServiceResponse(username, pgt string) *xmlServiceResponse {
	return &xmlServiceResponse{
		Success: &xmlAuthenticationSuccess{
			User:                username,
			ProxyGrantingTicket: pgt,
		},
	}
}
