const API_URL = process.env.API_URL;

export async function getOrderKeys() : Promise<string[]> {
  const response = await fetch(`${API_URL}/sale-offer/order-keys`, {
    method: "GET",
  });

  if (!response.ok) {
    throw new Error("Failed to fetch order keys");
  }

  const data = await response.json();

  return data.order_keys;
}