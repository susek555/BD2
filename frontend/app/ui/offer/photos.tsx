'use client'

import { useState } from "react";

export default function Photos({ imagesURLs }: { imagesURLs: string[] }) {
    const [currentIndex, setCurrentIndex] = useState(0);

    const handlePrev = () => {
        setCurrentIndex((prevIndex) => (prevIndex === 0 ? imagesURLs.length - 1 : prevIndex - 1));
    };

    const handleNext = () => {
        setCurrentIndex((prevIndex) => (prevIndex === imagesURLs.length - 1 ? 0 : prevIndex + 1));
    };

    return (
        <div className="flex flex-col items-center gap-4">
            <div className="relative w-full max-w-lg">
                <img
                    src={imagesURLs[currentIndex]}
                    alt={`Image ${currentIndex + 1}`}
                    className="w-full h-auto object-cover"
                />
                <button
                    onClick={handlePrev}
                    className="absolute left-0 top-1/2 transform -translate-y-1/2 bg-gray-800 text-white px-2 py-1 rounded"
                >
                    &#8592;
                </button>
                <button
                    onClick={handleNext}
                    className="absolute right-0 top-1/2 transform -translate-y-1/2 bg-gray-800 text-white px-2 py-1 rounded"
                >
                    &#8594;
                </button>
            </div>
            <p>
                {currentIndex + 1} / {imagesURLs.length}
            </p>
        </div>
    );
}