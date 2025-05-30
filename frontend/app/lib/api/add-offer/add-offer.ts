import { RegularOfferData } from "../../definitions/offer-form";

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
    return response.status;
  } else {
    const errorText = await response.text();
    throw new Error(`Failed to post offer: ${response.status} – ${errorText}`);
  }
}