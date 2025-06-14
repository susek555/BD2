import { RegularOfferData, AuctionOfferData } from "../../definitions/offer-form";

export async function editRegularOffer(
    data: RegularOfferData,
    id: number
): Promise<number> {
  console.log("Editing offer data:", data);

  const response = await fetch(`/api/edit-offer/${id}`, {
  method: "PUT",
  body: JSON.stringify(data),
  headers: {
    "Content-Type": "application/json",
  },
});

  if (response.status === 201 || response.status === 200) {
    const responseData = await response.json();
    return responseData.id;
  } else {
    const errorText = await response.text();
    throw new Error(`Failed to post offer: ${response.status} – ${errorText}`);
  }
}

export async function editAuction(
    data: AuctionOfferData,
    id: number
): Promise<number> {
  console.log("Editing auction data:", data);

  const response = await fetch(`/api/edit-auction/${id}`, {
  method: "PUT",
  body: JSON.stringify(data),
  headers: {
    "Content-Type": "application/json",
  },
});

  if (response.status === 201 || response.status === 200) {
    const responseData = await response.json();
    return responseData.id;
  } else {
    const errorText = await response.text();
    throw new Error(`Failed to post offer: ${response.status} – ${errorText}`);
  }
}