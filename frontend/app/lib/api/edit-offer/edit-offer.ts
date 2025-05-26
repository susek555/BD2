import { RegularOfferData } from "../../definitions/offer-form";

export async function editRegularOffer(data: RegularOfferData, id: string): Promise<number> {
    const response = await fetch(`/api/edit-offer/${id}`, {
        method: "PUT",
        body: JSON.stringify(data),
        headers: {
            "Content-Type": "application/json",
        },
    });

    if (response.status === 201) {
        return response.status;
    } else {
        const errorText = await response.text();
        throw new Error(`Failed to edit offer: ${response.status} – ${errorText}`);
    }
}