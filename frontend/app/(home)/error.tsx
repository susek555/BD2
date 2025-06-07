'use client'

import { useEffect } from "react";

export default function Error({ error }: { error: Error & { digest?: string } }) {
    useEffect(() => {
        console.error("Error occurred:", error, error.digest);
    }, [error]);

    return (
        <div className="flex flex-col items-center justify-center h-full">
            <div className="my-10" />
            <h2 className="text-2xl font-bold text-red-600">Error</h2>
            <p className="text-lg text-gray-700">{error.message}</p>
            <button
                onClick={() => window.location.href = '/'}
                className="mt-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors"
            >
                Back to Home
            </button>
        </div>
    );
}