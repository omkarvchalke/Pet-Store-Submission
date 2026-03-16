export const GET_AVAILABLE_PETS = `
  query GetAvailablePets($storeSlug: String!) {
    customerAvailablePets(storeSlug: $storeSlug) {
      id
      storeId
      name
      species
      age
      pictureUrl
      description
      createdAt
      status
      soldAt
    }
  }
`;

export const PURCHASE_PET = `
  mutation PurchasePet($petId: String!) {
    customerPurchasePet(petId: $petId) {
      success
      message
    }
  }
`;

export const CHECKOUT_CART = `
  mutation CheckoutCart($petIds: [String]) {
    customerCheckoutCart(petIds: $petIds) {
      success
      message
      unavailablePetNames
    }
  }
`;