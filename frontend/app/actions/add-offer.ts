import { OfferDetailsFormSchema, AddOfferFormState, OfferPricingFormSchema, CombinedOfferFormSchema } from "@/app/lib/definitions/offer-form";
import { permanentRedirect } from "next/navigation";

export async function addOffer(
    state: AddOfferFormState,
    formData: FormData,
    detailsPart: boolean
): Promise<AddOfferFormState> {
    console.log("Add Offer form data:", Object.fromEntries(formData.entries()));

    const formDataObj = Object.fromEntries(formData.entries());

    const normalizedData = {
        ...formDataObj,
        production_year: formDataObj.production_year ? parseInt(formDataObj.production_year as string) || undefined : undefined,
        mileage: formDataObj.mileage ? parseInt(formDataObj.mileage as string) || undefined : undefined,
        number_of_doors: formDataObj.number_of_doors ? parseInt(formDataObj.number_of_doors as string) || undefined : undefined,
        number_of_seats: formDataObj.number_of_seats ? parseInt(formDataObj.number_of_seats as string) || undefined : undefined,
        number_of_gears: formDataObj.number_of_gears ? parseInt(formDataObj.number_of_gears as string) || undefined : undefined,
        engine_power: formDataObj.engine_power ? parseInt(formDataObj.engine_power as string) || undefined : undefined,
        engine_capacity: formDataObj.engine_capacity ? parseInt(formDataObj.engine_capacity as string) || undefined : undefined,
        price: formDataObj.price ? parseInt(formDataObj.price as string) || undefined : undefined,
        margin: formDataObj.margin
            ? (parseInt((formDataObj.margin as string).replace('%', '')) || undefined)
            : undefined,
        is_auction: formDataObj.is_auction === "true" ? true : false,
        buy_now_auction_price: formDataObj.buy_now_auction_price
            ? parseInt(formDataObj.buy_now_auction_price as string) || undefined
            : undefined
    }

    let validatedFields = null;

    if (detailsPart) {
        validatedFields = OfferDetailsFormSchema.safeParse(normalizedData);
    } else {
        validatedFields = CombinedOfferFormSchema.safeParse(normalizedData);
    }


    if (!validatedFields.success) {
        console.log("Add Offer validation errors:", validatedFields.error.flatten().fieldErrors);
        return {
            errors: validatedFields.error.flatten().fieldErrors,
            values: normalizedData as AddOfferFormState['values']
        };
    }

    if (detailsPart) {
        return {
            errors: {},
            values: normalizedData as AddOfferFormState['values']
        };
    }


    // convert model from "producer model" to "model"
    if ('model' in validatedFields.data) {
        const firstSpaceIndex = validatedFields.data.model.indexOf(' ');
        if (firstSpaceIndex !== -1) {
            validatedFields.data.model = validatedFields.data.model.substring(firstSpaceIndex + 1);
        }
    }

    // This property exists, the error is incorrect
    const { is_auction, ...offerData } = validatedFields.data;
    validatedFields.data = offerData;

    console.log("Is auction:", is_auction);
    console.log("Add Offer validated fields:", validatedFields.data);

    try {
        
    } catch (error) {

    }







    return state;

    // permanentRedirect("/");
    //TODO redirect to my offers
}