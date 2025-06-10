'use client'

import { BasePriceButton } from "@/app/ui/offer/[id]/price-buttons/base-price-button";
import Link from "next/link";
import { useState } from "react";
import { PencilSquareIcon, TrashIcon } from "@heroicons/react/20/solid";
import { deleteListingAction } from "@/app/actions/listing-actions";
import { permanentRedirect } from "next/navigation";

export default function OwnerView({ can_edit, can_delete, offer_id, isAuction }: { can_edit: boolean, can_delete: boolean, offer_id: string, isAuction: boolean }) {
    const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

    async function handleDelete() {
        console.log("Delete offer with ID:", offer_id);
        await deleteListingAction(offer_id);
        alert("Offer deleted successfully. Now you will be redirected to your listings page.");
        permanentRedirect("/account/listings");
        setShowDeleteConfirm(false);
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
                        <BasePriceButton onClick={() => setShowDeleteConfirm(true)}>
                            <p className="text-bold text-xl">Delete</p>
                            <TrashIcon className="ml-auto w-5 text-gray-50" />
                        </BasePriceButton>
                    )}
                </div>
            </div>

            {/* Delete Confirmation Modal */}
            {showDeleteConfirm && (
                <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
                    <div className="bg-white p-6 rounded-lg shadow-lg max-w-md w-full">
                        <h2 className="text-xl font-bold mb-4">Confirm Delete</h2>
                        <p className="mb-6">Are you sure you want to delete this offer? This action cannot be undone.</p>
                        <div className="flex justify-end gap-3">
                            <button
                                onClick={() => setShowDeleteConfirm(false)}
                                className="px-4 py-2 rounded bg-gray-200 hover:bg-gray-300"
                            >
                                Cancel
                            </button>
                            <button
                                onClick={handleDelete}
                                className="px-4 py-2 rounded bg-red-500 text-white hover:bg-red-600"
                            >
                                Delete
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </>
    )
}