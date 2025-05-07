import { Suspense } from "react"
import { CarImageSkeleton } from "../skeletons"

export default async function OffersTable() {
    // const searchParams = useSearchParams();
    return (
        <div className="flex flex-col gap-4">
            <div className="h-10 md:h-40 bg-gray-200 rounded w-15/16">
                <div className="flex flex-row gap-4">
                    <div className="p-2.5">
                        <Suspense>
                            <CarImageSkeleton />
                        </Suspense>
                    </div>
                </div>
            </div>
        </div>
    )
}