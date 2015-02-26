package cas

import (
	"encoding/xml"
	"testing"
	"time"
)

func TestUnmarshalFailureServiceResponse(t *testing.T) {
	s := `<cas:serviceResponse xmlns:cas="http://www.yale.edu/tp/cas">
	<cas:authenticationFailure code="INVALID_TICKET">
			Ticket ST-1856339-aA5Yuvrxzpv8Tau1cYQ7 not recognized
	</cas:authenticationFailure>
</cas:serviceResponse>`

	_, err := ParseServiceResponse([]byte(s))
	if err == nil {
		t.Fatalf("Expected ParseServiceResponse to return error, got <nil>")
	}

	if err.Error() != "INVALID_TICKET: Ticket ST-1856339-aA5Yuvrxzpv8Tau1cYQ7 not recognized" {
		t.Errorf("Expected err to be <INVALID_TICKET: Ticket ST-1856339-aA5Yuvrxzpv8Tau1cYQ7 not recognized>, got <%v>",
			err.Error())
	}

	if authErr, ok := err.(AuthenticationError); ok {
		if authErr.Code != INVALID_TICKET {
			t.Errorf("Expected Code to be <INVALID_TICKET>, got <%v>", authErr.Code)
		}

		expected := "Ticket ST-1856339-aA5Yuvrxzpv8Tau1cYQ7 not recognized"
		if authErr.Message != expected {
			t.Errorf("Expected Message to be <%v>, got <%v>", expected, authErr.Message)
		}
	}
}

func TestUnmarshalSuccessfulServiceResponse(t *testing.T) {
	s := `<cas:serviceResponse xmlns:cas="http://www.yale.edu/tp/cas">
   <cas:authenticationSuccess>
     <cas:user>username</cas:user>
     <cas:proxyGrantingTicket>PGTIOU-84678-8a9d3r389439</cas:proxyGrantingTicket>
   </cas:authenticationSuccess>
</cas:serviceResponse>`

	sr, err := ParseServiceResponse([]byte(s))
	if err != nil {
		t.Errorf("Unmarshal service response failed: %v", err)
	}

	if sr.User != "username" {
		t.Errorf("Expected User to be <username>, got <%s>", sr.User)
	}

	if sr.ProxyGrantingTicket != "PGTIOU-84678-8a9d3r389439" {
		t.Errorf("Expected ProxyGrantingTicket to be <PGTIOU-84678-8a9d3r389439>, got <%s>",
			sr.ProxyGrantingTicket)
	}
}

func TestUnmarshalSuccessfulServiceResponseWithAttributes(t *testing.T) {
	s := `<?xml version="1.0"?>
<cas:serviceResponse xmlns:cas="http://www.yale.edu/tp/cas">
  <cas:authenticationSuccess>
    <cas:user>username</cas:user>
    <cas:attributes>
      <cas:authenticationDate>2015-02-10T14:28:42Z</cas:authenticationDate>
      <cas:longTermAuthenticationRequestTokenUsed>false</cas:longTermAuthenticationRequestTokenUsed>
      <cas:isFromNewLogin>true</cas:isFromNewLogin>
    </cas:attributes>
    <cas:proxyGrantingTicket>PGTIOU-84678-8a9d...</cas:proxyGrantingTicket>
  </cas:authenticationSuccess>
</cas:serviceResponse>`

	sr, err := ParseServiceResponse([]byte(s))
	if err != nil {
		t.Errorf("Unmarshal service response failed: %v", err)
	}

	authDate := time.Date(2015, 2, 10, 14, 28, 42, 0, time.UTC)
	if sr.AuthenticationDate != authDate {
		t.Errorf("Expected AuthenticationDate to be <%v>, got <%v>", authDate,
			sr.AuthenticationDate)
	}

	if sr.IsRememberedLogin != false {
		t.Errorf("Expected IsRememberedLogin to be false, got <%v>",
			sr.IsRememberedLogin)
	}

	if sr.IsNewLogin != true {
		t.Errorf("Expected IsNewLogin to be true, got <%v>",
			sr.IsNewLogin)
	}
}

