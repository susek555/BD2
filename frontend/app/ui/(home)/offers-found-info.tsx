export default function OffersFoundInfo({ totalOffers} : {totalOffers: number}) {
    return(
        <div className="flex flex-row gap-4">
            <p className="font-bold">Offers found: </p>
            <p className="font-bold">{totalOffers.toString()} </p>
            {/* <div className="h-5 w-5 bg-gray-200 rounded animate-pulse" /> */}
        </div>
    )
}