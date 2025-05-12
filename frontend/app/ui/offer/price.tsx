import { AuctionData } from "@/app/lib/definitions";

type PriceData = {
    price: number;
    isAuction: boolean;
    auction?: AuctionData;
    isActive: boolean;
}

export default function Price( data : { data : PriceData}) {
    const { price, isAuction, auction, isActive } = data.data;

    return (
        <div className="flex flex-col gap-4 w-full md:w-120 h-full md:h-63 border border-gray-300">
            {isActive ? (
                <div className="p-4">
                    {isAuction && auction ? (
                        <div>
                            <p className="text-lg font-bold">Auction Price: ${auction.currentBid}</p>
                            <p className="text-sm text-gray-500">Ends on: {new Date(auction.endDate).toLocaleString()}</p>
                        </div>
                    ) : (
                        <p className="text-lg font-bold">Price: ${price}</p>
                    )}
                </div>
            ) : (
                <div className="p-4 flex justify-center items-center flex-col h-full">
                    <p className="font-bold text-red-500 text-2xl">Offer is no longer active</p>
                </div>
            )}
        </div>
    );
}