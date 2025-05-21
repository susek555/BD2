import { BellAlertIcon } from "@heroicons/react/20/solid";
import { BaseAccountButton } from "@/app/ui/(topbar)/base-account-buttons/base-account-button";

export function NotificationsButtonSkeleton() {
    return (
        <>
            <BaseAccountButton
                className="w-12 flex items-center justify-center relative"
            >
                <BellAlertIcon className="w-5 text-gray-50" />
                <div className="absolute -top-1 -right-1 bg-gray-200 text-white rounded-full w-5 h-5 flex items-center justify-center">
                    <div className="w-3 h-3 border-2 border-gray-400 border-t-gray-600 rounded-full animate-spin"></div>
                </div>
            </BaseAccountButton>
        </>
    )
}

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

export function SideBarSkeleton() {
    return (
        <div className="flex h-full flex-col px-3 py-2 md:px-2 rounded-r-lg border-black border-[2px]">
            <div className="flex grow flex-row justify-between space-x-2 md:flex-col md:space-x-0 md:space-y-2">
                <div className="h-20 w-full bg-gray-200 rounded animate-pulse" /> {/* OfferType */}
                <div className="h-20 w-full bg-gray-200 rounded animate-pulse" /> {/* Sorting */}
                <div className="h-40 w-full bg-gray-200 rounded animate-pulse" /> {/* ProducersAndModels */}
                <div className="h-60 w-full bg-gray-200 rounded animate-pulse" /> {/* Filters */}
                <div className="h-10 w-full bg-gray-200 rounded animate-pulse" /> {/* ApplyButton */}
                <div className="hidden h-auto w-full grow rounded-md bg-gray-50 md:block"></div>
            </div>
        </div>
    )
};

export function CarImageHomePageSkeleton() {
    return (
        <img src="/(home)/car_placeholder.png" alt="Car placeholder" className="h-35 w-70 object-cover rounded" />
    )
}

export function CarImageOfferPageSkeleton() {
    return (
        <div className="flex flex-col items-center gap-4">
            <div className="w-full flex flex-col items-center border border-gray-300 relative bg-gray-100">
                <div className="md:h-130 relative w-full max-w-lg overflow-hidden flex justify-center items-center">
                    <div className="loader animate-spin rounded-full border-4 border-gray-300 border-t-gray-600 w-12 h-12"></div>
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

export function NameOfferPageSkeleton() {
    return (
        <div className="h-10 w-70 bg-gray-200 rounded animate-pulse" />
    )
}

export function PriceOfferPageSkeleton() {
    return (
        <div className="flex flex-col gap-4 w-full md:w-120 h-full md:h-63 border border-gray-300">
            <div className="h-full w-full bg-gray-200 animate-pulse" />
        </div>
    )
}

export function OfferDescriptionSkeleton() {
    return (
        <div className="flex flex-col gap-4 w-full h-full md:h-80 border border-gray-300">
            <div className="h-full w-full bg-gray-200 animate-pulse" />
        </div>
    )
}

