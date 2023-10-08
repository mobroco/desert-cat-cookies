import Container from "react-bootstrap/Container";
import Button from 'react-bootstrap/Button';
import Col from 'react-bootstrap/Col';
import Form from 'react-bootstrap/Form';
import Row from 'react-bootstrap/Row';
import Card from 'react-bootstrap/Card';
import {useState} from "react";

export default function EstimateForm() {
  const [firstName, setFirstName] = useState();
  const [lastName, setLastName] = useState();
  const [email, setEmail] = useState();
  const [phoneNumber, setPhoneNumber] = useState();
  const [cookieTheme, setCookieTheme] = useState();
  const [cookieQuantity, setCookieQuantity] = useState();
  const [pickupDate, setPickupDate] = useState(new Date(Date.now() + 14 * 24 * 60 * 60 * 1000).toISOString().split('T')[0]);
  const [anythingElse, setAnythingElse] = useState();

  return (
    <Container className="mb-5">
      <Card>
        <Card.Body>
          <Card.Title>Estimate Request</Card.Title>
          <Card.Subtitle className="text-black-50">Minimum two weeks notice for all orders. Local Phoenix area only!</Card.Subtitle>

          <Form onSubmit={(e) => {
            e.preventDefault()
            fetch("/estimates", {
              method: 'post',
              body: JSON.stringify({
                "first_name": firstName,
                "last_name":lastName,
                "email": email,
                "phone_number": phoneNumber,
                "cookie_theme": cookieTheme,
                "cookie_quantity": cookieQuantity,
                "pickup_date": pickupDate,
                "anything_else": anythingElse,
              }),
            })
              .finally(() => window.location.replace("/?message=thanks"));
          }}>
            <Row className="my-3">
              <Form.Group as={Col} controlId="formGridFirstName">
                <Form.Label>First Name</Form.Label>
                <Form.Control type="text"
                              onChange={(e) => setFirstName(e.target.value)}/>
              </Form.Group>

              <Form.Group as={Col} controlId="formGridLastName">
                <Form.Label>Last Name</Form.Label>
                <Form.Control type="text"
                              onChange={(e) => setLastName(e.target.value)}/>
              </Form.Group>
            </Row>

            <Row className="mb-3">
              <Form.Group as={Col} controlId="formGridEmail">
                <Form.Label>Email</Form.Label>
                <Form.Control type="email"
                              onChange={(e) => setEmail(e.target.value)} />
              </Form.Group>

              <Form.Group as={Col} controlId="formGridPhoneNumber">
                <Form.Label>Phone Number</Form.Label>
                <Form.Control type="text"
                              onChange={(e) => setPhoneNumber(e.target.value)} />
              </Form.Group>
            </Row>

            <Row className="mb-3">
              <Form.Group controlId="formGridCookieTheme">
                <Form.Label>Cookie Theme</Form.Label>
                <Form.Control type="text"
                              onChange={(e) => setCookieTheme(e.target.value)} />
              </Form.Group>
            </Row>

            <Row className="mb-3">
              <Form.Group as={Col} controlId="formGridCookieQuantity">
                <Form.Label>Cookie Quantity</Form.Label>
                <Form.Select defaultValue="Choose..."
                             onChange={(e) => setCookieQuantity(e.target.value)} >
                  <option>12</option>
                  <option>24</option>
                  <option>36</option>
                  <option>48</option>
                  <option>More than 48...</option>
                </Form.Select>
              </Form.Group>

              <Form.Group as={Col} controlId="formGridPickUpDate">
                <Form.Label>Desired Pick-Up Date</Form.Label>
                <Form.Control type="date"
                              placeholder={pickupDate}
                              value={pickupDate}
                              onChange={(e) => setPickupDate(e.target.value)}/>

              </Form.Group>
            </Row>

            <Row className="mb-3">
              <Form.Group id="formGridAnythingElse">
                <Form.Label>Anything else?</Form.Label>
                <Form.Control as="textarea" rows={4}
                              onChange={(e) => setAnythingElse(e.target.value)} />
              </Form.Group>
            </Row>

            <Row className="mb-3">
              <Button variant="primary" type="submit" className="m-auto w-50 w-lg-25">
                Submit Request
              </Button>
            </Row>
          </Form>
        </Card.Body>
      </Card>
    </Container>
  )
}