func TestUnmarshalSuccessfulServiceResponseWithUserAttributesRecommendedForm(t *testing.T) {
	s := `<?xml version="1.0"?>
<cas:serviceResponse xmlns:cas="http://www.yale.edu/tp/cas">
  <cas:authenticationSuccess>
    <cas:user>username</cas:user>
    <cas:attributes>
      <cas:authenticationDate>2015-02-10T14:28:42Z</cas:authenticationDate>
      <cas:longTermAuthenticationRequestTokenUsed>false</cas:longTermAuthenticationRequestTokenUsed>
      <cas:isFromNewLogin>true</cas:isFromNewLogin>
      <cas:memberOf>Group1</cas:memberOf>
      <cas:memberOf>Group2</cas:memberOf>
      <cas:memberOf>Group3</cas:memberOf>
      <cas:userAttributes>
        <cas:attribute name="firstname">John</cas:attribute>
        <cas:attribute name="lastname">Doe</cas:attribute>
        <cas:attribute name="title">Mr.</cas:attribute>
        <cas:attribute name="email">jdoe@example.org</cas:attribute>
        <cas:attribute name="affiliation">staff</cas:attribute>
        <cas:attribute name="affiliation">faculty</cas:attribute>
      </cas:userAttributes>
    </cas:attributes>
    <cas:proxyGrantingTicket>PGTIOU-84678-8a9d...</cas:proxyGrantingTicket>
  </cas:authenticationSuccess>
</cas:serviceResponse>`

	sr, err := ParseServiceResponse([]byte(s))
	if err != nil {
		t.Errorf("Unmarshal service response failed: %v", err)
	}

	memberOf := []string{"Group1", "Group2", "Group3"}
	for i, member := range sr.MemberOf {
		if memberOf[i] != member {
			t.Errorf("Expected MemberOf[%d] to be <%v>, got <%v>", i, memberOf[i],
				member)
		}
	}

	if len(sr.Attributes) != 5 {
		t.Errorf("Expected Attributes to have 5 items, got %v: %v",
			len(sr.Attributes), sr.Attributes)
	}

	if v := sr.Attributes.Get("firstname"); v != "John" {
		t.Errorf("Expected firstname attribute to be <John>, got <%v>", v)
	}

	if v := sr.Attributes.Get("lastname"); v != "Doe" {
		t.Errorf("Expected lastname attribute to be <Doe>, got <%v>", v)
	}

	if v := sr.Attributes.Get("title"); v != "Mr." {
		t.Errorf("Expected title attribute to be <Mr.>, got <%v>", v)
	}

	if v := sr.Attributes.Get("email"); v != "jdoe@example.org" {
		t.Errorf("Expected email attribute to be <jdoe@example.org>, got <%v>", v)
	}

	if v := sr.Attributes.Get("affiliation"); v != "staff" {
		t.Errorf("Expected affiliation attribute to be <staff>, got <%v>", v)
	}

	expectedAffiliations := []string{"staff", "faculty"}
	if affiliations, ok := sr.Attributes["affiliation"]; ok {
		for i, affiliation := range affiliations {
			if expectedAffiliations[i] != affiliation {
				t.Errorf("Expected affiliation attribute to be <%v>, got <%v>",
					expectedAffiliations[i], affiliation)
			}
		}
	} else {
		t.Errorf("Expected affiliation attribute to exist, but its !ok")
	}
}

