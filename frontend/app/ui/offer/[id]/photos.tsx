'use client';

// To mock images via url, run http-server . --cors in public/offer

import { useEffect, useState } from "react";

export default function Photos({ imagesURLs }: { imagesURLs: string[] }) {
    const [currentIndex, setCurrentIndex] = useState(0);
    const [imageSrc, setImageSrc] = useState<string[]>([...imagesURLs]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        setLoading(true);

        const img = new Image();
        img.src = imageSrc[currentIndex];

        img.onload = () => {
            setLoading(false);
        };

        img.onerror = () => {
            const fallbackImage = "/(home)/car_placeholder.png";
            setImageSrc((prev) => {
                const newSrc = [...prev];
                newSrc[currentIndex] = fallbackImage;
                return newSrc;
            });
            console.log(`Failed to load image: ${imagesURLs[currentIndex]}. Placeholder was loaded.`);
            setLoading(false);
        };
    }, [currentIndex, imagesURLs, imageSrc]);

    const handlePrev = () => {
        setCurrentIndex((prevIndex) =>
            prevIndex === 0 ? imagesURLs.length - 1 : prevIndex - 1
        );
    };

    const handleNext = () => {
        setCurrentIndex((prevIndex) =>
            prevIndex === imagesURLs.length - 1 ? 0 : prevIndex + 1
        );
    };

    return (
        <div className="flex flex-col items-center gap-4">
            <div className="relative w-full md:h-130 border border-gray-300 flex justify-center items-center">
                {loading && (
                    <div className="absolute inset-0 flex items-center justify-center bg-gray-100">
                        <div className="loader animate-spin rounded-full border-4 border-gray-300 border-t-gray-600 w-12 h-12"></div>
                    </div>
                )}
                {!loading && (
                    <img
                        src={imageSrc[currentIndex]}
                        alt={`Image ${currentIndex + 1}`}
                        className="max-w-full max-h-full"
                    />
                )}
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
