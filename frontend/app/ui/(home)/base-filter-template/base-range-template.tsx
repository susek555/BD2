"use client";

import React, { useState } from 'react';

interface BaseRangeTemplateProps {
    name: string;
    onChange: (name:string, range: { min: number; max: number }) => void;
}

export function BaseRangeTemplate({ name, onChange }: BaseRangeTemplateProps) {
    const [isOpen, setIsOpen] = useState(false);
    const [range, setRange] = useState<{ min: number; max: number }>({ min: 0, max: 0 });

    const toggleDropdown = () => setIsOpen(!isOpen);

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        const updatedRange = { ...range, [name]: Number(value) };
        setRange(updatedRange);
        onChange(name, updatedRange);
    };

    return (
        <div className="base-range-template border border-black-300 rounded px-2 py-1">
            <button
            className="flex justify-between items-center w-full"
            onClick={toggleDropdown}
            >
            <span>{name}</span>
            <span>{isOpen ? '▲' : '▼'}</span>
            </button>
            {isOpen && (
                <div className="range-inputs mt-2 flex items-center justify-between w-full">
                    <input
                        type="number"
                        id="min"
                        name="min"
                        placeholder="Min"
                        value={range.min || ''}
                        onChange={handleInputChange}
                        className="border border-gray-300 rounded px-1 py-1 w-full"
                        step="1"
                        inputMode="numeric"
                    />
                    <span className="mx-2 text-center flex-shrink-0 px-2">-</span>
                    <input
                        type="number"
                        id="max"
                        name="max"
                        placeholder="Max"
                        value={range.max || ''}
                        onChange={handleInputChange}
                        className="border border-gray-300 rounded px-1 py-1 w-full"
                        step="1"
                        inputMode="numeric"
                    />
                </div>
            )}
        </div>
    );
}