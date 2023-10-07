import Container from "react-bootstrap/Container";

export default function Hero() {
  return (
    <Container fluid className="text-center bg-light">
      <Container className="p-5 bg-light">
        <h1 className="display-4 fw-bold">Desert Cat Cookies</h1>
        <hr/>
        <p>Tempe, Arizona</p>
      </Container>
    </Container>
  )
}