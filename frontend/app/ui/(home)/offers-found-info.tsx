export default function OffersFoundInfo({ totalOffers} : {totalOffers: number}) {
    return(
        <div className="flex flex-row gap-4">
            <p className="font-bold">Offers found: </p>
            <p className="font-bold">{totalOffers.toString()} </p>
        </div>
    )
}