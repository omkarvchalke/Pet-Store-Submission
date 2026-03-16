import { useState } from "react";
import type { Pet } from "../types/graphql";

type Props = {
  pet: Pet;
  onBuyNow: (petId: string) => void;
  onAddToCart: (pet: Pet) => void;
  loading?: boolean;
};

export default function PetCard({ pet, onBuyNow, onAddToCart, loading }: Props) {
  const [imageFailed, setImageFailed] = useState(false);
  const imageSrc = pet.pictureUrl ? `http://localhost:8080${pet.pictureUrl}` : "";

  return (
    <div className="pet-card">
      <div className="pet-image-wrap">
        {imageSrc && !imageFailed ? (
          <img
            src={imageSrc}
            alt={pet.name}
            className="pet-image"
            onError={() => setImageFailed(true)}
          />
        ) : (
          <div className="pet-image-fallback">Pet photo preview</div>
        )}
      </div>

      <div className="pet-header">
        <h3 className="pet-name">{pet.name}</h3>
        <span className={`species-badge species-${pet.species}`}>{pet.species}</span>
      </div>

      <div className="pet-meta">{pet.age} human years old</div>
      <p className="pet-description">{pet.description}</p>

      <div className="button-row">
        <button className="btn btn-primary" onClick={() => onBuyNow(pet.id)} disabled={loading}>
          {loading ? "Processing..." : "Purchase"}
        </button>

        <button className="btn btn-secondary" onClick={() => onAddToCart(pet)} disabled={loading}>
          {loading ? "Please wait..." : "Add to Cart"}
        </button>
      </div>
    </div>
  );
}
