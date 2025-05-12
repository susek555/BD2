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
        <div className="flex flex-col">
            <div className="text-2xl font-bold text-gray-900">
                {isAuction ? "Auction" : "Price"}: {price} €
            </div>
            {isAuction && (
                <div className="text-sm text-gray-500">
                    Auction ends on: {auction && auction.endDate.toLocaleDateString()}
                </div>
            )}
        </div>
    );
}