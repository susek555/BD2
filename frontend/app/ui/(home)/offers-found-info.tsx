import { fetchTotalOffers } from "@/app/lib/data";
import { SearchParams } from "@/app/lib/definitions";

export default async function OffersFoundInfo({ params } : { params: SearchParams }) {
    const totalOffers = await fetchTotalOffers(params);

    return(
        <div className="flex flex-row gap-4">
            <p className="font-bold">Offers found: </p>
            <p className="font-bold">{totalOffers.toString()} </p>
        </div>
    )
}