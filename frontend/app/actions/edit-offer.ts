import { OfferFormState, RegularOfferData } from "@/app/lib/definitions/offer-form";
import { permanentRedirect } from "next/navigation";
import { editRegularOffer } from "../lib/api/edit-offer/edit-offer";

export async function editOffer(
    state: OfferFormState,
    id: string
): Promise<OfferFormState> {
    let validatedFields = state.values!;

    // This property exists, the error is incorrect
    const { is_auction, ...offerData } = validatedFields;
    validatedFields = offerData;

    console.log("Is auction:", is_auction);
    console.log("Edit Offer validated fields:", validatedFields.data);


    try {
        if (is_auction) {
            console.log("Editing auction offer");
        } else {
            const regularOfferData: RegularOfferData = validatedFields as RegularOfferData;
            await editRegularOffer(regularOfferData, id);
        }

        permanentRedirect("/account/listings");
    } catch (error) {
        if (typeof window !== 'undefined') {
            alert(`Error editing offer: ${error instanceof Error ? error.message : String(error)}`);
        }
        return {
            errors: {},
            values: { ...validatedFields, is_auction } as OfferFormState['values']
        }
    }
}