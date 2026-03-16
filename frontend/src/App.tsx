import StorePage from "./pages/StorePage";

function App() {
  const path = window.location.pathname;
  const storeSlug = path.startsWith("/store/") ? path.replace("/store/", "") : "happy-paws";

  return <StorePage storeSlug={storeSlug} />;
}

export default App;
