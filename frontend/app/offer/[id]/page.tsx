import { fetchOfferDetails } from "@/app/lib/data";
import Photos from "@/app/ui/offer/photos";
import Price from "@/app/ui/offer/price";
import { Suspense,lazy } from "react";

export default async function Page(props: { params: Promise<{id: string }> }) {
    const { params } = props;
    const { id } = await params;

    const offer = await fetchOfferDetails(id);

    return (
        <>
            <div className="my-3" />
                <div className="flex flex-row gap-10 md:px-20">
                    <div className="flex flex-col gap-4 w-full">
                        <div className="w-full h-full border-0 border-gray-300">
                            <div className="w-full static">
                                <Photos imagesURLs={offer.imagesURLs} />
                            </div>
                        </div>
                    </div>
                    <div className="md:w-120 h-full">
                        <div className="p-4">
                            <h1 className="text-3xl font-bold">{offer.name}</h1>
                        </div>
                        <div className="my-50" />
                        <Price data={{ price: offer.price, isAuction: offer.isAuction, auction: offer.auctionData,  isActive: offer.isActive }} />
                    </div>
            </div>
        </>
    );
}