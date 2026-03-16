import type { Pet } from "../types/graphql";

type Props = {
  items: Pet[];
  onRemove: (petId: string) => void;
  onCheckout: () => void;
  loading?: boolean;
};

export default function CartPanel({ items, onRemove, onCheckout, loading }: Props) {
  return (
    <aside className="cart-panel">
      <div className="cart-title-row">
        <h2 className="cart-title">Cart</h2>
        <span className="cart-count">{items.length} selected</span>
      </div>

      {items.length === 0 ? (
        <p className="cart-empty">Your cart is empty. Add a pet to continue.</p>
      ) : (
        <>
          <ul className="cart-list">
            {items.map((pet) => (
              <li key={pet.id} className="cart-item">
                <div>
                  <div className="cart-item-name">{pet.name}</div>
                  <div className="cart-item-meta">
                    {pet.species} • {pet.age} years
                  </div>
                </div>

                <button className="btn btn-secondary" onClick={() => onRemove(pet.id)}>
                  Remove
                </button>
              </li>
            ))}
          </ul>

          <button className="btn-checkout" onClick={onCheckout} disabled={loading}>
            {loading ? "Processing..." : "Checkout"}
          </button>
        </>
      )}
    </aside>
  );
}