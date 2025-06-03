import { API_URL } from "@/app/lib/constants";

export async function getMyCurrentBid(auction_id: number, user_id: number) {

  const response = await fetch(`${API_URL}/bid/highest/auction/${auction_id}/bidder/${user_id}`, {
    method: "GET",
  });

  if (!response.ok) {
    throw new Error("Failed to fetch my current bid");
  }

  const data = await response.json();

  return data.amount;
}

// Place bid

export type PlaceBidData = {
  amount: number;
  auction_id: number;
}

export async function placeBid(bidData: PlaceBidData) {
  const response = await fetch(`/api/bid`, {
    method: "POST",
    body: JSON.stringify(bidData),
  });

  if (!response.ok) {
    throw new Error("Failed to place bid");
  }

  const data = await response.json();
  return data;
}