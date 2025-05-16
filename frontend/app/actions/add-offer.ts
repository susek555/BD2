import { AddOfferFormSchema, AddOfferFormState } from "../lib/definitions";
import { permanentRedirect } from "next/navigation";

export async function addOffer(
    state: AddOfferFormState,
    formData: FormData
): Promise<AddOfferFormState> {
    console.log("Add Offer form data:", Object.fromEntries(formData.entries()));

    const formDataObj = Object.fromEntries(formData.entries());

    const normalizedData = {
        ...formDataObj,
        productionYear: parseInt(formDataObj.productionYear as string),
        mileage: parseInt(formDataObj.mileage as string),
        numberOfDoors: parseInt(formDataObj.numberOfDoors as string),
        numberOfSeats: parseInt(formDataObj.numberOfSeats as string),
        power: parseInt(formDataObj.power as string),
        engineDisplacement: parseInt(formDataObj.engineDisplacement as string),
        dateOfFirstRegistration: formDataObj.dateOfFirstRegistration
            ? new Date(formDataObj.dateOfFirstRegistration as string)
            : undefined,
        price: parseInt(formDataObj.price as string),
        isAuction: formDataObj.isAuction === "true" ? true : false,
        auctionEndDate: formDataObj.auctionEndDate
            ? new Date(formDataObj.auctionEndDate as string)
            : undefined,
        buyNowAuctionPrice: formDataObj.buyNowAuctionPrice
            ? parseInt(formDataObj.buyNowAuctionPrice as string)
            : undefined
    }

    const validatedFields = AddOfferFormSchema.safeParse(normalizedData);
    console.log("Add Offer validation result:", validatedFields);


    if (!validatedFields.success) {
        console.log("Add Offer validation errors:", validatedFields.error.flatten().fieldErrors);
        return {
            errors: validatedFields.error.flatten().fieldErrors,
            values: normalizedData as AddOfferFormState['values']
        };
    }


    return state;

    // permanentRedirect("/");
    //TODO redirect to my offers
}