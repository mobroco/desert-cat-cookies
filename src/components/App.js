import { default as DesertCatCookies } from "./desertcatcookies/Index";
import { default as GreasyShadows } from "./greasyshadows/Index";


import { QueryClient, QueryClientProvider } from 'react-query';
import Login from "./Login";
import Support from "./Support";

const queryClient = new QueryClient()

export default function App() {
  const urlSearchParams = new URLSearchParams(window.location.search);
  const params = Object.fromEntries(urlSearchParams.entries());
  const parts = window.location.pathname.split("/");

  console.log(parts)
  let site = <h1>???</h1>
  if (parts[1] === 'login') {
    site = <Login params={params} parts={parts}/>
  } else if (parts[1] === 'support') {
    site = <Support params={params} parts={parts}/>
  } else if (window.seed === 'desert-cat-cookies') {
    site = <DesertCatCookies params={params} parts={parts} />
  } else if (window.seed === 'greasy-shadows') {
    site = <GreasyShadows params={params} parts={parts}/>
  }
  return (
    <QueryClientProvider client={queryClient}>
      { site }
    </QueryClientProvider>
  )
}