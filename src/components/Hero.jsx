import Container from "react-bootstrap/Container";
import Image from "react-bootstrap/Image";

export default function Hero() {
  return (
    <Container fluid className="text-center bg-light">
      <Image  src="/public/images/logo.webp" alt="Desert Cat Cookies" fluid />
    </Container>
  )
}