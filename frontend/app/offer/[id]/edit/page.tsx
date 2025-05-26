import { fetchOfferFormData } from "@/app/lib/data/filters-sidebar/data";
import { OfferForm, offerActionEnum } from "@/app/ui/offer/add-offer-form";
import { Suspense } from "react";


export default async function Page(props: { params: Promise<{id: string }> }) {
    const { params } = props;
    const { id } = await params;

    const formData = await fetchOfferFormData();

    return (
        <div className="px-20 py-10">
            <h1 className="text-3xl font-bold pl-10">Edit Offer</h1>
            <div className="m-5"/>
            <Suspense fallback={<div>Loading...</div>}>
                <OfferForm inputsData = {formData} apiAction={offerActionEnum.EDIT_OFFER} id={id}/>
            </Suspense>
        </div>
    );
}