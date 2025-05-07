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
                            <CarImageSkeleton /> {/* TODO implement fetching image from internet */}
                        </Suspense>
                    </div>
                    <div className="flex py-2.5 h-40 flex-col w-full">
                        <p className="font-bold text-2xl">Car name</p>
                        <div className="px-4 py-3 flex flex-row gap-10">
                            <p className="text-bg">Production year: 2000</p>
                            <p className="text-bg">Mileage: 150000 km</p>
                        </div>
                        <div className="px-4 flex flex-row gap-10">
                            <p className="text-bg">First owner</p>
                        </div>
                        <div className="flex justify-end items-end h-full self-end">
                            <p className="px-3 font-bold text-2xl">10000 zł</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}