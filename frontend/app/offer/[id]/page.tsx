import { fetchOfferDetails } from "@/app/lib/data";
import Photos from "@/app/ui/offer/photos";
import { CarImageOfferPageSkeleton } from "@/app/ui/skeletons";
import { Suspense,lazy } from "react";

export default async function Page(props: { params: Promise<{id: string }> }) {
    const { params } = props;
    const { id } = await params;

    const offer = await fetchOfferDetails(id);

    return (
        <>
            <div className="my-3" />
            <div className="flex flex-row">
                <div className="mx-10" />
                <div className="flex flex-col gap-4">
                    <div className="md:h-120 md:w-200 border-0 border-gray-300">
                        {/* <Suspense fallback={<CarImageOfferPageSkeleton />}>
                            <Photos imagesURLs={offer.imagesURLs} />
                        </Suspense> */}
                        <Photos imagesURLs={offer.imagesURLs} />
                        {/* Add your offer details here */}
                    </div>
                </div>
            </div>
        </>
    );
}