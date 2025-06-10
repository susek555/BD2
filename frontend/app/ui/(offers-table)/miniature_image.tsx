'use client'

import { useEffect, useState } from "react";

interface MiniatureImageProps {
    url: string;
}

export default function MainImage({ url }: MiniatureImageProps) {
    const [imageSrc, setImageSrc] = useState<string>("");
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const img = new Image();
        img.src = url;

        img.onload = () => {
            setImageSrc(url);
            setLoading(false);
        };

        img.onerror = () => {
            const fallbackImage = "/(home)/car_placeholder.png";
            setImageSrc(fallbackImage);
            console.log(`Failed to load image: ${url}. Placeholder was loaded.`);
            setLoading(false);
        };
    }, [url]);

    return (
        <div className="relative h-35 w-70 object-cover rounded">
            {loading && (
                <div className="absolute inset-0 flex items-center justify-center bg-gray-100">
                    <div className="loader animate-spin rounded-full border-4 border-gray-300 border-t-gray-600 w-12 h-12"></div>
                </div>
            )}
            {!loading && (
                <img
                    src={imageSrc}
                    alt="Car image"
                    className="h-35 w-70 object-cover rounded"
                />
            )}
        </div>
    );
}