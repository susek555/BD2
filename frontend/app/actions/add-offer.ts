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
        production_year: parseInt(formDataObj.production_year as string),
        mileage: parseInt(formDataObj.mileage as string),
        number_of_doors: parseInt(formDataObj.number_of_doors as string),
        number_of_seats: parseInt(formDataObj.number_of_seats as string),
        number_of_gears: parseInt(formDataObj.number_of_gears as string),
        engine_power: parseInt(formDataObj.engine_power as string),
        engine_capacity: parseInt(formDataObj.engine_capacity as string),
        price: parseInt(formDataObj.price as string),
        is_auction: formDataObj.is_auction === "true" ? true : false,
        buy_now_auction_price: formDataObj.buy_now_auction_price
            ? parseInt(formDataObj.buy_now_auction_price as string)
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