func TestUnmarshalSuccessfulServiceResponseWithUserAttributesXmlAnyForm(t *testing.T) {
	s := `<?xml version="1.0"?>
<cas:serviceResponse xmlns:cas="http://www.yale.edu/tp/cas">
  <cas:authenticationSuccess>
    <cas:user>username</cas:user>
    <cas:attributes>
      <cas:authenticationDate>2015-02-10T14:28:42Z</cas:authenticationDate>
      <cas:longTermAuthenticationRequestTokenUsed>false</cas:longTermAuthenticationRequestTokenUsed>
      <cas:isFromNewLogin>true</cas:isFromNewLogin>
      <cas:memberOf>Group1</cas:memberOf>
      <cas:memberOf>Group2</cas:memberOf>
      <cas:memberOf>Group3</cas:memberOf>
      <cas:userAttributes>
        <cas:firstname>John</cas:firstname>
        <cas:lastname>Doe</cas:lastname>
        <cas:title>Mr.</cas:title>
        <cas:email>jdoe@example.org</cas:email>
        <cas:affiliation>staff</cas:affiliation>
        <cas:affiliation>faculty</cas:affiliation>
      </cas:userAttributes>
    </cas:attributes>
    <cas:proxyGrantingTicket>PGTIOU-84678-8a9d...</cas:proxyGrantingTicket>
  </cas:authenticationSuccess>
</cas:serviceResponse>`

	sr, err := ParseServiceResponse([]byte(s))
	if err != nil {
		t.Errorf("Unmarshal service response failed: %v", err)
	}

	memberOf := []string{"Group1", "Group2", "Group3"}
	for i, member := range sr.MemberOf {
		if memberOf[i] != member {
			t.Errorf("Expected MemberOf[%d] to be <%v>, got <%v>", i, memberOf[i],
				member)
		}
	}

	if len(sr.Attributes) != 5 {
		t.Errorf("Expected Attributes to have 5 items, got %v: %v",
			len(sr.Attributes), sr.Attributes)
	}

	if v := sr.Attributes.Get("firstname"); v != "John" {
		t.Errorf("Expected firstname attribute to be <John>, got <%v>", v)
	}

	if v := sr.Attributes.Get("lastname"); v != "Doe" {
		t.Errorf("Expected lastname attribute to be <Doe>, got <%v>", v)
	}

	if v := sr.Attributes.Get("title"); v != "Mr." {
		t.Errorf("Expected title attribute to be <Mr.>, got <%v>", v)
	}

	if v := sr.Attributes.Get("email"); v != "jdoe@example.org" {
		t.Errorf("Expected email attribute to be <jdoe@example.org>, got <%v>", v)
	}

	if v := sr.Attributes.Get("affiliation"); v != "staff" {
		t.Errorf("Expected affiliation attribute to be <staff>, got <%v>", v)
	}

	expectedAffiliations := []string{"staff", "faculty"}
	if affiliations, ok := sr.Attributes["affiliation"]; ok {
		for i, affiliation := range affiliations {
			if expectedAffiliations[i] != affiliation {
				t.Errorf("Expected affiliation attribute to be <%v>, got <%v>",
					expectedAffiliations[i], affiliation)
			}
		}
	} else {
		t.Errorf("Expected affiliation attribute to exist, but its !ok")
	}
}

