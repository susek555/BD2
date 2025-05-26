import { fetchAddOfferFormData } from "@/app/lib/data/filters-sidebar/data";
import OfferForm from "@/app/ui/offer/add-offer-form";

export default async function Page() {

    const formData = await fetchAddOfferFormData();

    return (
        <div className="px-20 py-10">
            <h1 className="text-3xl font-bold pl-10">Add Offer</h1>
            <div className="m-5"/>
            <OfferForm inputsData = {formData}/>
        </div>
    );
}