import { fetchOfferFormData } from "@/app/lib/data/filters-sidebar/data";
import { OfferFormState } from "@/app/lib/definitions/offer-form";
import { OfferForm, offerActionEnum } from "@/app/ui/offer/add-offer-form";
import { Suspense } from "react";


export default async function Page(props: { params: Promise<{id: string }> }) {
    const { params } = props;
    const { id } = await params;

    const formData = await fetchOfferFormData();

    //TODO fetch initial data from the API - probably need a wrapper

    const initialData: Partial<OfferFormState['values']> = {
        manufacturer: "Audi",
        model: "Audi A3",
        color: "Black",
        fuel_type: "Petrol",
        transmission: "Manual",
        drive: "FWD",
        production_year: 2020,
        mileage: 15000,
        number_of_doors: 4,
        number_of_seats: 5,
        number_of_gears: 6,
        engine_power: 150,
        engine_capacity: 2000,
        registration_date: "2020-01-01",
        registration_number: "ABC123",
        vin: "WVWZZZ1JZ9W123456",
        description: "This is a test offer description.",
        price: 20000,
        is_auction: false,
        margin: 8,
    }

    return (
        <div className="px-20 py-10">
            <h1 className="text-3xl font-bold pl-10">Edit Offer</h1>
            <div className="m-5"/>
            <Suspense fallback={<div>Loading...</div>}>
                <OfferForm
                    inputsData = {formData}
                    initialValues={initialData}
                    apiAction={offerActionEnum.EDIT_OFFER}
                />
            </Suspense>
        </div>
    );
}