func TestUnmarshalSuccessfulServiceResponseWithRubycasExtraAttributes(t *testing.T) {
	s := `<?xml version="1.0"?>
<cas:serviceResponse xmlns:cas="http://www.yale.edu/tp/cas">
	<cas:authenticationSuccess>
		<cas:user>username</cas:user>
		<cas:attributes>
			<cas:authenticationDate>2015-02-10T14:28:42Z</cas:authenticationDate>
			<cas:longTermAuthenticationRequestTokenUsed>false</cas:longTermAuthenticationRequestTokenUsed>
			<cas:isFromNewLogin>true</cas:isFromNewLogin>
		</cas:attributes>
		<cas:proxyGrantingTicket>PGTIOU-84678-8a9d...</cas:proxyGrantingTicket>
		<firstname>John</firstname>
		<lastname>Doe</lastname>
		<title>Mr.</title>
		<email>jdoe@example.org</email>
		<affiliation><![CDATA[---
- staff
- faculty
]]></affiliation>
	</cas:authenticationSuccess>
</cas:serviceResponse>`

	sr, err := ParseServiceResponse([]byte(s))
	if err != nil {
		t.Errorf("Unmarshal service response failed: %v", err)
	}

	memberOf := []string{"Group1", "Group2", "Group3"}
	for i, member := range sr.MemberOf {
		if memberOf[i] != member {
			t.Errorf("Expected MemberOf[%d] to be <%v>, got <%v>", i, memberOf[i],
				member)
		}
	}

	if len(sr.Attributes) != 5 {
		t.Errorf("Expected Attributes to have 5 items, got %v: %v",
			len(sr.Attributes), sr.Attributes)
	}

	if v := sr.Attributes.Get("firstname"); v != "John" {
		t.Errorf("Expected firstname attribute to be <John>, got <%v>", v)
	}

	if v := sr.Attributes.Get("lastname"); v != "Doe" {
		t.Errorf("Expected lastname attribute to be <Doe>, got <%v>", v)
	}

	if v := sr.Attributes.Get("title"); v != "Mr." {
		t.Errorf("Expected title attribute to be <Mr.>, got <%v>", v)
	}

	if v := sr.Attributes.Get("email"); v != "jdoe@example.org" {
		t.Errorf("Expected email attribute to be <jdoe@example.org>, got <%v>", v)
	}

	if v := sr.Attributes.Get("affiliation"); v != "staff" {
		t.Errorf("Expected affiliation attribute to be <staff>, got <%v>", v)
	}

	expectedAffiliations := []string{"staff", "faculty"}
	if affiliations, ok := sr.Attributes["affiliation"]; ok {
		for i, affiliation := range affiliations {
			if expectedAffiliations[i] != affiliation {
				t.Errorf("Expected affiliation attribute to be <%v>, got <%v>",
					expectedAffiliations[i], affiliation)
			}
		}
	} else {
		t.Errorf("Expected affiliation attribute to exist, but its !ok")
	}
}

func TestFailureServiceResponse(t *testing.T) {
	sr := failureServiceResponse("INVALID_TICKET", "Ticket ST-1856339-aA5Yuvrxzpv8Tau1cYQ7 not recognized")
	s, err := sr.marshalXML(2)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<serviceResponse xmlns="http://www.yale.edu/tp/cas">
  <authenticationFailure code="INVALID_TICKET">Ticket ST-1856339-aA5Yuvrxzpv8Tau1cYQ7 not recognized</authenticationFailure>
</serviceResponse>`

	if string(s) != expected {
		t.Errorf("Expected marshalled results to match. Expected:\n%s\nGot:\n%s", expected, s)
	}
}

func TestSuccessfulServiceResponse(t *testing.T) {
	sr := successServiceResponse("username", "")
	s, err := sr.marshalXML(2)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<serviceResponse xmlns="http://www.yale.edu/tp/cas">
  <authenticationSuccess>
    <user>username</user>
  </authenticationSuccess>
</serviceResponse>`

	if string(s) != expected {
		t.Errorf("Expected marshalled results to match. Expected:\n%s\nGot:\n%s", expected, s)
	}
}

func TestSuccessfulServiceResponseWithProxyGrantingTicket(t *testing.T) {
	sr := successServiceResponse("username", "PGTIOU-84678-8a9d3r389439")
	s, err := sr.marshalXML(2)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<serviceResponse xmlns="http://www.yale.edu/tp/cas">
  <authenticationSuccess>
    <user>username</user>
    <proxyGrantingTicket>PGTIOU-84678-8a9d3r389439</proxyGrantingTicket>
  </authenticationSuccess>
</serviceResponse>`

	if string(s) != expected {
		t.Errorf("Expected marshalled results to match. Expected:\n%s\nGot:\n%s", expected, s)
	}
}

func TestSuccessfulServiceResponseWithProxies(t *testing.T) {
	sr := successServiceResponse("username", "")
	sr.Success.Proxies = &xmlProxies{Proxies: []string{"test.host"}}

	s, err := sr.marshalXML(2)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<serviceResponse xmlns="http://www.yale.edu/tp/cas">
  <authenticationSuccess>
    <user>username</user>
    <proxies>
      <proxy>test.host</proxy>
    </proxies>
  </authenticationSuccess>
</serviceResponse>`

	if string(s) != expected {
		t.Errorf("Expected marshalled results to match. Expected:\n%s\nGot:\n%s", expected, s)
	}
}

