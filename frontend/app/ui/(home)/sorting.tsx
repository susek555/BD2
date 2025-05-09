import React from "react";

export default function Sorting() {
    const [isKeyDropdownOpen, setIsKeyDropdownOpen] = React.useState(false);
    const [isOrderDropdownOpen, setIsOrderDropdownOpen] = React.useState(false);
    const [selectedKey, setSelectedKey] = React.useState("name");
    const [selectedOrder, setSelectedOrder] = React.useState("Asc");

    const toggleKeyDropdown = () => setIsKeyDropdownOpen(!isKeyDropdownOpen);
    const toggleOrderDropdown = () => setIsOrderDropdownOpen(!isOrderDropdownOpen);

    const handleKeyChange = (key: string) => {
        setSelectedKey(key);
        setIsKeyDropdownOpen(false);
    };

    const handleOrderChange = (order: string) => {
        setSelectedOrder(order);
        setIsOrderDropdownOpen(false);
    };

    return (
        <div>
            <h2 className="text-bg px-2">Sort by:</h2>
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
                            {["name", "date", "price"].map((key) => (
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
                        <span>{selectedOrder}</span>
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