'use client'

import React, { useState } from 'react';
import { BasePriceButton } from './price-buttons/base-price-button';
import { CurrencyDollarIcon } from '@heroicons/react/20/solid';

export default function BidForm({ currentBid }: { currentBid: number }) {
    const [bidValue, setBidValue] = useState('');
    const [errors, setErrors] = useState<string | null>(null);
    const [showConfirmDialog, setShowConfirmDialog] = useState(false);

    const handleSubmit = (event: React.FormEvent) => {
        event.preventDefault();
        setErrors(null);

        if (parseFloat(bidValue) > currentBid) {
            setShowConfirmDialog(true);
        } else {
            setErrors('Your bid must be higher than the current bid.');
        }
    };

    const confirmBid = () => {
        console.log('Bid is valid:', bidValue);
        // Add logic to process the valid bid
        setShowConfirmDialog(false);
    };

    const cancelBid = () => {
        setShowConfirmDialog(false);
    };

    return (
        <>
            <form
                onSubmit={handleSubmit}
                className="flex flex-col items-center gap-2"
            >
                <label>
                    <input
                        type="number"
                        value={bidValue}
                        onChange={(e) => setBidValue(e.target.value)}
                        placeholder="Enter your bid"
                        required
                        className="border md:w-80 border-gray-300 p-2 rounded"
                    />
                </label>
                <BasePriceButton>
                    <p className="font-bold text-xl">Bid</p>
                    <CurrencyDollarIcon className="ml-auto w-5 text-gray-50" />
                </BasePriceButton>
                {errors && (
                    <div
                        id="bid-error"
                        aria-live="polite"
                        aria-atomic="true"
                        className="mt-1 text-sm text-red-500"
                    >
                        {errors}
                    </div>
                )}
            </form>

            {showConfirmDialog && (
                <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50">
                    <div className="bg-white p-6 rounded shadow-lg">
                        <p className="mb-4">Are you sure you want to place a bid of ${bidValue}?</p>
                        <div className="flex justify-end gap-2">
                            <button
                                onClick={cancelBid}
                                className="px-4 py-2 bg-gray-300 rounded"
                            >
                                Cancel
                            </button>
                            <button
                                onClick={confirmBid}
                                className="px-4 py-2 bg-blue-500 text-white rounded"
                            >
                                Confirm
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </>
    );
};