func TestSuccessfulServiceResponseWithBasicAttributes(t *testing.T) {
	sr := successServiceResponse("username", "")
	sr.Success.Attributes = &xmlAttributes{
		AuthenticationDate:                     time.Date(2015, 02, 10, 14, 28, 42, 0, time.UTC),
		LongTermAuthenticationRequestTokenUsed: false,
		IsFromNewLogin:                         true,
	}

	s, err := sr.marshalXML(2)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<serviceResponse xmlns="http://www.yale.edu/tp/cas">
  <authenticationSuccess>
    <user>username</user>
    <attributes>
      <authenticationDate>2015-02-10T14:28:42Z</authenticationDate>
      <longTermAuthenticationRequestTokenUsed>false</longTermAuthenticationRequestTokenUsed>
      <isFromNewLogin>true</isFromNewLogin>
    </attributes>
  </authenticationSuccess>
</serviceResponse>`

	if string(s) != expected {
		t.Errorf("Expected marshalled results to match. Expected:\n%s\nGot:\n%s", expected, s)
	}
}

func TestSuccessfulServiceResponseWithMemberOfAttributes(t *testing.T) {
	sr := successServiceResponse("username", "")
	sr.Success.Attributes = &xmlAttributes{
		AuthenticationDate:                     time.Date(2015, 02, 10, 14, 28, 42, 0, time.UTC),
		LongTermAuthenticationRequestTokenUsed: false,
		IsFromNewLogin:                         true,
		MemberOf:                               []string{"staff", "faculty", "testing"},
	}

	s, err := sr.marshalXML(2)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<serviceResponse xmlns="http://www.yale.edu/tp/cas">
  <authenticationSuccess>
    <user>username</user>
    <attributes>
      <authenticationDate>2015-02-10T14:28:42Z</authenticationDate>
      <longTermAuthenticationRequestTokenUsed>false</longTermAuthenticationRequestTokenUsed>
      <isFromNewLogin>true</isFromNewLogin>
      <memberOf>staff</memberOf>
      <memberOf>faculty</memberOf>
      <memberOf>testing</memberOf>
    </attributes>
  </authenticationSuccess>
</serviceResponse>`

	if string(s) != expected {
		t.Errorf("Expected marshalled results to match. Expected:\n%s\nGot:\n%s", expected, s)
	}
}

