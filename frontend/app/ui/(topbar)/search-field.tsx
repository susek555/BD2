'use client';

import { MagnifyingGlassIcon } from '@heroicons/react/24/outline';
import { useSearchParams, usePathname, useRouter } from 'next/navigation';
import { useDebouncedCallback } from 'use-debounce';

export default function Search({ placeholder }: { placeholder: string }) {
  const searchParams = useSearchParams();
  const pathname = usePathname();
  const { replace } = useRouter();

  const handleSearch = useDebouncedCallback((term) => {
    const params = new URLSearchParams(searchParams);

    params.set('page', '1'); // Reset to the first page
    if (term) {
      params.set('query', term);
    } else {
      params.delete('query');
    }
    // replace(`${pathname}?${params.toString()}`);
    window.history.replaceState(null, '', `${pathname}?${params.toString()}`);
  }, 300);

  const handleConfirmSearch = () => {
    const params = new URLSearchParams(searchParams);
    replace(`${pathname}?${params.toString()}`);
  }

  return (
    <div className="relative flex flex-1 h-12 flex-shrink-0">
      <label htmlFor="search" className="sr-only">
        Search
      </label>
      <input
        className="peer block w-full rounded-b-lg border border-gray-200 py-[9px] pl-10 text-sm outline-2 placeholder:text-gray-500"
        placeholder={placeholder}
        onChange={(e) => {
          handleSearch(e.target.value);
        }}
        defaultValue={searchParams.get('query')?.toString()}
      />
      <MagnifyingGlassIcon className="absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-gray-500 peer-focus:text-gray-900" />
      <button
        type="button"
        className="absolute right-2 top-1/2 -translate-y-1/2 px-3 py-1 rounded bg-blue-500 text-white text-sm hover:bg-blue-600"
        onClick={handleConfirmSearch}
      >
        Search
      </button>
    </div>
  );
}