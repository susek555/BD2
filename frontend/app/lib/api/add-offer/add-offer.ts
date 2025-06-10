import { RegularOfferData, AuctionOfferData } from "../../definitions/offer-form";

export async function postRegularOffer(data: RegularOfferData): Promise<number> {
  console.log("Posting offer data:", data);

  const response = await fetch("/api/add-offer", {
  method: "POST",
  body: JSON.stringify(data),
  headers: {
    "Content-Type": "application/json",
  },
});

  if (response.status === 201) {
    const responseData = await response.json();
    return responseData.id;
  } else {
    const errorText = await response.text();
    throw new Error(`Failed to post offer: ${response.status} – ${errorText}`);
  }
}

export async function postAuction(data: AuctionOfferData): Promise<number> {
  console.log("Posting auction data:", data);

  const response = await fetch("/api/add-auction", {
    method: "POST",
    body: JSON.stringify(data),
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (response.status === 201) {
    const responseData = await response.json();
    return responseData.id;
  } else {
    const errorText = await response.text();
    throw new Error(`Failed to post auction: ${response.status} – ${errorText}`);
  }
}

export async function publishOffer(id: number): Promise<void> {
  console.log("Publishing offer with ID:", id);

  const response = await fetch(`/api/publish-offer/${id}`, {
    method: "PUT",
  });

  if (response.status !== 200) {
    const errorText = await response.text();
    throw new Error(`Failed to publish offer: ${response.status} – ${errorText}`);
  }
}