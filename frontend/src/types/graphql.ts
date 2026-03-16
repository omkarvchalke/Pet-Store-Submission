export type Pet = {
  id: string;
  storeId: string;
  name: string;
  species: "CAT" | "DOG" | "BIRD";
  age: number;
  pictureUrl: string;
  description: string;
  createdAt: string;
  status: "AVAILABLE" | "SOLD";
  soldAt?: string | null;
};

export type PurchasePetResult = {
  success: boolean;
  message: string;
};

export type CheckoutCartResult = {
  success: boolean;
  message: string;
  unavailablePetNames: string[];
};
