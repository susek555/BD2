'use client';

import { useState } from "react";
import { useEffect } from "react";

export default function Photos({ imagesURLs }: { imagesURLs: string[] }) {
    const [currentIndex, setCurrentIndex] = useState(0);
    const [imageError, setImageError] = useState<boolean[]>(Array(imagesURLs.length).fill(false));

    // Zmiana obrazu na poprzedni
    const handlePrev = () => {
        setCurrentIndex((prevIndex) => (prevIndex === 0 ? imagesURLs.length - 1 : prevIndex - 1));
    };

    // Zmiana obrazu na następny
    const handleNext = () => {
        setCurrentIndex((prevIndex) => (prevIndex === imagesURLs.length - 1 ? 0 : prevIndex + 1));
    };

    // Obsługa błędu ładowania obrazu
    const handleError = (index: number) => {
        setImageError((prev) => {
            const newErrors = [...prev];
            newErrors[index] = true;
            return newErrors;
        });
        console.error(`Nie udało się załadować obrazu: ${imagesURLs[index]}`);
    };

    return (
        <div className="flex flex-col items-center gap-4">
            <div className="relative w-full">
                {imageError[currentIndex] ? (
                    <div className="w-full h-120 flex items-center justify-center bg-gray-300 text-red-500">
                        Błąd ładowania obrazu
                    </div>
                ) : (
                    <img
                        src={imagesURLs[currentIndex]}
                        alt={`Image ${currentIndex + 1}`}
                        onError={() => handleError(currentIndex)}
                        className="w-full h-auto md:h-120 object-cover"
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
