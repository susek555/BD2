"use client";

import React, { useState } from "react";
import { usePathname, useSearchParams } from "next/navigation";

export default function Sorting( {sortingOptions}: { sortingOptions: string[] }) {
    const searchParams = useSearchParams();
    const pathname = usePathname();

    const [isKeyDropdownOpen, setIsKeyDropdownOpen] = useState(false);
    const [isOrderDropdownOpen, setIsOrderDropdownOpen] = useState(false);

    const selectedKey = searchParams.get("sortKey") || "Base";
    const isDescSelected = searchParams.get("isSortDesc") || false;

    const toggleKeyDropdown = () => setIsKeyDropdownOpen(!isKeyDropdownOpen);
    const toggleOrderDropdown = () => setIsOrderDropdownOpen(!isOrderDropdownOpen);

    const handleKeyChange = (key: string) => {
        const params = new URLSearchParams(searchParams);
        if (key === "Base") {
            params.delete("sortKey");
        } else {
            params.set("sortKey", key);
        }
        params.set("page", "1"); // Reset to the first page
        window.history.replaceState(null, '', `${pathname}?${params.toString()}`);
        setIsKeyDropdownOpen(false);
    };

    const handleOrderChange = (order: string) => {
        const params = new URLSearchParams(searchParams);
        if (order === "Desc") {
            params.set("isSortDesc", "true");
        } else {
            params.delete("isSortDesc");
        }
        params.set("page", "1"); // Reset to the first page
        window.history.replaceState(null, '', `${pathname}?${params.toString()}`);
        setIsOrderDropdownOpen(false);
    };

    return (
        <div>
            <h2 className='mb-2 text-sm font-semibold text-gray-700'>Order by</h2>
            <div style={{ display: 'flex', gap: '1rem' }}>
                <div
                    className="base-filter-template border border-black-300 rounded px-2 py-1 relative"
                    style={{ flex: 2 }}
                >
                    <button
                        className="flex justify-between items-center w-full"
                        onClick={toggleKeyDropdown}
                    >
                        <span>{selectedKey}</span>
                        <span>{isKeyDropdownOpen ? '▲' : '▼'}</span>
                    </button>
                    {isKeyDropdownOpen && (
                        <div
                            className="filter-options mt-2 absolute bg-white border border-gray-300 rounded shadow-md w-full z-10"
                        >
                            {sortingOptions.map((key) => (
                                <div
                                    key={key}
                                    className="filter-option flex justify-between items-center px-2 py-1 hover:bg-gray-100 cursor-pointer"
                                    onClick={() => handleKeyChange(key)}
                                >
                                    <span>{key}</span>
                                </div>
                            ))}
                        </div>
                    )}
                </div>
                <div
                    className="base-filter-template border border-black-300 rounded px-2 py-1 relative"
                    style={{ flex: 1 }}
                >
                    <button
                        className="flex justify-between items-center w-full"
                        onClick={toggleOrderDropdown}
                    >
                        <span>{isDescSelected === "true" ? "Desc" : "Asc"}</span>
                        <span>{isOrderDropdownOpen ? '▲' : '▼'}</span>
                    </button>
                    {isOrderDropdownOpen && (
                        <div
                            className="filter-options mt-2 absolute bg-white border border-gray-300 rounded shadow-md w-full z-10"
                        >
                            {["Asc", "Desc"].map((order) => (
                                <div
                                    key={order}
                                    className="filter-option flex justify-between items-center px-2 py-1 hover:bg-gray-100 cursor-pointer"
                                    onClick={() => handleOrderChange(order)}
                                >
                                    <span>{order}</span>
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}