import Hero from "./Hero";
import EstimateForm from "./EstimateForm";
import Notification from "./Notification";

export default function App() {
  const queryParameters = new URLSearchParams(window.location.search)
  const message = queryParameters.get("message")
  return (
    <>
      {message && <Notification message={message} />}
      <Hero />
      <EstimateForm />
    </>
  )
}