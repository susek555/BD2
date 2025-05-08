export function OffersTableSkeleton() {
    return (
        <div>
        <div className="h-10 md:h-40 bg-gray-200 rounded animate-pulse w-15/16" />
        <div className="h-10 md:h-40 bg-gray-200 rounded animate-pulse w-15/16" />
        <div className="h-10 md:h-40 bg-gray-200 rounded animate-pulse w-15/16" />
        <div className="h-10 md:h-40 bg-gray-200 rounded animate-pulse w-15/16" />
        </div>
    );
}

export function OffersFoundSkeleton() {
    return (
        <div className="flex flex-row gap-4">
            <p className="font-bold">Offers found: </p>
            <div className="h-5 w-5 bg-gray-200 rounded animate-pulse" />
        </div>
    )
}

export function CarImageSkeleton() {
    return (
        <img src="/(home)/car_placeholder.png" alt="Car placeholder" className="h-35 w-70 object-cover rounded" />
    )
}