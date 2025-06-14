import { OfferFormState, RegularOfferData, AuctionOfferData } from "@/app/lib/definitions/offer-form";
import { permanentRedirect } from "next/navigation";
import { editAuction, editRegularOffer } from "../lib/api/edit-offer/edit-offer";
import { uploadImages } from "../lib/api/images/upload";
import { deleteImages } from "../lib/api/images/delete";

export async function editOffer(
    id: string,
    state: OfferFormState,
    imagesToDelete: string[],
): Promise<OfferFormState> {
    let validatedFields = state.values!;

    // This property exists, the error is incorrect
    const { is_auction, images = [], ...offerData } = validatedFields;
    validatedFields = offerData;

    console.log("Is auction:", is_auction);
    console.log("Edit Offer validated fields:", validatedFields.data);
    console.log("Images to delete:", imagesToDelete);


    try {
        if (is_auction) {
            const AuctionOfferData: AuctionOfferData = validatedFields as AuctionOfferData;
            await editAuction(AuctionOfferData, parseInt(id));
        } else {
            const regularOfferData: RegularOfferData = validatedFields as RegularOfferData;
            await editRegularOffer(regularOfferData, parseInt(id));
        }
        console.log("Edited Offer ID:", id);
        if (imagesToDelete && imagesToDelete.length > 0 ) {
            await deleteImages(imagesToDelete);
        }
        if (images && images.length > 0 && images[0].name !== "") {
            await uploadImages(images, parseInt(id));
        }

    } catch (error) {
        if (typeof window !== 'undefined') {
            alert(`Error editing offer: ${error instanceof Error ? error.message : String(error)}`);
        }
        return {
            errors: {
                upload_offer: [`Failed to edit offer: ${error instanceof Error ? error.message : String(error)}`]
            },
            values: { ...validatedFields, is_auction } as OfferFormState['values']
        }
    }

    permanentRedirect("/account/listings");
}