import { AuctionData } from "@/app/lib/definitions";
import { BasePriceButton } from "./price-buttons/base-price-button";
import { CurrencyDollarIcon } from "@heroicons/react/20/solid";
import TimeLeft from "./time-left";
import Link from "next/link";
import BidForm from "./bid-form";

type PriceData = {
    price: number;
    isAuction: boolean;
    auction?: AuctionData;
    isActive: boolean;
}

export default function Price( data : { data : PriceData}) {
    const { price, isAuction, auction, isActive } = data.data;

    return (
        <>
            <div
                className={`flex flex-col gap-4 w-full md:w-120 h-full md:h-63 ${
                isAuction ? "border-blue-500 border-2" : "border-gray-300 border"
                }`}
            >
                {isActive ? (
                isAuction ? (
                    <div className="flex justify-center items-center flex-col h-full gap-8">
                        <div className="flex flex-row gap-20">
                            <TimeLeft endDate={auction!.endDate}/>
                            <div className="flex flex-col gap-2 justify-center items-center">
                                <p className="text-2xl">Current bid</p>
                                <p className="font-bold text-2xl">{auction!.currentBid.toString()} PLN</p>
                            </div>
                        </div>
                        <BidForm currentBid={auction!.currentBid} />
                    </div>
                ) : (
                    <div className="flex justify-center items-center flex-col h-full gap-8">
                        <p className="font-bold text-3xl">{price.toString()} PLN</p>
                        <Link href="/account">
                            <BasePriceButton>
                                <p className="text-bold text-xl">Buy Now</p>
                                <CurrencyDollarIcon className="ml-auto w-5 text-gray-50" />
                            </BasePriceButton>
                        </Link>
                    </div>
                )
                ) : (
                <div className="flex justify-center items-center flex-col h-full">
                    <p className="font-bold text-red-500 text-2xl">Offer is no longer active</p>
                </div>
                )}
            </div>
            { isAuction && price ? (
                <>
                    <div className="mt-4"></div>
                    <div className="flex flex-col gap-4 w-full md:w-120 h-full md:h-63 border-blue-500 border-2">
                        <div className="flex justify-center items-center flex-col h-full gap-8">
                            <p className="font-bold text-3xl">{price.toString()} PLN</p>
                            <Link href="/account">
                                <BasePriceButton>
                                    <p className="text-bold text-xl">Buy Now</p>
                                    <CurrencyDollarIcon className="ml-auto w-5 text-gray-50" />
                                </BasePriceButton>
                            </Link>
                        </div>
                    </div>
                </>
            ) : ( <></>
            ) }
        </>
    );
}