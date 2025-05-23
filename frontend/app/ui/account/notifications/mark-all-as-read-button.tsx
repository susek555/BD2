'use client'

import { usePathname, useRouter, useSearchParams } from "next/navigation";

export default function MarkAllAsReadButton() {
  const searchParams = useSearchParams();
  const pathname = usePathname();
  const { replace } = useRouter();

  const handleClick = () => {
    //TODO Update through API

    const params = new URLSearchParams(searchParams);
    replace(`${pathname}?${params.toString()}`);
  }

    return (
        <button
            className="bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600"
            onClick={() => handleClick()}
        >
            Mark all as read
        </button>
    )
}