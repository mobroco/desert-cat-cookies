import {QueryClient, QueryClientProvider} from "react-query";
import Hello from "./Hello";
const queryClient = new QueryClient()

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Hello />
    </QueryClientProvider>
  )
}