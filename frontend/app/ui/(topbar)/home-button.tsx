'use client'

import Link from 'next/link';

export default function HomeButton() {
    return (
        <Link
            className="mb-2 flex w-20 items-end justify-start rounded-br-lg bg-blue-600 md:w-80 md:h-12"
            href="/"
            prefetch={false}
            onClick={() => {
            // Force a hard refresh when clicking HOME
                window.location.href = '/';
            }}
        >
            <div className="flex h-full w-full items-center justify-center rounded-md bg-blue-600 p-0 text-white">
            <p className="font-bold">HOME</p>
            </div>
        </Link>
    )
};