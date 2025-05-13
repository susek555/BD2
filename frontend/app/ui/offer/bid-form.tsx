'use client'

import React, { useState } from 'react';
import { BasePriceButton } from './price-buttons/base-price-button';
import { CurrencyDollarIcon } from '@heroicons/react/20/solid';

export default function BidForm({ currentBid }: { currentBid: number }) {
    const [bidValue, setBidValue] = useState('');
    const [errors, setErrors] = useState<string | null>(null);

    const handleSubmit = (event: React.FormEvent) => {
        event.preventDefault();
        setErrors(null);

        if (parseFloat(bidValue) > currentBid) {
            console.log('Bid is valid:', bidValue);
            // Add logic to process the valid bid
        } else {
            setErrors('Your bid must be higher than the current bid.');
        }
    };

    return (
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
    );
};