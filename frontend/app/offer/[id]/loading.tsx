// app/offer/[id]/loading.tsx
import { CarImageOfferPageSkeleton } from "@/app/ui/skeletons";

export default function Loading() {
    return (
        <>
            <div className="my-3" />
            <div className="flex flex-row">
                <div className="mx-10" />
                <div className="flex flex-col gap-4">
                    <div className="md:h-120 md:w-200 border-0 border-gray-300">
                        <CarImageOfferPageSkeleton />
                    </div>
                </div>
            </div>
        </>
    );
}
