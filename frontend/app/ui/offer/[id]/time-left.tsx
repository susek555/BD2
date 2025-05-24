'use client';

import React, { useState, useEffect } from 'react';

export default function TimeLeft({ endDate }: { endDate: Date }) {
    const [timeLeft, setTimeLeft] = useState<string>("");

    useEffect(() => {
        const updateCountdown = () => {
            const now = new Date();
            const diff = new Date(endDate).getTime() - now.getTime();

            if (diff <= 0) {
                setTimeLeft("Time is up!");
                return;
            }

            const days = Math.floor(diff / (1000 * 60 * 60 * 24));
            const hours = Math.floor((diff / (1000 * 60 * 60)) % 24);
            const minutes = Math.floor((diff / (1000 * 60)) % 60);
            const seconds = Math.floor((diff / 1000) % 60);

            setTimeLeft(
                `${days}d ${hours}h ${minutes}m ${seconds}s`
            );
        };

        updateCountdown();
        const timerId = setInterval(updateCountdown, 1000);

        return () => clearInterval(timerId); // Cleanup on unmount
    }, [endDate]);

    return (
        <div className="flex flex-col gap-2 justify-center items-center">
            <p className="text-2xl">Time left</p>
            <p className="font-bold text-lg">{timeLeft}</p>
        </div>
    );
}
