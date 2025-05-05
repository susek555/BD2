"use client";

import React, { useState } from 'react';

interface BaseFilterTemplateProps {
    name: string;
    options: string[];
    onChange: (selected: string[]) => void;
}

export function BaseFilterTemplate({ name, options, onChange }: BaseFilterTemplateProps) {
    const [isOpen, setIsOpen] = useState(false);
    const [selectedOptions, setSelectedOptions] = useState<string[]>([]);

    const toggleDropdown = () => setIsOpen(!isOpen);

    const handleCheckboxChange = (option: string) => {
        const updatedSelectedOptions = selectedOptions.includes(option)
            ? selectedOptions.filter((item) => item !== option)
            : [...selectedOptions, option];

        setSelectedOptions(updatedSelectedOptions);
        onChange(updatedSelectedOptions);
    };

    return (
        <div className="base-filter-template border border-black-300 rounded px-2 py-1">
            <button
            className="flex justify-between items-center w-full"
            onClick={toggleDropdown}
            >
            <span>{name}</span>
            <span>{isOpen ? '▲' : '▼'}</span>
            </button>
            {isOpen && (
            <div className="filter-options mt-2">
                {options.map((option) => (
                <div
                    key={option}
                    className="filter-option flex justify-between items-center border border-gray-300 rounded px-2 py-1 mb-1"
                    onClick={() => handleCheckboxChange(option)}
                >
                    <span>{option}</span>
                    <input
                    type="checkbox"
                    className="w-6 h-6 accent-green-500"
                    checked={selectedOptions.includes(option)}
                    onChange={(e) => e.stopPropagation()}
                    />
                </div>
                ))}
            </div>
            )}
        </div>
    );
}