func TestSuccessfulServiceResponseWithNamedUserAttributes(t *testing.T) {
	sr := successServiceResponse("username", "")
	sr.Success.Attributes = &xmlAttributes{
		AuthenticationDate:                     time.Date(2015, 02, 10, 14, 28, 42, 0, time.UTC),
		LongTermAuthenticationRequestTokenUsed: false,
		IsFromNewLogin:                         true,
		MemberOf:                               []string{"staff", "faculty", "testing"},
		UserAttributes: &xmlUserAttributes{
			Attributes: []*xmlNamedAttribute{
				&xmlNamedAttribute{Name: "firstname", Value: "Enoch"},
				&xmlNamedAttribute{Name: "lastname", Value: "Root"},
			},
		},
	}

	s, err := sr.marshalXML(2)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<serviceResponse xmlns="http://www.yale.edu/tp/cas">
  <authenticationSuccess>
    <user>username</user>
    <attributes>
      <authenticationDate>2015-02-10T14:28:42Z</authenticationDate>
      <longTermAuthenticationRequestTokenUsed>false</longTermAuthenticationRequestTokenUsed>
      <isFromNewLogin>true</isFromNewLogin>
      <memberOf>staff</memberOf>
      <memberOf>faculty</memberOf>
      <memberOf>testing</memberOf>
      <userAttributes>
        <attribute name="firstname">Enoch</attribute>
        <attribute name="lastname">Root</attribute>
      </userAttributes>
    </attributes>
  </authenticationSuccess>
</serviceResponse>`

	if string(s) != expected {
		t.Errorf("Expected marshalled results to match. Expected:\n%s\nGot:\n%s", expected, s)
	}
}

func TestSuccessfulServiceResponseWithAnyUserAttributes(t *testing.T) {
	sr := successServiceResponse("username", "")
	sr.Success.Attributes = &xmlAttributes{
		AuthenticationDate:                     time.Date(2015, 02, 10, 14, 28, 42, 0, time.UTC),
		LongTermAuthenticationRequestTokenUsed: false,
		IsFromNewLogin:                         true,
		MemberOf:                               []string{"staff", "faculty", "testing"},
		UserAttributes: &xmlUserAttributes{
			AnyAttributes: []*xmlAnyAttribute{
				&xmlAnyAttribute{
					XMLName: xml.Name{Local: "firstname"},
					Value:   "Enoch"},
				&xmlAnyAttribute{
					XMLName: xml.Name{Local: "lastname"},
					Value:   "Root"},
			},
		},
	}

	s, err := sr.marshalXML(2)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<serviceResponse xmlns="http://www.yale.edu/tp/cas">
  <authenticationSuccess>
    <user>username</user>
    <attributes>
      <authenticationDate>2015-02-10T14:28:42Z</authenticationDate>
      <longTermAuthenticationRequestTokenUsed>false</longTermAuthenticationRequestTokenUsed>
      <isFromNewLogin>true</isFromNewLogin>
      <memberOf>staff</memberOf>
      <memberOf>faculty</memberOf>
      <memberOf>testing</memberOf>
      <userAttributes>
        <firstname>Enoch</firstname>
        <lastname>Root</lastname>
      </userAttributes>
    </attributes>
  </authenticationSuccess>
</serviceResponse>`

	if string(s) != expected {
		t.Errorf("Expected marshalled results to match. Expected:\n%s\nGot:\n%s", expected, s)
	}
}

func TestSuccessfulServiceResponseWithRubyCasAttributes(t *testing.T) {
	sr := successServiceResponse("username", "")
	sr.Success.ExtraAttributes = []*xmlAnyAttribute{
		&xmlAnyAttribute{
			XMLName: xml.Name{Local: "firstname"},
			Value:   "Enoch",
		},
		&xmlAnyAttribute{
			XMLName: xml.Name{Local: "lastname"},
			Value:   "Root",
		},
		&xmlAnyAttribute{
			XMLName: xml.Name{Local: "groups"},
			Value: `---
- staff
- faculty
- testing`,
		},
	}

	s, err := sr.marshalXML(2)
	if err != nil {
		t.Fatal(err)
	}

	// The XML encoded YAML format should still work with RubyCAS
	// At the moment though in order for RubyCAS to parse the XML it would
	// have to have the "cas:" XML prefix on all of the CAS elements.
	expected := `<serviceResponse xmlns="http://www.yale.edu/tp/cas">
  <authenticationSuccess>
    <user>username</user>
    <firstname>Enoch</firstname>
    <lastname>Root</lastname>
    <groups>---&#xA;- staff&#xA;- faculty&#xA;- testing</groups>
  </authenticationSuccess>
</serviceResponse>`

	if string(s) != expected {
		t.Errorf("Expected marshalled results to match. Expected:\n%s\nGot:\n%s", expected, s)
	}
}
