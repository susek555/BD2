'use client'

import React, { useState } from 'react';
import { BasePriceButton } from './price-buttons/base-price-button';
import { CurrencyDollarIcon } from '@heroicons/react/20/solid';

const BidForm: React.FC = () => {
    const [bidValue, setBidValue] = useState('');

    const handleSubmit = (event: React.FormEvent) => {
        event.preventDefault();
        console.log('Bid submitted:', bidValue);
        // Add logic to handle the bid submission
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
        </form>
    );
};

export default BidForm;