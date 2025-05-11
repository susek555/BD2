export function OffersTableSkeleton() {
    return (
        <div className="flex flex-col gap-4">
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

export function CarImageHomePageSkeleton() {
    return (
        <img src="/(home)/car_placeholder.png" alt="Car placeholder" className="h-35 w-70 object-cover rounded" />
    )
}

export function CarImageOfferPageSkeleton() {
    return (
        <div className="flex flex-col items-center gap-4">
            <div className="w-full flex flex-col items-center border border-gray-300" >
                <div className="relative w-full max-w-lg overflow-hidden">
                    <img src="/(home)/car_placeholder.png" alt="Car placeholder" className="w-full h-auto md:h-120 object-cover rounded scale-140" />
                </div>
            </div>
            <div className="flex items-center gap-2">
                <button
                    className="bg-gray-800 text-white px-2 py-1 rounded"
                >
                    &#8592;
                </button>
                <p>
                    0 / 0
                </p>
                <button
                    className="bg-gray-800 text-white px-2 py-1 rounded"
                >
                    &#8594;
                </button>
            </div>
        </div>
    )
}

