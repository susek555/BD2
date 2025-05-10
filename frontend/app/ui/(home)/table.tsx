import { Suspense } from "react"
import { CarImageSkeleton } from "../skeletons"
import { SearchParams } from "@/app/lib/definitions"
import { fetchOffers } from "@/app/lib/data";

export default async function OffersTable({ params } : { params: SearchParams }) {

    const offers = await fetchOffers(params);

    return (
        <div className="flex flex-col gap-4">
            {offers.map((offer) => (
                <a
                    key={offer.id}
                    href={`/${offer.id}`}
                    className={`h-10 md:h-40 bg-gray-200 rounded w-15/16 ${offer.isAuction ? 'border-4 border-blue-500' : ''}`}
                    style={{ borderColor: offer.isAuction ? 'blue' : 'transparent', borderWidth: offer.isAuction ? '2px' : '0' }}
                >
                    <div className="flex flex-row gap-4">
                        <div className="p-2.5">
                            <Suspense>
                                <CarImageSkeleton /> {/* TODO implement fetching image from internet */}
                            </Suspense>
                        </div>
                        <div className="flex py-2.5 h-40 flex-col w-full">
                            <div className="flex flex-row gap-10 items-center">
                                <p className="font-bold text-2xl whitespace-nowrap">{offer.name}</p>
                                {offer.isAuction && (
                                    <div className="flex justify-end w-full">
                                        <p className="font-bold text-2xl whitespace-nowrap px-3">Auction</p>
                                    </div>
                                )}
                            </div>
                            <div className="px-4 py-3 flex flex-row gap-10">
                                <p className="text-bg">Production year: <span className="font-bold">{offer.productionYear.toString()}</span></p>
                                <p className="text-bg">Mileage: <span className="font-bold">{offer.mileage.toString()} km</span></p>
                            </div>
                            <div className="px-4 flex flex-row gap-10">
                                <p className="text-bg">Color: <span className="font-bold">{offer.color}</span></p>
                            </div>
                            <div className="flex justify-end items-end h-full self-end">
                                <p className="px-3 font-bold text-2xl">{offer.price.toString()} zł</p>
                            </div>
                        </div>
                    </div>
                </a>
            ))}
        </div>
    )
}

