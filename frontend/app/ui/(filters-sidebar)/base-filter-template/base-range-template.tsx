import React, { useState } from 'react';

interface BaseRangeTemplateProps {
    fieldName: string;
    rangeBase: {
        min: number | null;
        max: number | null;
    };
    onChange: (name: string, range: { min: number; max: number }) => void;
}

export function BaseRangeTemplate({ fieldName, rangeBase, onChange }: BaseRangeTemplateProps) {
    const [isOpen, setIsOpen] = useState(false);
    const [range, setRange] = useState<{ min: number; max: number }>({ min: rangeBase.min || 0, max: rangeBase.max || 0 });

    const toggleDropdown = () => setIsOpen(!isOpen);

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        const updatedRange = { ...range, [name]: Number(value) };
        setRange(updatedRange);
        onChange(fieldName, updatedRange);
    };

    return (
        <div className="base-range-template border border-black-300 rounded px-2 py-1">
            <button
                className="flex justify-between items-center w-full"
                onClick={toggleDropdown}
            >
                <span>{fieldName}</span>
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
            {range.min !== 0 && range.max !== 0 && range.max < range.min && (
                <p className="text-red-500 text-sm mt-1">Max value cannot be less than Min value.</p>
            )}
        </div>
    );
}