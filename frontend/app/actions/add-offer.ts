import { AddOfferFormSchema, AddOfferFormState } from "../lib/definitions";
import { permanentRedirect } from "next/navigation";

export async function addOffer(
    state: AddOfferFormState,
    formData: FormData
): Promise<AddOfferFormState> {
    console.log("Add Offer form data:", Object.fromEntries(formData.entries()));

    const formDataObj = Object.fromEntries(formData.entries());
    const validatedFields = AddOfferFormSchema.safeParse(formDataObj);
    console.log("Add Offer validation result:", validatedFields);


    if (!validatedFields.success) {
        console.log("Add Offer validation errors:", validatedFields.error.flatten().fieldErrors);
        return {
            errors: validatedFields.error.flatten().fieldErrors,
            values: formDataObj as AddOfferFormState['values']
        };
    }


    return state;

    // permanentRedirect("/");
    //TODO redirect to my offers
}

export function setProducer(
    producer: string,
    newProducer: string,
): string {
    return newProducer;
}