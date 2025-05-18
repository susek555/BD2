"use client";

import React, { useState } from "react";
import { usePathname, useSearchParams } from "next/navigation";

export default function OfferType({ offerTypes }: { offerTypes: string[] }) {
    const searchParams = useSearchParams();
    const pathname = usePathname();

    const [isDropdownOpen, setIsDropdownOpen] = useState(false);
    const options = offerTypes;
    const selectedOption = searchParams.get("offerType") || "All";

    const toggleDropdown = () => {
        setIsDropdownOpen(!isDropdownOpen);
    };

    const handleOptionChange = (option: string) => {
        const params = new URLSearchParams(searchParams);
        if (option === "All") {
            params.delete("offerType");
        } else {
            params.set("offerType", option);
        }
        params.set("page", "1"); // Reset to the first page
        window.history.replaceState(null, '', `${pathname}?${params.toString()}`);
        setIsDropdownOpen(false);
};

return (
        <div>
            <h2 className="text-bg px-2">Offer Type:</h2>
            <div style={{ display: 'flex', gap: '1rem' }}>
                <div
                    className="base-filter-template border border-black-300 rounded px-2 py-1 relative"
                    style={{ flex: 2 }}
                >
                    <button
                        className="flex justify-between items-center w-full"
                        onClick={toggleDropdown}
                    >
                        <span>{selectedOption}</span>
                        <span>{isDropdownOpen ? '▲' : '▼'}</span>
                    </button>
                    {isDropdownOpen && (
                        <div
                            className="filter-options mt-2 absolute bg-white border border-gray-300 rounded shadow-md w-full z-10"
                        >
                            {options.map((option) => (
                                <div
                                    key={option}
                                    className="filter-option flex justify-between items-center px-2 py-1 hover:bg-gray-100 cursor-pointer"
                                    onClick={() => handleOptionChange(option)}
                                >
                                    <span>{option}</span>
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
