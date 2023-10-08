import Alert from "react-bootstrap/Alert";
import Container from "react-bootstrap/Container";

export default function Notification({message}) {
  return (
    <Container fluid className="text-center">
      { message === "thanks" && <Alert variant="primary" className="m-3">
        <span>Thank you for submitting an estimate request!</span>
        <br/>
        <span className="small">You should receive a text or email soon.</span></Alert>}
      { message === "error" && <Alert variant="danger" className="m-3">Something went wrong!</Alert>}
    </Container>

  )
}