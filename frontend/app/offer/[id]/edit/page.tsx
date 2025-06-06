import { fetchOfferFormData } from "@/app/lib/data/filters-sidebar/data";
import { fetchOfferDetails } from "@/app/lib/data/offer/data";
import { editFormWrapper } from "@/app/lib/definitions/offer-form";
import { OfferForm, offerActionEnum } from "@/app/ui/offer/add-offer-form";
import { Suspense } from "react";


export default async function Page(props: { params: Promise<{id: string }> }) {
    const { params } = props;
    const { id } = await params;

    const offer = await fetchOfferDetails(id);
    if (!offer) {
    throw new Error('Offer not found');
    }
    if (!offer.can_edit) {
        throw new Error('You are not allowed to edit this offer');
    }

    const formData = await fetchOfferFormData();

    //TODO fetch initial data from the API - probably need a wrapper

    const initialData = editFormWrapper(offer);

    const imagesURLs = offer.imagesURLs || [];

    console.log('initialData', initialData);

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