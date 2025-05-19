'use client'

import { PencilSquareIcon } from "@heroicons/react/20/solid"
import { BasePriceButton } from "@/app/ui/offer/[id]/price-buttons/base-price-button";
import Link from "next/link";

export default function OwnerView( { can_edit, can_delete, offer_id, isAuction} : {can_edit: boolean, can_delete: boolean, offer_id: string, isAuction: boolean}) {
    function handleDelete() {
        console.log("Delete offer with ID:", offer_id);
    }

    return (
        <>
            <div
                className={`flex flex-col gap-4 w-full md:w-120 border-2 p-4 ${
                isAuction ? "border-blue-500" : "border-gray-300"
                }`}
            >
                <div className="flex justify-center items-center flex-col h-full gap-5">
                    <p className="text-2xl">This is your offer</p>
                    {can_edit && (
                        <Link href={`/offer/${offer_id}/edit`}>
                            <BasePriceButton>
                                <p className="text-bold text-xl">Edit</p>
                                <PencilSquareIcon className="ml-auto w-5 text-gray-50" />
                            </BasePriceButton>
                        </Link>
                    )}
                    {can_delete && (
                        <button onClick={() => handleDelete()}>
                            <BasePriceButton>
                                <p className="text-bold text-xl">Delete</p>
                                <PencilSquareIcon className="ml-auto w-5 text-gray-50" />
                            </BasePriceButton>
                        </button>
                    )}
                </div>
            </div>
        </>

    )
}