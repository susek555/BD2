import { OfferFormState, RegularOfferData } from "@/app/lib/definitions/offer-form";
import { postRegularOffer } from "@/app/lib/api/add-offer/add-offer";
import { UploadImages } from "../lib/api/images/upload";

export async function addOffer(
    state: OfferFormState,
): Promise<OfferFormState> {
    let validatedFields = state.values!;

    // This property exists, the error is incorrect
    const { is_auction, images, ...offerData } = validatedFields;
    validatedFields = offerData;

    console.log("Is auction:", is_auction);
    console.log("Add Offer validated fields:", validatedFields.data);


    try {
        if (is_auction) {
            console.log("Adding auction offer");
        } else {
            const regularOfferData: RegularOfferData = validatedFields as RegularOfferData;
            const id = await postRegularOffer(regularOfferData);
            console.log("Posted Offer ID:", id);
            await UploadImages(images!, id);
        }

        return state;

    } catch (error) {
        if (typeof window !== 'undefined') {
            alert(`Error adding offer: ${error instanceof Error ? error.message : String(error)}`);
        }
        return {
            errors: {
                upload_offer: [`Failed to upload offer: ${error instanceof Error ? error.message : String(error)}`]
            },
            values: { ...validatedFields, is_auction } as OfferFormState['values']
        }
    }
}