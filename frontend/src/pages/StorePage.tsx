import { useEffect, useState } from "react";
import { graphqlRequest } from "../api/client";
import { CHECKOUT_CART, GET_AVAILABLE_PETS, PURCHASE_PET } from "../api/queries";
import CartPanel from "../components/CartPanel";
import PetCard from "../components/PetCard";
import { useCart } from "../hooks/useCart";
import type { CheckoutCartResult, Pet, PurchasePetResult } from "../types/graphql";

type Props = {
  storeSlug: string;
};

export default function StorePage({ storeSlug }: Props) {
  const [pets, setPets] = useState<Pet[]>([]);
  const [loading, setLoading] = useState(false);
  const [initialLoading, setInitialLoading] = useState(true);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  const { items, addToCart, removeFromCart, clearCart } = useCart();

  const loadPets = async () => {
    try {
      setError("");
      const data = await graphqlRequest<{ customerAvailablePets: Pet[] }>(
        GET_AVAILABLE_PETS,
        { storeSlug }
      );
      setPets(data.customerAvailablePets);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load pets.");
    } finally {
      setInitialLoading(false);
    }
  };

  useEffect(() => {
    loadPets();
  }, [storeSlug]);

  const handleBuyNow = async (petId: string) => {
    try {
      setLoading(true);
      setMessage("");
      setError("");

      const data = await graphqlRequest<{ customerPurchasePet: PurchasePetResult }>(
        PURCHASE_PET,
        { petId }
      );

      if (data.customerPurchasePet.success) {
        setMessage(data.customerPurchasePet.message);
        removeFromCart(petId);
      } else {
        setError(data.customerPurchasePet.message);
      }

      await loadPets();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Purchase failed.");
    } finally {
      setLoading(false);
    }
  };

  const handleAddToCart = (pet: Pet) => {
    addToCart(pet);
    setMessage(`${pet.name} was added to your cart.`);
    setError("");
  };

  const handleCheckout = async () => {
    try {
      setLoading(true);
      setMessage("");
      setError("");

      const petIds = items.map((item) => item.id);

      const data = await graphqlRequest<{ customerCheckoutCart: CheckoutCartResult }>(
        CHECKOUT_CART,
        { petIds }
      );

      if (data.customerCheckoutCart.success) {
        setMessage(data.customerCheckoutCart.message);
        clearCart();
      } else {
        setError(data.customerCheckoutCart.message);
      }

      await loadPets();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Checkout failed.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="store-shell">
      <div className="store-container">
        <header className="store-hero">
          <span className="store-eyebrow">Customer storefront</span>
          <h1 className="store-title">Happy Paws Pet Store</h1>
          <p className="store-subtitle">
            Browse available pets, bring home your favorite companion instantly,
            or add multiple pets to your cart for checkout.
          </p>

          <div className="store-meta">
            <span className="store-pill">Store: {storeSlug}</span>
            <span className="store-pill">Available pets: {pets.length}</span>
            <span className="store-pill">Cart items: {items.length}</span>
          </div>

          {message ? <div className="alert alert-success">{message}</div> : null}
          {error ? <div className="alert alert-error">{error}</div> : null}
        </header>

        <div className="store-grid">
          <section>
            <h2 className="section-title">Available Pets</h2>

            {initialLoading ? (
              <div className="empty-state">
                <span className="loading-text">Loading pets...</span>
              </div>
            ) : pets.length === 0 ? (
              <div className="empty-state">No pets available right now.</div>
            ) : (
              <div className="pets-grid">
                {pets.map((pet) => (
                  <PetCard
                    key={pet.id}
                    pet={pet}
                    onBuyNow={handleBuyNow}
                    onAddToCart={handleAddToCart}
                    loading={loading}
                  />
                ))}
              </div>
            )}
          </section>

          <CartPanel
            items={items}
            onRemove={removeFromCart}
            onCheckout={handleCheckout}
            loading={loading}
          />
        </div>
      </div>
    </div>
  );
}
