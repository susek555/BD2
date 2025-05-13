// app/offer/[id]/loading.tsx
import { CarImageOfferPageSkeleton, NameOfferPageSkeleton, OfferDescriptionSkeleton, PriceOfferPageSkeleton } from "@/app/ui/skeletons";

export default function Loading() {
    return (
        <>
            <div className="my-3" />
            <div className="flex flex-row gap-10 md:px-20">
                <div className="flex flex-col gap-4 w-full">
                    <div className="w-full h-full border-0 border-gray-300">
                        <CarImageOfferPageSkeleton />
                    </div>
                    <OfferDescriptionSkeleton />
                </div>
                <div className="md:w-120 h-full md:h-130">
                    <div className="p-4">
                        <NameOfferPageSkeleton />
                    </div>
                    <div className="my-10" />
                    <PriceOfferPageSkeleton />
                    <div className="my-4" />
                    <PriceOfferPageSkeleton />
                </div>
            </div>

        </>
    );
}
