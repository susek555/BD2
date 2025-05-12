import { AuctionData } from "@/app/lib/definitions";
import { BasePriceButton } from "./price-buttons/base-price-button";
import { CurrencyDollarIcon } from "@heroicons/react/20/solid";

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
                isAuction ? (
                    <p>Auction</p>
                ) : (
                    <div className="flex justify-center items-center flex-col h-full gap-8">
                        <p className="font-bold text-3xl">{price.toString()} PLN</p>
                        <BasePriceButton>
                            <p className="text-bold text-xl">Buy Now</p>
                            <CurrencyDollarIcon className="ml-auto w-5 text-gray-50" />
                        </BasePriceButton>
                    </div>
                )
            ) : (
                <div className="flex justify-center items-center flex-col h-full">
                    <p className="font-bold text-red-500 text-2xl">Offer is no longer active</p>
                </div>
            )}
        </div>
    );
}