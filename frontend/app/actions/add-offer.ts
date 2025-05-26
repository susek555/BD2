import { offerFormState, RegularOfferData } from "@/app/lib/definitions/offer-form";
import { permanentRedirect } from "next/navigation";
import { postRegularOffer } from "@/app/lib/api/add-offer/add-offer";

export async function addOffer(
    state: offerFormState,
): Promise<offerFormState> {
    let validatedFields = state.values!;

    // This property exists, the error is incorrect
    const { is_auction, ...offerData } = validatedFields;
    validatedFields = offerData;

    console.log("Is auction:", is_auction);
    console.log("Add Offer validated fields:", validatedFields.data);


    try {
        if (is_auction) {
            console.log("Adding auction offer");
        } else {
            const regularOfferData: RegularOfferData = validatedFields as RegularOfferData;
            await postRegularOffer(regularOfferData);
        }

        permanentRedirect("/account/listings");
    } catch (error) {
        if (typeof window !== 'undefined') {
            alert(`Error adding offer: ${error instanceof Error ? error.message : String(error)}`);
        }
        return {
            errors: {},
            values: { ...validatedFields, is_auction } as offerFormState['values']
        }
    }
}