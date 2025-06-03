'use client'

import React, { useState } from 'react';
import { BasePriceButton } from './price-buttons/base-price-button';
import { CurrencyDollarIcon } from '@heroicons/react/20/solid';
import ConfirmationModal from '../../(common)/confirm-modal';
import { useRouter } from 'next/navigation';
import { placeBid } from '@/app/lib/api/offer/bid';

export default function BidForm({ currentBid }: { currentBid: number }) {
    const router = useRouter();

    const [bidValue, setBidValue] = useState('');
    const [errors, setErrors] = useState<string | null>(null);
    const [showConfirmDialog, setShowConfirmDialog] = useState(false);

    const handleSubmit = (event: React.FormEvent) => {
        event.preventDefault();
        setErrors(null);

        if (parseInt(bidValue) > currentBid) {
            setShowConfirmDialog(true);
        } else {
            setErrors('Your bid must be higher than the current bid.');
        }
    };

    const confirmBid = async () => {
        console.log('Bid is valid:', bidValue);
        //TODO: Add logic to process the valid bid

        try {
            await placeBid({
                amount: parseInt(bidValue),
                auction_id: parseInt(window.location.pathname.split('/').pop() || '0'),
            });
            setShowConfirmDialog(false);
            alert('Your bid successful!');

        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        } catch (error) {
            alert('An error occurred while processing your bid. Please try again.');
        }
        router.refresh();
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

            <ConfirmationModal
                title='Confirm Bid'
                message='Are you sure you want to bid this offer? This action cannot be undone.'
                confirmText='Bid'
                onConfirm={confirmBid}
                onCancel={() => setShowConfirmDialog(false)}
                isOpen={showConfirmDialog}
                bg_color='bg-blue-500'
                bg_color_hover='bg-blue-600'
            />
        </>
    );
};