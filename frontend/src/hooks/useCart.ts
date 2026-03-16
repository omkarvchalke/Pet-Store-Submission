import { useMemo, useState } from "react";
import type { Pet } from "../types/graphql";

export function useCart() {
  const [items, setItems] = useState<Pet[]>([]);

  const addToCart = (pet: Pet) => {
    setItems((prev) => (prev.some((p) => p.id === pet.id) ? prev : [...prev, pet]));
  };

  const removeFromCart = (petId: string) => {
    setItems((prev) => prev.filter((p) => p.id !== petId));
  };

  const clearCart = () => setItems([]);

  const totalItems = useMemo(() => items.length, [items]);

  return { items, addToCart, removeFromCart, clearCart, totalItems };
}