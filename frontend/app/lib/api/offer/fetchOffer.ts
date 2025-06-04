import { fetchWithRefresh } from "../fetchWithRefresh";

const API_URL = process.env.API_URL;

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export async function getOfferData(id : number) : Promise<any> {
  const response = await fetchWithRefresh(`${API_URL}/sale-offer/id/${id}`, {
    method: "GET",
  });

  if (!response.ok) {
    throw new Error("Failed to fetch offer data");
  }

  const data = await response.json();

return data;
}

export async function getSellerName(sellerId: number): Promise<string> {
  const response = await fetch(`${API_URL}/users/id/${sellerId}`, {
    method: "GET",
  });

  if (!response.ok) {
    throw new Error("Failed to fetch seller name");
  }

  const data = await response.json();
  return data.username;
}