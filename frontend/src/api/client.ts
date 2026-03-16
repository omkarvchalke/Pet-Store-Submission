const GRAPHQL_URL = import.meta.env.VITE_GRAPHQL_URL;
const username = import.meta.env.VITE_CUSTOMER_AUTH_USERNAME;
const password = import.meta.env.VITE_CUSTOMER_AUTH_PASSWORD;

const authHeader = `Basic ${btoa(`${username}:${password}`)}`;

export async function graphqlRequest<T>(
  query: string,
  variables?: Record<string, unknown>
): Promise<T> {
  const response = await fetch(GRAPHQL_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: authHeader,
    },
    body: JSON.stringify({ query, variables }),
  });

  const result = await response.json();

  if (result.errors?.length) {
    throw new Error(result.errors[0].message || "GraphQL request failed");
  }

  return result.data;
}