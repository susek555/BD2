import { AuctionData } from "@/app/lib/definitions/reviews";
import { BasePriceButton } from "./price-buttons/base-price-button";
import { CurrencyDollarIcon, ArrowRightIcon } from "@heroicons/react/20/solid";
import TimeLeft from "./time-left";
import Link from "next/link";
import BidForm from "./bid-form";
import { authConfig } from "@/app/lib/authConfig";
import { getServerSession } from "next-auth/next";

type PriceData = {
    id: string;
    price: number;
    isAuction: boolean;
    auction?: AuctionData;
    isActive: boolean;
    priceOnly: boolean;
}

export default async function Price( data : { data : PriceData}) {
    const session = await getServerSession(authConfig);
    const loggedIn = !!session;
    const { id, price, isAuction, auction, isActive, priceOnly } = data.data;

    return (
        <>
            <div
                className={`flex flex-col gap-4 w-full md:w-120 border-2 p-4 ${
                isAuction ? "border-blue-500" : "border-gray-300"
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
                        {loggedIn ? (
                            !priceOnly ? (
                                <BidForm currentBid={auction!.currentBid} />
                            ) : (
                                <></>
                            )
                        ) : (
                            <Link href="/login">
                                <BasePriceButton>
                                    <p className="text-bold text-xl">Log in to Bid</p>
                                    <ArrowRightIcon className="ml-auto w-5 text-gray-50" />
                                </BasePriceButton>
                            </Link>
                        )}
                    </div>
                ) : (
                    <div className="flex justify-center items-center flex-col h-full gap-8">
                        <p className="font-bold text-3xl">{price.toString()} PLN</p>
                        {loggedIn ? (
                            !priceOnly ? (
                                <Link href={`/offer/${id}/buynow`}>
                                    <BasePriceButton>
                                        <p className="text-bold text-xl">Buy Now</p>
                                        <CurrencyDollarIcon className="ml-auto w-5 text-gray-50" />
                                    </BasePriceButton>
                                </Link>
                            ) : (
                                <></>
                            )
                        ) : (
                            <Link href="/login">
                                <BasePriceButton>
                                    <p className="text-bold text-xl">Log in to Buy</p>
                                    <ArrowRightIcon className="ml-auto w-5 text-gray-50" />
                                </BasePriceButton>
                            </Link>
                        )}
                    </div>
                )
                ) : (
                <div className="flex justify-center items-center flex-col h-full">
                    <p className="font-bold text-red-500 text-2xl">Offer is no longer active</p>
                </div>
                )}
            </div>
            {/* If the offer is an auction and has a price, show the buy now button */}
            { isAuction && price && isActive ? (
                <>
                    <div className="mt-4"></div>
                    <div className="flex flex-col gap-4 w-full md:w-120 p-4 border-blue-500 border-2">
                        <div className="flex justify-center items-center flex-col h-full gap-8">
                            <p className="font-bold text-3xl">{price.toString()} PLN</p>
                            {loggedIn ? (
                                !priceOnly ? (
                                    <Link href={`/offer/${id}/buynow`}>
                                        <BasePriceButton>
                                            <p className="text-bold text-xl">Buy Now</p>
                                            <CurrencyDollarIcon className="ml-auto w-5 text-gray-50" />
                                        </BasePriceButton>
                                    </Link>
                                ) : (
                                    <></>
                                )
                            ) : (
                                <Link href="/login">
                                    <BasePriceButton>
                                        <p className="text-bold text-xl">Log in to Buy</p>
                                        <ArrowRightIcon className="ml-auto w-5 text-gray-50" />
                                    </BasePriceButton>
                                </Link>
                            )}
                        </div>
                    </div>
                </>
            ) : ( <></>
            ) }
        </>
    );
}