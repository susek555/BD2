export default async function FoundInfo({ title, totalOffers: totalCount } : { title: string, totalOffers: number }) {
    return(
        <div className="flex flex-row gap-4">
            <p className="font-bold">{title}: </p>
            <p className="font-bold">{totalCount.toString()} </p>
        </div>
    )
}