'use client'

import { usePathname, useRouter, useSearchParams } from "next/navigation";

export default function ApplyButton() {
  const searchParams = useSearchParams();
  const pathname = usePathname();
  const { replace } = useRouter();

  const handleApply = () => {
    const params = new URLSearchParams(searchParams);
    replace(`${pathname}?${params.toString()}`);
  }

    return (
        <button
            className="bg-blue-500 text-white font-bold py-2 px-4 rounded"
            onClick={() => handleApply()}
        >
            Apply
        </button>
    )
}