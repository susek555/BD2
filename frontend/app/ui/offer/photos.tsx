'use client';

// To mock images via url, run http-server . --cors in public/offer

import { useState } from "react";

export default function Photos({ imagesURLs }: { imagesURLs: string[] }) {
    const [currentIndex, setCurrentIndex] = useState(0);
    const [imageSrc, setImageSrc] = useState<string[]>([...imagesURLs]);

    const handlePrev = () => {
        setCurrentIndex((prevIndex) => (prevIndex === 0 ? imagesURLs.length - 1 : prevIndex - 1));
    };

    const handleNext = () => {
        setCurrentIndex((prevIndex) => (prevIndex === imagesURLs.length - 1 ? 0 : prevIndex + 1));
    };

    const handleError = (index: number) => {
        const fallbackImage = "/(home)/car_placeholder.png";
        setImageSrc((prev) => {
            const newSrc = [...prev];
            newSrc[index] = fallbackImage;
            return newSrc;
        });
        console.log(`Failed to load image: ${imagesURLs[index]}. Placeholder was loaded.`);
    };

    return (
        <div className="flex flex-col items-center gap-4">
            <div className="relative w-full md:h-120 border border-gray-300 flex justify-center items-center">
                <img
                    src={imageSrc[currentIndex]}
                    alt={`Image ${currentIndex + 1}`}
                    onError={() => handleError(currentIndex)}
                    className="max-w-full max-h-full"
                />
            </div>
            <div className="flex items-center gap-2">
                <button
                    onClick={handlePrev}
                    className="bg-gray-800 text-white px-2 py-1 rounded"
                >
                    &#8592;
                </button>
                <p>
                    {currentIndex + 1} / {imagesURLs.length}
                </p>
                <button
                    onClick={handleNext}
                    className="bg-gray-800 text-white px-2 py-1 rounded"
                >
                    &#8594;
                </button>
            </div>
        </div>
    );
}
