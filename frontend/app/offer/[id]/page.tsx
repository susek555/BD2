import { fetchOfferDetails } from "@/app/lib/data";
import Photos from "@/app/ui/offer/photos";
import { Suspense } from "react";

export default async function Page(props: { params: Promise<{id: string }> }) {
    const { params } = props;
    const { id } = await params;

    const offer = await fetchOfferDetails(id);

    return (
        <div>
            <h1>Offer ID: {id}</h1>
            <Suspense fallback={<div>Loading photos...</div>}>
                <Photos imagesURLs={offer.imagesURLs} />
            </Suspense>
            {/* Add your offer details here */}
        </div>
    );
}